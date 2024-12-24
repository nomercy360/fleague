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
	ID                 string    `db:"id"`
	FirstName          *string   `db:"first_name"`
	LastName           *string   `db:"last_name"`
	Username           string    `db:"username"`
	AvatarURL          *string   `db:"avatar_url"`
	LanguageCode       *string   `db:"language_code"`
	ChatID             int64     `db:"chat_id"`
	ReferralCode       string    `db:"referral_code"`
	ReferredBy         *int      `db:"referred_by"`
	CreatedAt          time.Time `db:"created_at"`
	TotalPoints        int       `db:"total_points"`
	TotalPredictions   int       `db:"total_predictions"`
	CorrectPredictions int       `db:"correct_predictions"`
	GlobalRank         int       `db:"global_rank"`
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
		INSERT INTO users (id, first_name, last_name, username, language_code, chat_id, avatar_url, referral_code)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.db.Exec(query, user.ID, user.FirstName, user.LastName, user.Username, user.LanguageCode, user.ChatID, user.AvatarURL, user.ReferralCode)
	return err
}

func (s *storage) getUserBy(query string, args ...interface{}) (*User, error) {
	var user User
	row := s.db.QueryRowContext(context.Background(), query, args...)

	if err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.LanguageCode,
		&user.ChatID,
		&user.CreatedAt,
		&user.TotalPoints,
		&user.TotalPredictions,
		&user.CorrectPredictions,
		&user.AvatarURL,
		&user.ReferredBy,
		&user.ReferralCode,
		&user.GlobalRank,
	); err != nil && IsNoRowsError(err) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *storage) GetUserByChatID(chatID int64) (*User, error) {
	query := `
		SELECT id, first_name, last_name, username, language_code, chat_id, created_at, total_points, total_predictions, correct_predictions, avatar_url, referred_by, referral_code, global_rank
		FROM (
		         SELECT id, first_name, last_name, username, language_code, chat_id, created_at, total_points, total_predictions, correct_predictions, avatar_url, referred_by, referral_code,
		                RANK() OVER (ORDER BY total_points DESC) AS global_rank
		         FROM users
		     ) ranked_users
		WHERE chat_id = ?`

	return s.getUserBy(query, chatID)
}

func (s *storage) GetUserByID(id string) (*User, error) {
	query := `
		SELECT id, first_name, last_name, username, language_code, chat_id, created_at, total_points, total_predictions, correct_predictions, avatar_url, referred_by, referral_code, global_rank
		FROM (
		         SELECT id, first_name, last_name, username, language_code, chat_id, created_at, total_points, total_predictions, correct_predictions, avatar_url, referred_by, referral_code,
		                RANK() OVER (ORDER BY total_points DESC) AS global_rank
		         FROM users
		     ) ranked_users
		WHERE id = ?`

	return s.getUserBy(query, id)
}

func (s *storage) GetUserByUsername(uname string) (*User, error) {
	query := `
		SELECT id, first_name, last_name, username, language_code, chat_id, created_at, total_points, total_predictions, correct_predictions, avatar_url, referred_by, referral_code, global_rank
		FROM (
		         SELECT id, first_name, last_name, username, language_code, chat_id, created_at, total_points, total_predictions, correct_predictions, avatar_url, referred_by, referral_code,
		                RANK() OVER (ORDER BY total_points DESC) AS global_rank
		         FROM users
		     ) ranked_users
		WHERE username = ?`

	return s.getUserBy(query, uname)
}

func (s *storage) UpdateUserPoints(ctx context.Context, userID string, points int) error {
	query := `
		UPDATE users
		SET total_points = total_points + ?,
		    correct_predictions = correct_predictions + 1
		WHERE id = ?`

	_, err := s.db.ExecContext(ctx, query, points, userID)
	return err
}

func (s *storage) UpdateUserPredictionCount(ctx context.Context, userID string) error {
	query := `
		UPDATE users
		SET total_predictions = total_predictions + 1
		WHERE id = ?`

	_, err := s.db.ExecContext(ctx, query, userID)
	return err
}

func (s *storage) ListUserReferrals(ctx context.Context, userID string) ([]User, error) {
	query := `
		SELECT id, first_name, last_name, username, language_code, chat_id, created_at, total_points, total_predictions, correct_predictions, avatar_url, referred_by, global_rank
		FROM (
		         SELECT id, first_name, last_name, username, language_code, chat_id, created_at, total_points, total_predictions, correct_predictions, avatar_url, referred_by,
		                RANK() OVER (ORDER BY total_points DESC) AS global_rank
		         FROM users
		     ) ranked_users
		WHERE referred_by = ?`

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Username,
			&user.LanguageCode,
			&user.ChatID,
			&user.CreatedAt,
			&user.TotalPoints,
			&user.TotalPredictions,
			&user.CorrectPredictions,
			&user.AvatarURL,
			&user.ReferredBy,
			&user.GlobalRank,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
