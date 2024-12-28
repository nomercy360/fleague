package api

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/user/project/internal/contract"
	"github.com/user/project/internal/db"
	"github.com/user/project/internal/terrors"
	"net/http"
	"sort"
)

func (a API) GetLeaderboard(c echo.Context) error {
	ctx := c.Request().Context()

	season, err := a.storage.GetActiveSeason(ctx)
	if err != nil {
		return terrors.InternalServer(err, "failed to get active season")
	}

	res, err := a.storage.GetLeaderboard(ctx, season.ID)

	if err != nil {
		return terrors.InternalServer(err, "failed to get leaderboard")
	}

	var leaderboard []contract.LeaderboardEntry
	for _, entry := range res {
		user, err := a.storage.GetUserByID(entry.UserID)
		if err != nil && !errors.Is(err, db.ErrNotFound) {
			return terrors.InternalServer(err, "failed to get user")
		} else if err != nil {
			continue
		}

		userProfile := contract.UserProfile{
			ID:               user.ID,
			FirstName:        user.FirstName,
			LastName:         user.LastName,
			Username:         user.Username,
			AvatarURL:        user.AvatarURL,
			FavoriteTeam:     user.FavoriteTeam,
			CurrentWinStreak: user.CurrentWinStreak,
			LongestWinStreak: user.LongestWinStreak,
		}

		leaderboard = append(leaderboard, contract.LeaderboardEntry{
			User:     userProfile,
			UserID:   entry.UserID,
			Points:   entry.Points,
			SeasonID: entry.SeasonID,
		})
	}

	// sort leaderboard by points
	sort.Slice(leaderboard, func(i, j int) bool {
		return leaderboard[i].Points > leaderboard[j].Points
	})

	return c.JSON(http.StatusOK, leaderboard)
}
