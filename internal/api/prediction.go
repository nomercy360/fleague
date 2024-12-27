package api

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/user/project/internal/contract"
	"github.com/user/project/internal/db"
	"github.com/user/project/internal/terrors"
	"net/http"
)

func (a API) SavePrediction(c echo.Context) error {
	var req contract.PredictionRequest
	if err := c.Bind(&req); err != nil {
		return terrors.BadRequest(err, "failed to decode request")
	}

	if err := req.Validate(); err != nil {
		return terrors.BadRequest(err, "failed to validate request")
	}

	ctx := c.Request().Context()

	uid := getUserID(c)

	match, err := a.storage.GetMatchByID(ctx, req.MatchID)
	if err != nil && errors.Is(err, db.ErrNotFound) {
		return terrors.BadRequest(nil, "match not found")
	} else if err != nil {
		return err
	}

	if match.Status != db.MatchStatusScheduled {
		return terrors.BadRequest(nil, "match is not scheduled")
	}

	prediction := db.Prediction{
		UserID:             uid,
		MatchID:            req.MatchID,
		PredictedOutcome:   req.PredictedOutcome,
		PredictedHomeScore: req.PredictedHomeScore,
		PredictedAwayScore: req.PredictedAwayScore,
	}

	if err := a.storage.SavePrediction(ctx, prediction); err != nil {
		return err
	}

	if err := a.storage.UpdateUserPredictionCount(ctx, uid); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
}

func (a API) GetUserPredictions(c echo.Context) error {
	ctx := c.Request().Context()
	uid := getUserID(c)

	resp, err := a.predictionsByUserID(ctx, uid, false)

	if err != nil {
		return terrors.InternalServer(err, "failed to get user predictions")
	}

	return c.JSON(http.StatusOK, resp)
}
