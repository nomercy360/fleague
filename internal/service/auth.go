package service

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	initdata "github.com/telegram-mini-apps/init-data-golang"
	"github.com/user/project/internal/contract"
	"github.com/user/project/internal/db"
	"time"
)

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

		create := db.User{
			FirstName: firstName,
			LastName:  lastName,
			Username:  username,
			ChatID:    data.User.ID,
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
		ID:           user.ID,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Username:     user.Username,
		LanguageCode: user.LanguageCode,
		ChatID:       user.ChatID,
		CreatedAt:    user.CreatedAt,
		Token:        token,
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
