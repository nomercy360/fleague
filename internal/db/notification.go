package db

import (
	"context"
	"github.com/user/project/internal/nanoid"
)

const (
	notificationTypeWeeklySummary = "weekly_summary"
	notificationTypeMatchReminder = "match_reminder"
)

func (s *Storage) LogNotification(ctx context.Context, userID, notificationType, relatedID string) error {
	query := `
        INSERT INTO notifications (id, user_id, notification_type, related_id)
        VALUES (?, ?, ?, ?)
    `
	_, err := s.db.ExecContext(ctx, query, nanoid.Must(), userID, notificationType, relatedID)
	return err
}

func (s *Storage) HasNotificationBeenSent(ctx context.Context, userID, notificationType, relatedID string) (bool, error) {
	query := `
        SELECT COUNT(*)
        FROM notifications
        WHERE user_id = ? AND notification_type = ? AND related_id = ?
    `
	var count int
	err := s.db.QueryRowContext(ctx, query, userID, notificationType, relatedID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

type WeeklyRecap struct {
	TotalPredictions int
	Wins             int
	Losses           int
	CurrentStreak    int
}

func (s *Storage) GetWeeklyRecap(ctx context.Context, userID string) (WeeklyRecap, error) {
	query := `
        SELECT
            COUNT(*) AS total_predictions,
            COALESCE(SUM(CASE WHEN p.points_awarded > 0 THEN 1 ELSE 0 END), 0) AS wins,
            COALESCE(SUM(CASE WHEN p.points_awarded = 0 THEN 1 ELSE 0 END), 0) AS losses,
            COALESCE(u.current_win_streak, 0) AS current_streak
        FROM predictions p
        JOIN users u ON p.user_id = u.id
        WHERE p.user_id = ? AND p.created_at >= DATE('now', '-7 days')
    `
	var recap WeeklyRecap
	err := s.db.QueryRowContext(ctx, query, userID).Scan(
		&recap.TotalPredictions,
		&recap.Wins,
		&recap.Losses,
		&recap.CurrentStreak,
	)
	if err != nil {
		return WeeklyRecap{}, err
	}
	return recap, nil
}
