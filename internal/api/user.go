package api

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/user/project/internal/contract"
	"github.com/user/project/internal/db"
	"github.com/user/project/internal/terrors"
	"net/http"
	"sort"
)

func GetContextUserID(c echo.Context) string {
	user, ok := c.Get("user").(*jwt.Token)
	if !ok || user == nil {
		return "" // or handle the error appropriately
	}

	claims, ok := user.Claims.(*contract.JWTClaims)
	if !ok || claims == nil {
		return "" // or handle the error appropriately
	}

	return claims.UID
}

func (a API) predictionsByUserID(ctx context.Context, uid string, onlyCompleted bool) ([]contract.PredictionResponse, error) {
	predictions, err := a.storage.GetPredictionsByUserID(ctx, uid, onlyCompleted)
	if err != nil {
		return nil, err
	}

	var res []contract.PredictionResponse
	for _, prediction := range predictions {
		match, err := a.storage.GetMatchByID(ctx, prediction.MatchID)
		if err != nil && errors.Is(err, db.ErrNotFound) {
			return nil, terrors.NotFound(err, "match not found")
		} else if err != nil {
			return nil, err
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
			Match:              toMatchResponse(match),
		})
	}

	// sort predictions by match date
	sort.Slice(res, func(i, j int) bool {
		return res[i].Match.MatchDate.After(res[j].Match.MatchDate)
	})

	return res, nil
}

func (a API) GetUserInfo(c echo.Context) error {
	username := c.Param("username")
	ctx := c.Request().Context()

	user, err := a.storage.GetUserByUsername(username)

	if err != nil && errors.Is(err, db.ErrNotFound) {
		return terrors.NotFound(err, "user not found")
	} else if err != nil {
		return terrors.InternalServer(err, "failed to get user")
	}

	userPredictions, err := a.predictionsByUserID(ctx, user.ID, false)

	if err != nil {
		return terrors.InternalServer(err, "failed to get user predictions")
	}

	ranks, err := a.storage.GetUserRank(ctx, user.ID)
	if err != nil {
		return terrors.InternalServer(err, "failed to get user rank")
	}

	resp := &contract.UserInfoResponse{
		User: contract.UserProfile{
			ID:                 user.ID,
			FirstName:          user.FirstName,
			LastName:           user.LastName,
			Username:           user.Username,
			AvatarURL:          user.AvatarURL,
			TotalPoints:        user.TotalPoints,
			TotalPredictions:   user.TotalPredictions,
			CorrectPredictions: user.CorrectPredictions,
			Ranks:              ranks,
			FavoriteTeam:       user.FavoriteTeam,
			CurrentWinStreak:   user.CurrentWinStreak,
			LongestWinStreak:   user.LongestWinStreak,
			Badges:             user.Badges,
		},
		Predictions: userPredictions,
	}

	return c.JSON(http.StatusOK, resp)
}

func (a API) ListMyReferrals(c echo.Context) error {
	res, err := a.storage.ListUserReferrals(c.Request().Context(), GetContextUserID(c))
	if err != nil {
		return terrors.InternalServer(err, "failed to get user referrals")
	}

	var users []contract.UserProfile
	for _, user := range res {
		users = append(users, contract.UserProfile{
			ID:                 user.ID,
			FirstName:          user.FirstName,
			LastName:           user.LastName,
			Username:           user.Username,
			AvatarURL:          user.AvatarURL,
			TotalPoints:        user.TotalPoints,
			TotalPredictions:   user.TotalPredictions,
			CorrectPredictions: user.CorrectPredictions,
		})
	}

	return c.JSON(http.StatusOK, users)
}

func (a API) UpdateUser(c echo.Context) error {
	var req contract.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return terrors.BadRequest(err, "failed to decode request")
	}

	if err := req.Validate(); err != nil {
		return terrors.BadRequest(err, "failed to validate request")
	}

	uid := GetContextUserID(c)
	ctx := c.Request().Context()

	user, err := a.storage.GetUserByID(uid)

	if err != nil && errors.Is(err, db.ErrNotFound) {
		return terrors.NotFound(err, "user not found")
	}

	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.FavoriteTeamID = req.FavoriteTeamID
	if req.LanguageCode != nil {
		user.LanguageCode = req.LanguageCode
	}

	if err := a.storage.UpdateUserInformation(ctx, user); err != nil {
		return terrors.InternalServer(err, "could not update user")
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "ok"})
}
