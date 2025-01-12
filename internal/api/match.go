package api

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/user/project/internal/contract"
	"github.com/user/project/internal/db"
	"github.com/user/project/internal/terrors"
	"net/http"
)

func calculateWinLoss(matches []db.Match, teamID string) []string {
	var results []string
	for _, match := range matches {
		if match.HomeTeamID == teamID {
			if *match.HomeScore > *match.AwayScore {
				results = append(results, "W")
			} else {
				results = append(results, "L")
			}
		} else if match.AwayTeamID == teamID {
			if *match.AwayScore > *match.HomeScore {
				results = append(results, "W")
			} else {
				results = append(results, "L")
			}
		}
	}
	return results
}

func (a API) GetMatchByID(c echo.Context) error {
	matchID := c.Param("id")
	ctx := c.Request().Context()

	match, err := a.storage.GetMatchByID(ctx, matchID)
	if err != nil && errors.Is(err, db.ErrNotFound) {
		return terrors.NotFound(err, "match not found")
	} else if err != nil {
		return terrors.InternalServer(err, "failed to get match")
	}

	homeLastMatches, err := a.storage.GetLastMatchesByTeamID(ctx, match.HomeTeamID, 5)
	if err != nil {
		return fmt.Errorf("failed to fetch last matches for home team: %w", err)
	}

	awayLastMatches, err := a.storage.GetLastMatchesByTeamID(ctx, match.AwayTeamID, 5)
	if err != nil {
		return fmt.Errorf("failed to fetch last matches for away team: %w", err)
	}

	homeResults, awayResults := calculateWinLoss(homeLastMatches, match.HomeTeamID), calculateWinLoss(awayLastMatches, match.AwayTeamID)

	return c.JSON(http.StatusOK, toMatchResponse(match, homeResults, awayResults))
}

func (a API) ListMatches(c echo.Context) error {
	ctx := c.Request().Context()
	uid := getUserID(c)
	matches, err := a.storage.GetActiveMatches(ctx, uid)

	if err != nil {
		return terrors.InternalServer(err, "failed to get active matches")
	}

	return c.JSON(http.StatusOK, matches)
}

func toMatchResponse(match db.Match, homeResults, awayResults []string) contract.MatchResponse {
	return contract.MatchResponse{
		ID:              match.ID,
		Tournament:      match.Tournament,
		HomeTeam:        match.HomeTeam,
		AwayTeam:        match.AwayTeam,
		MatchDate:       match.MatchDate,
		Status:          match.Status,
		AwayScore:       match.AwayScore,
		HomeScore:       match.HomeScore,
		Prediction:      match.Prediction,
		HomeOdds:        match.HomeOdds,
		DrawOdds:        match.DrawOdds,
		AwayOdds:        match.AwayOdds,
		HomeTeamResults: homeResults,
		AwayTeamResults: awayResults,
	}
}
