package service

import (
	"context"
	"errors"
	"github.com/user/project/internal/contract"
	"github.com/user/project/internal/db"
)

func (s Service) SavePrediction(ctx context.Context, req contract.PredictionRequest) error {
	uid := GetUserIDFromContext(ctx)

	match, err := s.storage.GetMatchByID(ctx, req.MatchID)
	if err != nil && errors.Is(err, db.ErrNotFound) {
		return contract.ErrMatchNotFound
	} else if err != nil {
		return err
	}

	if match.Status != "scheduled" {
		return contract.MatchStatusInvalid
	}

	prediction := db.Prediction{
		UserID:             uid,
		MatchID:            req.MatchID,
		PredictedOutcome:   req.PredictedOutcome,
		PredictedHomeScore: req.PredictedHomeScore,
		PredictedAwayScore: req.PredictedAwayScore,
	}

	return s.storage.SavePrediction(ctx, prediction)
}
