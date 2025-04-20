package api

import (
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/labstack/echo/v4"
	"github.com/user/project/internal/db"
	"log"
	"sync"
	"time"
)

func (a *API) BroadcastSubscriptionMessage(c echo.Context) error {
	isTestRun := c.QueryParam("test") == "true"
	if isTestRun {
		log.Println("Test run for broadcast message")
	}

	users, err := a.storage.GetAllUsers(c.Request().Context())
	if err != nil {
		log.Printf("Failed to get users: %v", err)
		return fmt.Errorf("failed to get users: %w", err)
	}

	messageRu := `🎉 *Друзья, важное обновление\!*

Теперь для ставок нужна подписка за *150 Telegram Stars* в месяц\. Это позволит нам награждать лучших предикторов *ценными призами* каждый месяц и улучшать бота\!  

🔥 Подписывайтесь и покажите свое *футбольное чутье*\!`

	messageEn := `
🎉 *Big update, friends\!*

To place bets, you\’ll need a subscription for *150 Telegram Stars* per month\. This helps us reward top predictors with *awesome prizes* monthly and keep improving the bot\!  

🔥 Subscribe now and unleash your *football instincts*\!`

	batchSize := 10
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, batchSize)

	for i := 0; i < len(users); i += batchSize {
		end := i + batchSize
		if end > len(users) {
			end = len(users)
		}
		batch := users[i:end]

		for _, user := range batch {
			wg.Add(1)
			semaphore <- struct{}{}

			go func(user db.User) {
				defer wg.Done()
				defer func() { <-semaphore }()

				message := messageRu
				if user.LanguageCode != nil && *user.LanguageCode == "en" {
					message = messageEn
				}

				msg := bot.SendPhotoParams{
					ChatID:    user.ChatID,
					Caption:   message,
					ParseMode: models.ParseModeMarkdown,
					Photo:     &models.InputFileString{Data: "https://assets.peatch.io/preview.png"},
					ReplyMarkup: &models.InlineKeyboardMarkup{
						InlineKeyboard: [][]models.InlineKeyboardButton{
							{
								{
									Text: "В приложение",
									WebApp: &models.WebAppInfo{
										URL: "https://fleague.mxksimdev.com",
									},
								},
							},
						},
					},
				}

				if user.ChatID != 927635965 && isTestRun {
					log.Printf("Test run skipped for user %s", user.ID)
					return
				}

				_, err := a.tg.SendPhoto(c.Request().Context(), &msg)
				if err != nil {
					log.Printf("Failed to send message to user %s: %v", user.ID, err)
				}
			}(user)
		}

		wg.Wait()

		if end < len(users) {
			time.Sleep(1 * time.Second)
		}
	}

	log.Printf("Broadcast message sent to %d users", len(users))
	return c.NoContent(200)
}
