package db

import (
	"context"
	"time"
)

type Prediction struct {
	ID                 int       `json:"id" db:"id"`
	UserID             int       `json:"user_id" db:"user_id"`
	MatchID            int       `json:"match_id" db:"match_id"`
	PredictedOutcome   *string   `json:"predicted_outcome" db:"predicted_outcome"`
	PredictedHomeScore *int      `json:"predicted_home_score" db:"predicted_home_score"`
	PredictedAwayScore *int      `json:"predicted_away_score" db:"predicted_away_score"`
	PointsAwarded      int       `json:"points_awarded" db:"points_awarded"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
}

func (s *storage) SavePrediction(ctx context.Context, prediction Prediction) error {
	query := `
		INSERT INTO predictions (user_id, match_id, predicted_outcome, predicted_home_score, predicted_away_score)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(user_id, match_id) DO UPDATE SET
		predicted_outcome = excluded.predicted_outcome,
		predicted_home_score = excluded.predicted_home_score,
		predicted_away_score = excluded.predicted_away_score,
		points_awarded = excluded.points_awarded`
	_, err := s.db.ExecContext(ctx, query,
		prediction.UserID,
		prediction.MatchID,
		prediction.PredictedOutcome,
		prediction.PredictedHomeScore,
		prediction.PredictedAwayScore,
	)

	return err
}

func (s *storage) GetUserPredictionByMatchID(ctx context.Context, uid, matchID int) (*Prediction, error) {
	query := `
		SELECT
			id,
			user_id,
			match_id,
			predicted_outcome,
			predicted_home_score,
			predicted_away_score,
			points_awarded,
			created_at
		FROM predictions
		WHERE user_id = ? AND match_id = ?`

	var prediction Prediction
	err := s.db.QueryRowContext(ctx, query, uid, matchID).Scan(
		&prediction.ID,
		&prediction.UserID,
		&prediction.MatchID,
		&prediction.PredictedOutcome,
		&prediction.PredictedHomeScore,
		&prediction.PredictedAwayScore,
		&prediction.PointsAwarded,
		&prediction.CreatedAt,
	)

	if err != nil && IsNoRowsError(err) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return &prediction, nil
}
