package api

import (
	"context"
	"fmt"
	telegram "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/labstack/echo/v4"
	"github.com/user/project/internal/db"
	"github.com/user/project/internal/nanoid"
	"github.com/user/project/internal/terrors"
	"log"
	"net/http"
	"strings"
	"time"
)

// SendInvoice отправляет счет пользователю для покупки подписки на месяц
func (a *API) SendInvoice(c echo.Context) error {
	ctx := c.Request().Context()
	uid := GetContextUserID(c)

	// Стоимость подписки в Telegram Stars (XTR)
	amount := 1 // 1 XTR за подписку на месяц

	invoice := telegram.CreateInvoiceLinkParams{
		Title:       "Monthly Subscription",
		Description: "Get access to predictions for 30 days",
		Payload:     fmt.Sprintf("subscription:%s", uid), // Payload для подписки
		Currency:    "XTR",
		Prices: []models.LabeledPrice{
			{Label: "Monthly Subscription", Amount: amount},
		},
	}

	link, err := a.tg.CreateInvoiceLink(ctx, &invoice)
	if err != nil {
		return terrors.InternalServer(err, "failed to send invoice")
	}

	return c.JSON(http.StatusOK, echo.Map{"status": "ok", "link": link})
}

// HandlePreCheckoutQuery обрабатывает запрос перед оплатой
func (a *API) HandlePreCheckoutQuery(update models.Update) error {
	var resp telegram.AnswerPreCheckoutQueryParams
	if update.PreCheckoutQuery == nil {
		return nil
	}

	query := update.PreCheckoutQuery
	ctx := context.Background()

	// Проверяем payload
	parts := strings.Split(query.InvoicePayload, ":")
	if len(parts) != 2 || parts[0] != "subscription" {
		resp = telegram.AnswerPreCheckoutQueryParams{
			PreCheckoutQueryID: query.ID,
			OK:                 false,
			ErrorMessage:       "Invalid subscription request",
		}
		_, err := a.tg.AnswerPreCheckoutQuery(ctx, &resp)
		if err != nil {
			log.Printf("failed to reject payment: %v\n", err)
		}
		return nil
	}

	uid := parts[1]

	// Проверяем существование пользователя
	_, err := a.storage.GetUserByID(uid)
	if err != nil {
		log.Printf("failed to get user: %v\n", err)
		resp = telegram.AnswerPreCheckoutQueryParams{
			PreCheckoutQueryID: query.ID,
			OK:                 false,
			ErrorMessage:       "User not found",
		}
		_, err = a.tg.AnswerPreCheckoutQuery(ctx, &resp)
		if err != nil {
			log.Printf("failed to reject payment: %v\n", err)
		}
		return nil
	}

	// Подтверждаем оплату
	resp = telegram.AnswerPreCheckoutQueryParams{
		PreCheckoutQueryID: query.ID,
		OK:                 true,
	}
	_, err = a.tg.AnswerPreCheckoutQuery(ctx, &resp)
	if err != nil {
		log.Printf("failed to approve payment: %v\n", err)
		return nil
	}

	return nil
}

// HandleSuccessfulPayment обрабатывает успешную оплату
func (a *API) HandleSuccessfulPayment(update models.Update) error {
	if update.Message == nil || update.Message.SuccessfulPayment == nil {
		return nil
	}

	payment := update.Message.SuccessfulPayment
	ctx := context.Background()

	// Извлекаем payload
	parts := strings.Split(payment.InvoicePayload, ":")
	if len(parts) != 2 || parts[0] != "subscription" {
		return nil // Игнорируем некорректный payload
	}

	uid := parts[1]

	// Активируем подписку на 30 дней
	user, err := a.storage.GetUserByID(uid)
	if err != nil {
		log.Printf("failed to get user: %v\n", err)
		return nil
	}

	// Рассчитываем новую дату окончания подписки
	now := time.Now()
	newExpiry := now.AddDate(0, 1, 0) // +1 месяц
	if user.SubscriptionActive && user.SubscriptionExpiry.After(now) {
		// Если подписка уже активна, добавляем 30 дней к текущей дате окончания
		newExpiry = user.SubscriptionExpiry.AddDate(0, 1, 0)
	}

	// Обновляем статус подписки пользователя
	err = a.storage.UpdateUserSubscription(ctx, uid, true, newExpiry)
	if err != nil {
		log.Printf("failed to update user subscription: %v\n", err)
		return nil
	}

	subscription := db.Subscription{
		ID:        nanoid.Must(),
		UserID:    uid,
		StartDate: now,
		EndDate:   newExpiry,
		IsPaid:    true,
		CreatedAt: now,
	}

	err = a.storage.SaveSubscription(ctx, subscription)
	if err != nil {
		log.Printf("failed to save subscription: %v\n", err)
		return nil
	}

	var messageText string
	switch *user.LanguageCode {
	case "ru":
		messageText = fmt.Sprintf(
			"Оплата прошла успешно! Ваша подписка активна до %s",
			newExpiry.Format("02.01.2006"),
		)
	default:
		messageText = fmt.Sprintf(
			"Payment successful! Your subscription is active until %s",
			newExpiry.Format("2006-01-02"),
		)
	}

	msg := telegram.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   messageText,
	}

	_, err = a.tg.SendMessage(ctx, &msg)
	if err != nil {
		log.Printf("failed to send message: %v\n", err)
	}

	// Для тестов возвращаем деньги (оставлено как было)
	go func() {
		_, err = a.tg.RefundStarPayment(ctx, &telegram.RefundStarPaymentParams{
			UserID:                  update.Message.Chat.ID,
			TelegramPaymentChargeID: payment.TelegramPaymentChargeID,
		})
		if err != nil {
			log.Printf("failed to refund payment: %v\n", err)
		}
	}()

	return nil
}

func (a *API) TelegramWebhook(c echo.Context) error {
	var update models.Update
	if err := c.Bind(&update); err != nil {
		return terrors.BadRequest(err, "failed to decode update")
	}

	if err := a.HandlePreCheckoutQuery(update); err != nil {
		return terrors.InternalServer(err, "failed to handle pre-checkout query")
	}

	if err := a.HandleSuccessfulPayment(update); err != nil {
		return terrors.InternalServer(err, "failed to handle successful payment")
	}

	return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
}
