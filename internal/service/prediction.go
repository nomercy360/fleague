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

	if match.Status != db.MatchStatusScheduled {
		return contract.MatchStatusInvalid
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

func (s Service) predictionsByUserID(ctx context.Context, uid int, onlyCompleted bool) ([]contract.PredictionResponse, error) {
	predictions, err := s.storage.GetPredictionsByUserID(ctx, uid, onlyCompleted)
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

	// sort predictions by match date
	sort.Slice(res, func(i, j int) bool {
		return res[i].Match.MatchDate.After(res[j].Match.MatchDate)
	})

	return res, nil
}

func (s Service) GetUserPredictions(ctx context.Context) ([]contract.PredictionResponse, error) {
	uid := GetUserIDFromContext(ctx)

	return s.predictionsByUserID(ctx, uid, false)
}
