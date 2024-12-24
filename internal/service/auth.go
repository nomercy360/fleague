package service

import (
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

func (s Service) TelegramAuth(query string) (*contract.UserAuthResponse, error) {
	expIn := 24 * time.Hour
	botToken := s.cfg.BotToken

	if err := initdata.Validate(query, botToken, expIn); err != nil {
		return nil, terrors.Unauthorized(err, "invalid init data from telegram")
	}

	data, err := initdata.Parse(query)

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

		imgUrl := fmt.Sprintf("%s/avatars/%d.svg", s.cfg.AssetsURL, rand.Intn(30)+1)

		if data.User.PhotoURL != "" {
			imgFile := fmt.Sprintf("fb/users/%s.jpg", nanoid.Must())
			imgUrl = fmt.Sprintf("%s/%s", s.cfg.AssetsURL, imgFile)
			go func() {
				if err = s.uploadImageToS3(data.User.PhotoURL, imgFile); err != nil {
					log.Printf("failed to upload user avatar to S3: %v", err)
				}
			}()
		}

		create := db.User{
			ID:           nanoid.Must(),
			Username:     username,
			ChatID:       data.User.ID,
			ReferralCode: nanoid.Must(),
			FirstName:    first,
			LastName:     last,
			LanguageCode: &lang,
			AvatarURL:    &imgUrl,
		}

		if err = s.storage.CreateUser(create); err != nil {
			return nil, terrors.InternalServer(err, "failed to create user")
		}

		user, err = s.storage.GetUserByChatID(data.User.ID)
		if err != nil {
			return nil, terrors.InternalServer(err, "failed to get user")
		}
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
		ReferralCode:       user.ReferralCode,
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
