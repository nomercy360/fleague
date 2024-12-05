package db

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

// User represents a user in the system
type User struct {
	ID           int       `db:"id"`
	FirstName    *string   `db:"first_name"` // Optional field
	LastName     *string   `db:"last_name"`  // Optional field
	Username     string    `db:"username"`
	LanguageCode *string   `db:"language_code"` // Optional field
	ChatID       int64     `db:"chat_id"`
	CreatedAt    time.Time `db:"created_at"`
}

func IsNoRowsError(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func (s *storage) AddUserToLeague(ctx context.Context, leagueID int, userID string) error {
	query := `
        INSERT INTO league_members (league_id, user_id, joined_at)
        VALUES (?, ?, ?)`

	_, err := s.db.ExecContext(ctx, query, leagueID, userID, time.Now())
	return err
}

func (s *storage) CreateUser(user User) error {
	query := `
		INSERT INTO users (id, first_name, last_name, username, language_code, chat_id, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := s.db.Exec(query, user.ID, user.FirstName, user.LastName, user.Username, user.LanguageCode, user.ChatID, time.Now())
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
