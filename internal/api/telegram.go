package api

import (
	"context"
	"fmt"
	telegram "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/labstack/echo/v4"
	"github.com/user/project/internal/db"
	"github.com/user/project/internal/terrors"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// SendInvoice sends an invoice to the user for purchasing prediction tokens
func (a *API) SendInvoice(c echo.Context) error {
	ctx := c.Request().Context()
	uid := GetContextUserID(c)

	//user, err := a.storage.GetUserByID(uid)
	//if err != nil {
	//	return terrors.InternalServer(err, "failed to get user")
	//}

	// Define token packages (e.g., 50 tokens for 10 XTR)
	amount := 1       // XTR (Telegram Stars)
	tokenAmount := 50 // Prediction tokens

	invoice := telegram.CreateInvoiceLinkParams{
		Title:       "Prediction Tokens Purchase",
		Description: fmt.Sprintf("Buy %d prediction tokens", tokenAmount),
		Payload:     fmt.Sprintf("purchase:%s:%d", uid, tokenAmount), // Internal payload
		Currency:    "XTR",
		Prices: []models.LabeledPrice{
			{Label: "Prediction Tokens", Amount: amount},
		},
	}

	link, err := a.tg.CreateInvoiceLink(ctx, &invoice)
	if err != nil {
		return terrors.InternalServer(err, "failed to send invoice")
	}

	return c.JSON(http.StatusOK, echo.Map{"status": "ok", "link": link})
}

// HandlePreCheckoutQuery processes the pre-checkout query
func (a *API) HandlePreCheckoutQuery(update models.Update) error {
	var resp telegram.AnswerPreCheckoutQueryParams
	if update.PreCheckoutQuery == nil {
		return nil
	}

	query := update.PreCheckoutQuery
	ctx := context.Background()

	// Validate payload
	parts := strings.Split(query.InvoicePayload, ":")
	if len(parts) != 3 || parts[0] != "purchase" {
		resp = telegram.AnswerPreCheckoutQueryParams{
			PreCheckoutQueryID: query.ID,
			OK:                 false,
			ErrorMessage:       "Invalid purchase request",
		}

		_, err := a.tg.AnswerPreCheckoutQuery(ctx, &resp)

		if err != nil {
			log.Printf("failed to reject payment: %v\n", err)
			return nil
		}

		return nil
	}

	uid := parts[1]
	//tokenAmount, _ := strconv.Atoi(parts[2])

	// Verify user exists
	_, err := a.storage.GetUserByID(uid)
	if err != nil {
		log.Printf("failed to get user: %v\n", err)
		return nil
	}

	// Approve the payment
	resp = telegram.AnswerPreCheckoutQueryParams{
		PreCheckoutQueryID: query.ID,
		OK:                 true,
	}

	ok, err := a.tg.AnswerPreCheckoutQuery(ctx, &resp)
	if err != nil {
		log.Printf("failed to approve payment: %v\n", err)
		return nil
	}

	if !ok {
		log.Printf("failed to approve payment: %v\n", err)
		return nil
	}

	return nil
}

// HandleSuccessfulPayment processes the successful payment
func (a *API) HandleSuccessfulPayment(update models.Update) error {
	if update.Message == nil || update.Message.SuccessfulPayment == nil {
		return nil
	}

	payment := update.Message.SuccessfulPayment
	ctx := context.Background()

	// Extract payload
	parts := strings.Split(payment.InvoicePayload, ":")
	if len(parts) != 3 || parts[0] != "purchase" {
		return nil // Silently ignore invalid payload
	}

	uid := parts[1]
	tokenAmount, _ := strconv.Atoi(parts[2])

	// Convert XTR to prediction tokens and update balance
	balance, err := a.storage.UpdateUserTokens(
		ctx,
		uid,
		tokenAmount,
		db.TokenTransactionTypePurchase,
	)

	if err != nil {
		log.Printf("failed to update user tokens: %v\n", err)
		return nil
	}

	msg := telegram.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("Payment successful! You've received %d prediction tokens. New balance: %d", tokenAmount, balance),
	}

	// Notify user
	_, err = a.tg.SendMessage(ctx, &msg)

	// for testing purposes do a refund in telegram
	go func() {
		_, err = a.tg.RefundStarPayment(ctx, &telegram.RefundStarPaymentParams{
			UserID:                  update.Message.Chat.ID,
			TelegramPaymentChargeID: payment.TelegramPaymentChargeID,
		})

		if err != nil {
			log.Printf("failed to refund payment: %v\n", err)
		}
	}()

	return err
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
