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
		INSERT INTO predictions (user_id, match_id, predicted_outcome, predicted_home_score, predicted_away_score)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(user_id, match_id) DO UPDATE SET
		predicted_outcome = excluded.predicted_outcome,
		predicted_home_score = excluded.predicted_home_score,
		predicted_away_score = excluded.predicted_away_score,
		points_awarded = excluded.points_awarded,
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

func (s *Storage) GetPredictionsByUserID(ctx context.Context, uid string, onlyCompleted bool) ([]Prediction, error) {
	// Base query
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
		WHERE user_id = ?
	`

	if onlyCompleted {
		query += " AND completed_at IS NOT NULL"
	}

	rows, err := s.db.QueryContext(ctx, query, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var predictions []Prediction
	for rows.Next() {
		var prediction Prediction
		if err := rows.Scan(
			&prediction.UserID,
			&prediction.MatchID,
			&prediction.PredictedOutcome,
			&prediction.PredictedHomeScore,
			&prediction.PredictedAwayScore,
			&prediction.PointsAwarded,
			&prediction.CreatedAt,
			&prediction.UpdatedAt,
			&prediction.CompletedAt,
		); err != nil {
			return nil, err
		}

		predictions = append(predictions, prediction)
	}

	return predictions, rows.Err()
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
