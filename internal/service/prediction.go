package service

import (
	"context"
	"errors"
	"github.com/user/project/internal/contract"
	"github.com/user/project/internal/db"
	"sort"
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

func (s Service) GetUserPredictions(ctx context.Context) ([]contract.PredictionResponse, error) {
	uid := GetUserIDFromContext(ctx)

	predictions, err := s.storage.GetPredictionsByUserID(ctx, uid)
	if err != nil {
		return nil, err
	}

	var res []contract.PredictionResponse
	for _, prediction := range predictions {
		match, err := s.storage.GetMatchByID(ctx, prediction.MatchID)
		if err != nil && errors.Is(err, db.ErrNotFound) {
			return nil, contract.ErrMatchNotFound
		} else if err != nil {
			return nil, err
		}

		homeTeam, err := s.storage.GetTeamByID(ctx, match.HomeTeamID)
		if err != nil {
			return nil, err
		}

		awayTeam, err := s.storage.GetTeamByID(ctx, match.AwayTeamID)
		if err != nil {
			return nil, err
		}

		res = append(res, contract.PredictionResponse{
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

	// sort predictions by creation date
	sort.Slice(res, func(i, j int) bool {
		return res[i].CreatedAt.After(res[j].CreatedAt)
	})

	return res, nil
}
