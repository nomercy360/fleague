package db

import (
	"context"
	"time"
)

type Prediction struct {
	UserID             string     `json:"user_id" db:"user_id"`
	MatchID            string     `json:"match_id" db:"match_id"`
	PredictedOutcome   *string    `json:"predicted_outcome" db:"predicted_outcome"`
	PredictedHomeScore *int       `json:"predicted_home_score" db:"predicted_home_score"`
	PredictedAwayScore *int       `json:"predicted_away_score" db:"predicted_away_score"`
	PointsAwarded      int        `json:"points_awarded" db:"points_awarded"`
	CompletedAt        *time.Time `json:"completed_at" db:"completed_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
}

const (
	MatchOutcomeHome = "home"
	MatchOutcomeAway = "away"
	MatchOutcomeDraw = "draw"
)

func (s *Storage) SavePrediction(ctx context.Context, prediction Prediction) error {
	query := `
		INSERT INTO predictions (
			user_id, match_id, predicted_outcome, predicted_home_score, predicted_away_score
		) VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(user_id, match_id) DO UPDATE SET
			predicted_outcome = excluded.predicted_outcome,
			predicted_home_score = excluded.predicted_home_score,
			predicted_away_score = excluded.predicted_away_score,
			updated_at = CURRENT_TIMESTAMP`
	_, err := s.db.ExecContext(ctx, query,
		prediction.UserID,
		prediction.MatchID,
		prediction.PredictedOutcome,
		prediction.PredictedHomeScore,
		prediction.PredictedAwayScore,
	)
	return err
}

func (s *Storage) DeletePrediction(ctx context.Context, userID, matchID string) error {
	query := `DELETE FROM predictions WHERE user_id = ? AND match_id = ?`
	_, err := s.db.ExecContext(ctx, query, userID, matchID)
	return err
}

func (s *Storage) GetUserPredictionByMatchID(ctx context.Context, uid, matchID string) (Prediction, error) {
	query := `
		SELECT
			user_id,
			match_id,
			predicted_outcome,
			predicted_home_score,
			predicted_away_score,
			points_awarded,
			created_at,
			updated_at,
			completed_at
		FROM predictions
		WHERE user_id = ? AND match_id = ?`

	var prediction Prediction
	err := s.db.QueryRowContext(ctx, query, uid, matchID).Scan(
		&prediction.UserID,
		&prediction.MatchID,
		&prediction.PredictedOutcome,
		&prediction.PredictedHomeScore,
		&prediction.PredictedAwayScore,
		&prediction.PointsAwarded,
		&prediction.CreatedAt,
		&prediction.UpdatedAt,
		&prediction.CompletedAt,
	)

	if err != nil && IsNoRowsError(err) {
		return Prediction{}, ErrNotFound
	} else if err != nil {
		return Prediction{}, err
	}

	return prediction, nil
}

func (s *Storage) GetPredictionsByUserID(ctx context.Context, uid string, opts ...PredictionFilter) ([]Prediction, error) {
	query := `
		SELECT
			p.user_id,
			p.match_id,
			p.predicted_outcome,
			p.predicted_home_score,
			p.predicted_away_score,
			p.points_awarded,
			p.created_at,
			p.updated_at,
			p.completed_at
		FROM predictions p
		JOIN matches m ON p.match_id = m.id
		WHERE user_id = ?`
	args := []interface{}{uid}

	filters := predictionFilters{}
	for _, opt := range opts {
		opt(&filters)
	}
	if filters.OnlyCompleted {
		query += " AND completed_at IS NOT NULL"
	}
	if !filters.StartTime.IsZero() {
		query += " AND m.match_date >= ?"
		args = append(args, filters.StartTime)
	}
	if !filters.EndTime.IsZero() {
		query += " AND m.match_date < ?"
		args = append(args, filters.EndTime)
	}
	query += " ORDER BY m.match_date ASC"
	if filters.Limit > 0 {
		query += " LIMIT ?"
		args = append(args, filters.Limit)
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	predictions := make([]Prediction, 0)
	for rows.Next() {
		var p Prediction
		err := rows.Scan(
			&p.UserID,
			&p.MatchID,
			&p.PredictedOutcome,
			&p.PredictedHomeScore,
			&p.PredictedAwayScore,
			&p.PointsAwarded,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.CompletedAt,
		)
		if err != nil {
			return nil, err
		}
		predictions = append(predictions, p)
	}
	return predictions, nil
}

type PredictionFilter func(*predictionFilters)

type predictionFilters struct {
	OnlyCompleted bool
	StartTime     time.Time
	EndTime       time.Time
	Limit         int
}

func WithOnlyCompleted() PredictionFilter {
	return func(f *predictionFilters) { f.OnlyCompleted = true }
}

func WithStartTime(start time.Time) PredictionFilter {
	return func(f *predictionFilters) { f.StartTime = start }
}

func WithEndTime(end time.Time) PredictionFilter {
	return func(f *predictionFilters) { f.EndTime = end }
}

func WithLimit(limit int) PredictionFilter {
	return func(f *predictionFilters) { f.Limit = limit }
}

func (s *Storage) GetPredictionsForMatch(ctx context.Context, matchID string) ([]Prediction, error) {
	query := `
		SELECT
			user_id,
			match_id,
			predicted_outcome,
			predicted_home_score,
			predicted_away_score,
			points_awarded,
			created_at,
			updated_at,
			completed_at
		FROM predictions
		WHERE match_id = ?`

	rows, err := s.db.QueryContext(ctx, query, matchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var predictions []Prediction
	for rows.Next() {
		var prediction Prediction
		err := rows.Scan(
			&prediction.UserID,
			&prediction.MatchID,
			&prediction.PredictedOutcome,
			&prediction.PredictedHomeScore,
			&prediction.PredictedAwayScore,
			&prediction.PointsAwarded,
			&prediction.CreatedAt,
			&prediction.UpdatedAt,
			&prediction.CompletedAt,
		)
		if err != nil {
			return nil, err
		}

		predictions = append(predictions, prediction)
	}

	return predictions, nil
}

func (s *Storage) UpdatePredictionResult(ctx context.Context, matchID, userID string, points int) error {
	query := `
		UPDATE predictions
		SET points_awarded = ?, completed_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE match_id = ? AND user_id = ?`
	_, err := s.db.ExecContext(ctx, query, points, matchID, userID)

	return err
}
