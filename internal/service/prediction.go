package service

import (
	"context"
	"errors"
	"github.com/user/project/internal/contract"
	"github.com/user/project/internal/db"
	"github.com/user/project/internal/terrors"
	"sort"
)

func (s Service) SavePrediction(ctx context.Context, uid string, req contract.PredictionRequest) error {
	match, err := s.storage.GetMatchByID(ctx, req.MatchID)
	if err != nil && errors.Is(err, db.ErrNotFound) {
		return terrors.BadRequest(nil, "match not found")
	} else if err != nil {
		return err
	}

	if match.Status != db.MatchStatusScheduled {
		return terrors.BadRequest(nil, "match is not scheduled")
	}

	prediction := db.Prediction{
		UserID:             uid,
		MatchID:            req.MatchID,
		PredictedOutcome:   req.PredictedOutcome,
		PredictedHomeScore: req.PredictedHomeScore,
		PredictedAwayScore: req.PredictedAwayScore,
	}

	if err := s.storage.SavePrediction(ctx, prediction); err != nil {
		return err
	}

	if err := s.storage.UpdateUserPredictionCount(ctx, uid); err != nil {
		return err
	}

	return nil
}

func (s Service) predictionsByUserID(ctx context.Context, uid string, onlyCompleted bool) ([]contract.PredictionResponse, error) {
	predictions, err := s.storage.GetPredictionsByUserID(ctx, uid, onlyCompleted)
	if err != nil {
		return nil, err
	}

	var res []contract.PredictionResponse
	for _, prediction := range predictions {
		match, err := s.storage.GetMatchByID(ctx, prediction.MatchID)
		if err != nil && errors.Is(err, db.ErrNotFound) {
			return nil, terrors.NotFound(err, "match not found")
		} else if err != nil {
			return nil, err
		}

		homeTeam, err := s.storage.GetTeamByID(ctx, match.HomeTeamID)
		if err != nil {
			return nil, terrors.InternalServer(err, "failed to get home team")
		}

		awayTeam, err := s.storage.GetTeamByID(ctx, match.AwayTeamID)
		if err != nil {
			return nil, terrors.InternalServer(err, "failed to get away team")
		}

		res = append(res, contract.PredictionResponse{
			UserID:             prediction.UserID,
			MatchID:            prediction.MatchID,
			PredictedOutcome:   prediction.PredictedOutcome,
			PredictedHomeScore: prediction.PredictedHomeScore,
			PredictedAwayScore: prediction.PredictedAwayScore,
			PointsAwarded:      prediction.PointsAwarded,
			CreatedAt:          prediction.CreatedAt,
			CompletedAt:        prediction.CompletedAt,
			Match:              toMatchResponse(match, homeTeam, awayTeam),
		})
	}

	// sort predictions by match date
	sort.Slice(res, func(i, j int) bool {
		return res[i].Match.MatchDate.After(res[j].Match.MatchDate)
	})

	return res, nil
}

func (s Service) GetUserPredictions(ctx context.Context, uid string) ([]contract.PredictionResponse, error) {
	return s.predictionsByUserID(ctx, uid, false)
}
