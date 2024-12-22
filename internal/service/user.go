package service

import (
	"context"
	"errors"
	"github.com/user/project/internal/contract"
	"github.com/user/project/internal/db"
)

func (s Service) GetUserInfo(ctx context.Context, username string) (*contract.UserInfoResponse, error) {
	user, err := s.storage.GetUserByUsername(username)

	if err != nil && errors.Is(err, db.ErrNotFound) {
		return nil, contract.ErrUserNotFound
	} else if err != nil {

	}

	userPredictions, err := s.predictionsByUserID(ctx, user.ID, true)

	if err != nil {
		return nil, err
	}

	return &contract.UserInfoResponse{
		User: contract.UserProfile{
			ID:                 user.ID,
			FirstName:          user.FirstName,
			LastName:           user.LastName,
			Username:           user.Username,
			AvatarURL:          user.AvatarURL,
			TotalPoints:        user.TotalPoints,
			TotalPredictions:   user.TotalPredictions,
			CorrectPredictions: user.CorrectPredictions,
			GlobalRank:         user.GlobalRank,
		},
		Predictions: userPredictions,
	}, nil
}
