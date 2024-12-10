package db

import (
	"context"
	"database/sql"
	"errors"
	"github.com/mattn/go-sqlite3"
	"time"
)

// User represents a user in the system
type User struct {
	ID           int       `db:"id"`
	FirstName    *string   `db:"first_name"`
	LastName     *string   `db:"last_name"`
	Username     string    `db:"username"`
	LanguageCode *string   `db:"language_code"`
	ChatID       int64     `db:"chat_id"`
	ReferralCode string    `db:"referral_code"`
	ReferredBy   *int      `db:"referred_by"`
	CreatedAt    time.Time `db:"created_at"`
}

func IsNoRowsError(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func IsUniqueViolationError(err error) bool {
	var sqliteErr sqlite3.Error
	if errors.As(err, &sqliteErr) {
		return errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique)
	}
	return false
}

func IsForeignKeyViolationError(err error) bool {
	var sqliteErr sqlite3.Error
	if errors.As(err, &sqliteErr) {
		return errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintForeignKey)
	}
	return false
}

func (s *storage) CreateUser(user User) error {
	query := `
		INSERT INTO users (first_name, last_name, username, language_code, chat_id, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`

	_, err := s.db.Exec(query, user.FirstName, user.LastName, user.Username, user.LanguageCode, user.ChatID, time.Now())
	return err
}

func (s *storage) GetUserByChatID(chatID int64) (*User, error) {
	query := `
		SELECT id, first_name, last_name, username, language_code, chat_id, created_at
		FROM users
		WHERE chat_id = ?`

	var user User
	row := s.db.QueryRowContext(context.Background(), query, chatID)

	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Username, &user.LanguageCode, &user.ChatID, &user.CreatedAt)

	if err != nil && IsNoRowsError(err) {
		return nil, ErrNotFound
	}

	return &user, err
}
