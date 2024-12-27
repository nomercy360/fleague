package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	initdata "github.com/telegram-mini-apps/init-data-golang"
	"github.com/user/project/internal/contract"
	"github.com/user/project/internal/db"
	"github.com/user/project/internal/nanoid"
	"github.com/user/project/internal/terrors"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func (s Service) TelegramAuth(req contract.AuthTelegramRequest) (*contract.UserAuthResponse, error) {
	expIn := 24 * time.Hour
	botToken := s.cfg.BotToken

	if err := initdata.Validate(req.Query, botToken, expIn); err != nil {
		return nil, terrors.Unauthorized(err, "invalid init data from telegram")
	}

	data, err := initdata.Parse(req.Query)

	if err != nil {
		return nil, terrors.Unauthorized(err, "cannot parse init data from telegram")
	}

	user, err := s.storage.GetUserByChatID(data.User.ID)
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

		// if referrer is not empty, get referrer user by ID
		var referrerID *string
		if req.ReferrerID != nil {
			referrer, err := s.storage.GetUserByID(*req.ReferrerID)
			if err != nil && errors.Is(err, db.ErrNotFound) {
				log.Printf("referrer not found: %v", err)
			} else if err != nil {
				log.Printf("failed to get referrer: %v", err)
			}

			if referrer != nil {
				referrerID = &referrer.ID

				// add 10 points to referrer
				if err = s.storage.UpdateUserPoints(context.Background(), referrer.ID, 10); err != nil {
					log.Printf("failed to update referrer points: %v", err)
				}
			}
		}

		imgUrl := fmt.Sprintf("%s/avatars/%d.svg", s.cfg.AssetsURL, rand.Intn(30)+1)

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

		if err = s.storage.CreateUser(create); err != nil {
			return nil, terrors.InternalServer(err, "failed to create user")
		}

		user, err = s.storage.GetUserByChatID(data.User.ID)
		if err != nil {
			return nil, terrors.InternalServer(err, "failed to get user")
		}

		//if data.User.PhotoURL != "" {
		//	go func() {
		//		imgFile := fmt.Sprintf("fb/users/%s.jpg", nanoid.Must())
		//		imgUrl := fmt.Sprintf("%s/%s", s.cfg.AssetsURL, imgFile)
		//		if err = s.uploadImageToS3(data.User.PhotoURL, imgFile); err != nil {
		//			log.Printf("failed to upload user avatar to S3: %v", err)
		//		}
		//
		//		if err = s.storage.UpdateUserAvatarURL(context.Background(), data.User.ID, imgUrl); err != nil {
		//			log.Printf("failed to update user avatar URL: %v", err)
		//		}
		//	}()
		//}
	} else if err != nil {
		return nil, terrors.InternalServer(err, "failed to get user")
	}

	token, err := generateJWT(user.ID, user.ChatID, s.cfg.JWTSecret)

	if err != nil {
		return nil, terrors.InternalServer(err, "jwt library error")
	}

	uresp := contract.UserResponse{
		ID:                 user.ID,
		FirstName:          user.FirstName,
		LastName:           user.LastName,
		Username:           user.Username,
		LanguageCode:       user.LanguageCode,
		ChatID:             user.ChatID,
		CreatedAt:          user.CreatedAt,
		TotalPoints:        user.TotalPoints,
		TotalPredictions:   user.TotalPredictions,
		CorrectPredictions: user.CorrectPredictions,
		AvatarURL:          user.AvatarURL,
		ReferredBy:         user.ReferredBy,
		GlobalRank:         user.GlobalRank,
	}

	return &contract.UserAuthResponse{
		Token: token,
		User:  uresp,
	}, nil
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

func (s Service) uploadImageToS3(imgURL string, fileName string) error {
	resp, err := http.Get(imgURL)

	if err != nil {
		return fmt.Errorf("failed to download file: %v", err)

	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	if _, err = s.s3Client.UploadFile(data, fileName); err != nil {
		return fmt.Errorf("failed to upload user avatar to S3: %v", err)
	}

	return nil
}
