package service

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	gonanoid "github.com/matoous/go-nanoid/v2"
	initdata "github.com/telegram-mini-apps/init-data-golang"
	"github.com/user/project/internal/contract"
	"github.com/user/project/internal/db"
	"io"
	"math/rand"
	"net/http"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateReferralCode(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := make([]byte, length)
	for i := range code {
		code[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(code)
}

func (s Service) TelegramAuth(query string) (*contract.UserAuthResponse, error) {
	expIn := 24 * time.Hour
	botToken := s.botToken

	if err := initdata.Validate(query, botToken, expIn); err != nil {
		return nil, contract.ErrUnauthorized
	}

	data, err := initdata.Parse(query)

	if err != nil {
		return nil, contract.ErrUnauthorized
	}

	user, err := s.storage.GetUserByChatID(data.User.ID)
	if err != nil && errors.Is(err, db.ErrNotFound) {
		var firstName, lastName *string

		if data.User.FirstName != "" {
			firstName = &data.User.FirstName
		}

		if data.User.LastName != "" {
			lastName = &data.User.LastName
		}

		username := data.User.Username
		if username == "" {
			username = "user_" + fmt.Sprintf("%d", data.User.ID)
		}

		var cdnPath string

		if data.User.PhotoURL == "" {
			imgFile := fmt.Sprintf("fb/users/%s.jpg", gonanoid.Must(8))
			cdnPath = fmt.Sprintf("https://assets.peatch.io/%s", imgFile)
			if err = s.uploadImageToS3(data.User.PhotoURL, imgFile); err != nil {
				return nil, err
			}
		} else {
			// get random one of 30 avatars
			cdnPath = fmt.Sprintf("https://assets.peatch.io/avatars/%d.svg", rand.Intn(30)+1)
		}

		create := db.User{
			FirstName:    firstName,
			LastName:     lastName,
			Username:     username,
			ChatID:       data.User.ID,
			AvatarURL:    &cdnPath,
			ReferralCode: GenerateReferralCode(6),
			ReferredBy:   nil,
		}

		lang := "ru"

		if data.User.LanguageCode != "ru" {
			lang = "en"
		}

		create.LanguageCode = &lang

		if err = s.storage.CreateUser(create); err != nil {
			return nil, err
		}

		user, err = s.storage.GetUserByChatID(data.User.ID)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	token, err := generateJWT(user.ID, user.ChatID)

	if err != nil {
		return nil, err
	}

	return &contract.UserAuthResponse{
		ID:                 user.ID,
		FirstName:          user.FirstName,
		LastName:           user.LastName,
		Username:           user.Username,
		LanguageCode:       user.LanguageCode,
		ChatID:             user.ChatID,
		CreatedAt:          user.CreatedAt,
		Token:              token,
		TotalPoints:        user.TotalPoints,
		TotalPredictions:   user.TotalPredictions,
		CorrectPredictions: user.CorrectPredictions,
		AvatarURL:          user.AvatarURL,
		GlobalRank:         user.GlobalRank,
	}, nil
}

type JWTClaims struct {
	jwt.RegisteredClaims
	UID    int   `json:"uid"`
	ChatID int64 `json:"chat_id"`
}

func generateJWT(id int, chatID int64) (string, error) {
	claims := &JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
		UID:    id,
		ChatID: chatID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte("secret"))
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
