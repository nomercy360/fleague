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

	seasons, err := a.storage.GetActiveSeasons(ctx)
	if err != nil {
		return terrors.InternalServer(err, "failed to get active seasons")
	}

	var monthlySeason, footballSeason *db.Season
	for _, season := range seasons {
		if season.Type == "monthly" {
			monthlySeason = &season
		} else if season.Type == "football" {
			footballSeason = &season
		}
	}

	getLeaderboardForSeason := func(season *db.Season) ([]contract.LeaderboardEntry, error) {
		if season == nil {
			return nil, nil
		}

		res, err := a.storage.GetLeaderboard(ctx, season.ID)
		if err != nil {
			return nil, terrors.InternalServer(err, "failed to get leaderboard")
		}

		var leaderboard []contract.LeaderboardEntry
		for _, entry := range res {
			user, err := a.storage.GetUserByID(entry.UserID)
			if err != nil && !errors.Is(err, db.ErrNotFound) {
				return nil, terrors.InternalServer(err, "failed to get user")
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

		// Sort leaderboard by points
		sort.Slice(leaderboard, func(i, j int) bool {
			return leaderboard[i].Points > leaderboard[j].Points
		})

		return leaderboard, nil
	}

	monthlyLeaderboard, err := getLeaderboardForSeason(monthlySeason)
	if err != nil {
		return err
	}

	footballLeaderboard, err := getLeaderboardForSeason(footballSeason)
	if err != nil {
		return err
	}

	response := map[string]interface{}{
		"monthly":  monthlyLeaderboard,
		"football": footballLeaderboard,
	}

	return c.JSON(http.StatusOK, response)
}
