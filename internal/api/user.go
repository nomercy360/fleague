package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/user/project/internal/contract"
	"github.com/user/project/internal/db"
	"github.com/user/project/internal/terrors"
	"net/http"
	"sort"
)

func getUserID(c echo.Context) string {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*contract.JWTClaims)
	return claims.UID
}

func (a API) ListMatches(c echo.Context) error {
	ctx := c.Request().Context()
	res, err := a.storage.GetActiveMatches(ctx)
	uid := getUserID(c)

	if err != nil {
		return terrors.InternalServer(err, "failed to get active matches")
	}

	var matches []contract.MatchResponse
	for _, match := range res {
		homeTeam, err := a.storage.GetTeamByID(ctx, match.HomeTeamID)
		if err != nil && errors.Is(err, db.ErrNotFound) {
			return terrors.NotFound(err, fmt.Sprintf("team with id %s not found", match.HomeTeamID))
		} else if err != nil {
			return terrors.InternalServer(err, "failed to get home team")
		}

		awayTeam, err := a.storage.GetTeamByID(ctx, match.AwayTeamID)
		if err != nil && errors.Is(err, db.ErrNotFound) {
			return terrors.NotFound(err, fmt.Sprintf("team with id %s not found", match.AwayTeamID))
		} else if err != nil {
			return terrors.InternalServer(err, "failed to get away team")
		}

		prediction, err := a.storage.GetUserPredictionByMatchID(ctx, uid, match.ID)
		if err != nil && !errors.Is(err, db.ErrNotFound) {
			return terrors.InternalServer(err, "failed to get user prediction")
		}

		resp := toMatchResponse(match, homeTeam, awayTeam)

		if prediction.UserID != "" {
			resp.Prediction = &prediction
		}

		matches = append(matches, resp)
	}

	// sort matches by date
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].MatchDate.Before(matches[j].MatchDate)
	})

	return c.JSON(http.StatusOK, matches)
}

func toMatchResponse(match db.Match, homeTeam db.Team, awayTeam db.Team) contract.MatchResponse {
	return contract.MatchResponse{
		ID:         match.ID,
		Tournament: match.Tournament,
		MatchDate:  match.MatchDate,
		Status:     match.Status,
		HomeTeam:   homeTeam,
		AwayTeam:   awayTeam,
		HomeScore:  match.HomeScore,
		AwayScore:  match.AwayScore,
	}
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

		homeTeam, err := a.storage.GetTeamByID(ctx, match.HomeTeamID)
		if err != nil {
			return nil, terrors.InternalServer(err, "failed to get home team")
		}

		awayTeam, err := a.storage.GetTeamByID(ctx, match.AwayTeamID)
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

	rank, err := a.storage.GetUserRank(ctx, user.ID)
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
			GlobalRank:         rank,
			FavoriteTeam:       user.FavoriteTeam,
			CurrentWinStreak:   user.CurrentWinStreak,
			LongestWinStreak:   user.LongestWinStreak,
		},
		Predictions: userPredictions,
	}

	return c.JSON(http.StatusOK, resp)
}

func (a API) ListMyReferrals(c echo.Context) error {
	res, err := a.storage.ListUserReferrals(c.Request().Context(), getUserID(c))
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

	uid := getUserID(c)
	ctx := c.Request().Context()

	user, err := a.storage.GetUserByID(uid)

	if err != nil && errors.Is(err, db.ErrNotFound) {
		return terrors.NotFound(err, "user not found")
	}

	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.FavoriteTeamID = req.FavoriteTeamID

	if err := a.storage.UpdateUserInformation(ctx, user); err != nil {
		return terrors.InternalServer(err, "could not update user")
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "ok"})
}
