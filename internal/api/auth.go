package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	initdata "github.com/telegram-mini-apps/init-data-golang"
	"github.com/user/project/internal/contract"
	"github.com/user/project/internal/db"
	"github.com/user/project/internal/nanoid"
	"github.com/user/project/internal/terrors"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"
)

func (a *API) TelegramAuth(c echo.Context) error {
	var req contract.AuthTelegramRequest
	if err := c.Bind(&req); err != nil {
		return terrors.BadRequest(err, "failed to bind request")
	}

	if err := req.Validate(); err != nil {
		return terrors.BadRequest(err, "failed to validate request")
	}

	log.Printf("AuthTelegram: %+v", req)

	expIn := 24 * time.Hour
	botToken := a.cfg.BotToken

	if err := initdata.Validate(req.Query, botToken, expIn); err != nil {
		return terrors.Unauthorized(err, "invalid init data from telegram")
	}

	data, err := initdata.Parse(req.Query)
	if err != nil {
		return terrors.Unauthorized(err, "cannot parse init data from telegram")
	}

	user, err := a.storage.GetUserByChatID(data.User.ID)
	if err != nil && errors.Is(err, db.ErrNotFound) {
		username := data.User.Username
		if username == "" {
			username = "user_" + fmt.Sprintf("%d", data.User.ID)
		}

		var first, last *string
		if data.User.FirstName != "" {
			first = &data.User.FirstName
		}
		if data.User.LastName != "" {
			last = &data.User.LastName
		}

		lang := "ru"
		if data.User.LanguageCode != "ru" {
			lang = "en"
		}

		var referrerID *string
		if req.ReferrerID != nil {
			referrer, err := a.storage.GetUserByID(*req.ReferrerID)
			if err != nil && errors.Is(err, db.ErrNotFound) {
				log.Printf("referrer not found: %v", err)
			} else if err != nil {
				log.Printf("failed to get referrer: %v", err)
			}
			if referrer.ID != "" {
				referrerID = &referrer.ID
				balance, err := a.storage.UpdateUserTokens(context.Background(), referrer.ID, 50, db.TokenTransactionTypeReferral)
				if err != nil {
					log.Printf("Failed to award referral bonus for user %s: %v", referrer.ID, err)
				}

				referrer.PredictionTokens = balance
			}
		}

		imgUrl := fmt.Sprintf("%s/avatars/%d.svg", a.cfg.AssetsURL, rand.Intn(30)+1)
		create := db.User{
			ID:           nanoid.Must(),
			Username:     username,
			ChatID:       data.User.ID,
			FirstName:    first,
			LastName:     last,
			LanguageCode: &lang,
			AvatarURL:    &imgUrl,
			ReferredBy:   referrerID,
		}

		if err = a.storage.CreateUser(create); err != nil {
			return terrors.InternalServer(err, "failed to create user")
		}

		user, err = a.storage.GetUserByChatID(data.User.ID)
		if err != nil {
			return terrors.InternalServer(err, "failed to get user")
		}
	} else if err != nil {
		return terrors.InternalServer(err, "failed to get user")
	}

	ctx := c.Request().Context()
	hasLoggedInToday, err := a.storage.HasLoggedInToday(ctx, user.ID)
	if err != nil {
		log.Printf("Failed to check daily login for user %s: %v", user.ID, err)
	}

	if err := a.storage.RecordUserLogin(ctx, user.ID); err != nil {
		log.Printf("Failed to record login for user %s: %v", user.ID, err)
	}

	if err != nil {
		log.Printf("Failed to check daily login for user %s: %v", user.ID, err)

	} else if !hasLoggedInToday && user.CreatedAt.Before(time.Now().Add(-24*time.Hour)) {
		balance, err := a.storage.UpdateUserTokens(ctx, user.ID, 5, db.TokenTransactionTypeDailyLogin)

		if err != nil {
			log.Printf("Failed to award daily login bonus for user %s: %v", user.ID, err)
		} else {
			log.Printf("Awarded 5 Prediction Tokens to user %s for daily login", user.ID)
		}

		user.PredictionTokens = balance
	}

	token, err := generateJWT(user.ID, user.ChatID, a.cfg.JWTSecret)
	if err != nil {
		return terrors.InternalServer(err, "jwt library error")
	}

	ranks, err := a.storage.GetUserRank(context.Background(), user.ID)
	if err != nil {
		return terrors.InternalServer(err, "failed to get user rank")
	}

	uresp := contract.UserResponse{
		ID:                 user.ID,
		FirstName:          user.FirstName,
		LastName:           user.LastName,
		Username:           user.Username,
		LanguageCode:       user.LanguageCode,
		ChatID:             user.ChatID,
		CreatedAt:          user.CreatedAt,
		TotalPredictions:   user.TotalPredictions,
		CorrectPredictions: user.CorrectPredictions,
		AvatarURL:          user.AvatarURL,
		ReferredBy:         user.ReferredBy,
		Ranks:              ranks,
		FavoriteTeam:       user.FavoriteTeam,
		CurrentWinStreak:   user.CurrentWinStreak,
		LongestWinStreak:   user.LongestWinStreak,
		Badges:             user.Badges,
		PredictionTokens:   user.PredictionTokens,
	}

	if uresp.TotalPredictions > 0 {
		accuracy := (float64(uresp.CorrectPredictions) / float64(uresp.TotalPredictions)) * 100
		uresp.PredictionAccuracy = math.Round(accuracy*100) / 100
	}

	resp := &contract.UserAuthResponse{
		Token: token,
		User:  uresp,
	}

	return c.JSON(http.StatusOK, resp)
}

func generateJWT(userID string, chatID int64, secretKey string) (string, error) {
	claims := &contract.JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
		UID:    userID,
		ChatID: chatID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return t, nil
}

func (a *API) uploadImageToS3(imgURL string, fileName string) error {
	resp, err := http.Get(imgURL)

	if err != nil {
		return fmt.Errorf("failed to download file: %v", err)

	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	if _, err = a.s3.UploadFile(data, fileName); err != nil {
		return fmt.Errorf("failed to upload user avatar to S3: %v", err)
	}

	return nil
}
