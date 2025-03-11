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
	"math"
	"net/http"
	"path/filepath"
	"sort"
	"time"
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

func (a *API) predictionsByUserID(ctx context.Context, uid string, onlyCompleted bool) ([]contract.PredictionResponse, error) {
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

func (a *API) GetUserInfo(c echo.Context) error {
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

	if resp.User.TotalPredictions > 0 {
		accuracy := (float64(resp.User.CorrectPredictions) / float64(resp.User.TotalPredictions)) * 100
		resp.User.PredictionAccuracy = math.Round(accuracy*100) / 100
	}

	return c.JSON(http.StatusOK, resp)
}

func (a *API) ListMyReferrals(c echo.Context) error {
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
			TotalPredictions:   user.TotalPredictions,
			CorrectPredictions: user.CorrectPredictions,
		})
	}

	return c.JSON(http.StatusOK, users)
}

func (a *API) UpdateUser(c echo.Context) error {
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

	if req.AvatarURL != nil {
		user.AvatarURL = req.AvatarURL
	}

	if err := a.storage.UpdateUserInformation(ctx, user); err != nil {
		return terrors.InternalServer(err, "could not update user")
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "ok"})
}

func (a *API) GetPresignedURL(c echo.Context) error {
	var req contract.PresignedURLRequest
	if err := c.Bind(&req); err != nil {
		return terrors.BadRequest(err, "failed to bind request")
	}

	if err := c.Validate(req); err != nil {
		return terrors.BadRequest(err, "failed to validate request")
	}

	uid := GetContextUserID(c)

	if uid == "" {
		return terrors.Unauthorized(nil, "unauthorized")
	}

	fileExt := filepath.Ext(req.FileName)
	if fileExt == "" {
		fileExt = ".jpg" // Default extension
	}

	fileName := fmt.Sprintf("wishes/%d%s", time.Now().Unix(), fileExt)

	url, err := a.s3.GetPresignedURL(fileName, 15*time.Minute)

	if err != nil {
		return terrors.InternalServer(err, "failed to get presigned url")
	}

	res := contract.PresignedURLResponse{
		URL:      url,
		FileName: fileName,
		CdnURL:   fmt.Sprintf("%s/%s", a.cfg.AssetsURL, fileName),
	}

	return c.JSON(http.StatusOK, res)
}

func (a *API) FollowUserHandler(c echo.Context) error {
	followingID := c.Param("user_id")
	followerID := GetContextUserID(c)

	if followerID == "" {
		return terrors.Unauthorized(nil, "unauthorized")
	}
	if followerID == followingID {
		return terrors.BadRequest(nil, "you cannot follow yourself")
	}

	ctx := c.Request().Context()
	err := a.storage.FollowUser(ctx, followerID, followingID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return terrors.NotFound(err, "user not found")
		}
		return terrors.InternalServer(err, "failed to follow user")
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "followed successfully"})
}

// UnfollowUserHandler allows a user to unfollow another user
func (a *API) UnfollowUserHandler(c echo.Context) error {
	followingID := c.Param("user_id")
	followerID := GetContextUserID(c)

	if followerID == "" {
		return terrors.Unauthorized(nil, "unauthorized")
	}

	ctx := c.Request().Context()
	err := a.storage.UnfollowUser(ctx, followerID, followingID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return terrors.NotFound(err, "user not found")
		}
		return terrors.InternalServer(err, "failed to unfollow user")
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "unfollowed successfully"})
}

func (a *API) GetFollowersHandler(c echo.Context) error {
	userID := c.Param("user_id")
	ctx := c.Request().Context()

	followers, err := a.storage.GetFollowers(ctx, userID)
	if err != nil {
		return terrors.InternalServer(err, "failed to get followers")
	}

	var users []contract.UserProfile
	for _, user := range followers {
		users = append(users, contract.UserProfile{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Username:  user.Username,
			AvatarURL: user.AvatarURL,
		})
	}

	return c.JSON(http.StatusOK, users)
}

// GetFollowingHandler retrieves a list of users the logged-in user is following
func (a *API) GetFollowingHandler(c echo.Context) error {
	userID := GetContextUserID(c)
	if userID == "" {
		return terrors.Unauthorized(nil, "unauthorized")
	}

	ctx := c.Request().Context()
	following, err := a.storage.GetFollowing(ctx, userID)
	if err != nil {
		return terrors.InternalServer(err, "failed to get following list")
	}

	var users []contract.UserProfile
	for _, user := range following {
		users = append(users, contract.UserProfile{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Username:  user.Username,
			AvatarURL: user.AvatarURL,
		})
	}

	return c.JSON(http.StatusOK, users)
}
