package db

import (
	"context"
	"time"
)

type Subscription struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	StartDate time.Time `json:"start_date" db:"start_date"`
	EndDate   time.Time `json:"end_date" db:"end_date"`
	IsPaid    bool      `json:"is_paid" db:"is_paid"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func (s *Storage) UpdateUserSubscription(ctx context.Context, uid string, active bool, expiry time.Time) error {
	query := `
		UPDATE users 
		SET subscription_active = ?, subscription_expiry = ?
		WHERE id = ?`
	_, err := s.db.ExecContext(ctx, query, active, expiry, uid)
	return err
}

func (s *Storage) SaveSubscription(ctx context.Context, sub Subscription) error {
	query := `
		INSERT INTO subscriptions (id, user_id, start_date, end_date, is_paid, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`
	_, err := s.db.ExecContext(ctx, query,
		sub.ID, sub.UserID, sub.StartDate, sub.EndDate, sub.IsPaid, sub.CreatedAt,
	)
	return err
}
