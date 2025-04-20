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
	PaymentID string    `json:"payment_id" db:"payment_id"`
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
		INSERT INTO subscriptions (id, user_id, start_date, end_date, is_paid, created_at, payment_id)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := s.db.ExecContext(ctx, query,
		sub.ID, sub.UserID, sub.StartDate, sub.EndDate, sub.IsPaid, sub.CreatedAt, sub.PaymentID,
	)
	return err
}

func (s *Storage) GetActiveSubscription(ctx context.Context, uid string) (Subscription, error) {
	query := `
		SELECT id, user_id, start_date, end_date, is_paid, created_at, payment_id
		FROM subscriptions
		WHERE user_id = ? AND end_date > CURRENT_TIMESTAMP`
	var sub Subscription
	err := s.db.QueryRowContext(ctx, query, uid).Scan(
		&sub.ID,
		&sub.UserID,
		&sub.StartDate,
		&sub.EndDate,
		&sub.IsPaid,
		&sub.CreatedAt,
		&sub.PaymentID,
	)

	if err != nil && IsNoRowsError(err) {
		return Subscription{}, ErrNotFound
	} else if err != nil {
		return Subscription{}, err
	}

	return sub, nil
}

func (s *Storage) SuspendSubscription(ctx context.Context, uid string) error {
	query := `
		UPDATE subscriptions 
		SET end_date = CURRENT_TIMESTAMP
		WHERE user_id = ? AND end_date > CURRENT_TIMESTAMP`
	_, err := s.db.ExecContext(ctx, query, uid)
	if err != nil {
		return err
	}

	query = `
		UPDATE users
		SET subscription_active = false, subscription_expiry = CURRENT_TIMESTAMP
		WHERE id = ?`

	_, err = s.db.ExecContext(ctx, query, uid)
	if err != nil {
		return err
	}

	return err
}
