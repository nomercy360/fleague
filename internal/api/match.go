package api

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/user/project/internal/contract"
	"github.com/user/project/internal/db"
	"github.com/user/project/internal/terrors"
	"net/http"
)

func (a API) GetMatchByID(c echo.Context) error {
	matchID := c.Param("id")
	ctx := c.Request().Context()

	match, err := a.storage.GetMatchByID(ctx, matchID)
	if err != nil && errors.Is(err, db.ErrNotFound) {
		return terrors.NotFound(err, "match not found")
	} else if err != nil {
		return terrors.InternalServer(err, "failed to get match")
	}

	stats, err := a.storage.GetPredictionStats(ctx, matchID)
	if err != nil {
		return terrors.InternalServer(err, "failed to get prediction stats")
	}

	response := toMatchResponse(match)
	response.PredictionStats = stats

	return c.JSON(http.StatusOK, toMatchResponse(match))
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

func toMatchResponse(match db.Match) contract.MatchResponse {
	return contract.MatchResponse{
		ID:         match.ID,
		Tournament: match.Tournament,
		HomeTeam:   match.HomeTeam,
		AwayTeam:   match.AwayTeam,
		MatchDate:  match.MatchDate,
		Status:     match.Status,
		AwayScore:  match.AwayScore,
		HomeScore:  match.HomeScore,
		Prediction: match.Prediction,
		HomeOdds:   match.HomeOdds,
		DrawOdds:   match.DrawOdds,
		AwayOdds:   match.AwayOdds,
	}
}
