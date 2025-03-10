package db

import (
	"context"
	"time"
)

type Survey struct {
	ID         string    `json:"id" db:"id"`
	UserID     string    `json:"user_id" db:"user_id"`
	Feature    string    `json:"feature" db:"feature"`
	Preference string    `json:"preference" db:"preference"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

func (s *Storage) SaveSurvey(ctx context.Context, survey Survey) error {
	query := `
        INSERT INTO surveys (id, user_id, feature, preference)
        VALUES (?, ?, ?, ?)
        ON CONFLICT DO NOTHING` // Prevents duplicate submissions if needed

	_, err := s.db.ExecContext(ctx, query,
		survey.ID,
		survey.UserID,
		survey.Feature,
		survey.Preference,
	)
	return err
}

func (s *Storage) GetSurveyByUserAndFeature(ctx context.Context, userID, feature string) (Survey, error) {
	query := `
        SELECT id, user_id, feature, preference, created_at
        FROM surveys
        WHERE user_id = ? AND feature = ?`

	var survey Survey
	err := s.db.QueryRowContext(ctx, query, userID, feature).Scan(
		&survey.ID,
		&survey.UserID,
		&survey.Feature,
		&survey.Preference,
		&survey.CreatedAt,
	)

	if err != nil && IsNoRowsError(err) {
		return Survey{}, ErrNotFound
	} else if err != nil {
		return Survey{}, err
	}

	return survey, nil
}

func (s *Storage) GetSurveyStats(ctx context.Context, feature string) (map[string]int, error) {
	query := `
        SELECT preference, COUNT(*) as count
        FROM surveys
        WHERE feature = ?
        GROUP BY preference`

	rows, err := s.db.QueryContext(ctx, query, feature)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make(map[string]int)
	for rows.Next() {
		var preference string
		var count int
		if err := rows.Scan(&preference, &count); err != nil {
			return nil, err
		}
		stats[preference] = count
	}

	return stats, rows.Err()
}
