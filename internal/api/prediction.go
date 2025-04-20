package api

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/user/project/internal/contract"
	"github.com/user/project/internal/db"
	"github.com/user/project/internal/terrors"
	"net/http"
	"time"
)

var ErrNoActiveSubscription = terrors.Forbidden(errors.New("no active subscription"), "no active subscription")

func (a *API) SavePrediction(c echo.Context) error {
	var req contract.PredictionRequest
	if err := c.Bind(&req); err != nil {
		return terrors.BadRequest(err, "failed to decode request")
	}
	if err := req.Validate(); err != nil {
		return terrors.BadRequest(err, "failed to validate request")
	}

	ctx := c.Request().Context()
	uid := GetContextUserID(c)

	//user, err := a.storage.GetUserByID(uid)
	//if err != nil {
	//	return terrors.InternalServer(err, "failed to get user")
	//}
	//
	//if !user.SubscriptionActive || user.SubscriptionExpiry.Before(time.Now()) {
	//	return ErrNoActiveSubscription
	//}

	match, err := a.storage.GetMatchByID(ctx, req.MatchID)
	if err != nil && errors.Is(err, db.ErrNotFound) {
		return terrors.BadRequest(nil, "match not found")
	} else if err != nil {
		return err
	}
	if match.Status != db.MatchStatusScheduled {
		return terrors.BadRequest(nil, "match is not scheduled")
	}

	// Сохраняем прогноз без учета токенов
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

	return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
}

func (a *API) CancelPrediction(c echo.Context) error {
	ctx := c.Request().Context()
	uid := GetContextUserID(c)

	matchID := c.Param("id")
	if matchID == "" {
		return terrors.BadRequest(nil, "match_id is required")
	}

	// Проверка статуса матча
	match, err := a.storage.GetMatchByID(ctx, matchID)
	if err != nil && errors.Is(err, db.ErrNotFound) {
		return terrors.BadRequest(nil, "match not found")
	} else if err != nil {
		return err
	}
	if match.Status != db.MatchStatusScheduled {
		return terrors.BadRequest(nil, "cannot cancel prediction for a match that has started or completed")
	}

	// Проверка существования прогноза
	_, err = a.storage.GetUserPredictionByMatchID(ctx, uid, matchID)
	if err != nil && errors.Is(err, db.ErrNotFound) {
		return terrors.BadRequest(nil, "no prediction found for this match")
	} else if err != nil {
		return err
	}

	// Удаляем прогноз без возврата токенов
	if err := a.storage.DeletePrediction(ctx, uid, matchID); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
}

func (a *API) GetUserPredictions(c echo.Context) error {
	ctx := c.Request().Context()
	uid := GetContextUserID(c)

	resp, err := a.predictionsByUserID(
		ctx,
		uid,
		db.WithStartTime(time.Now().Add(-7*24*time.Hour)),
		db.WithLimit(100),
	)

	if err != nil {
		return terrors.InternalServer(err, "failed to get user predictions")
	}
	return c.JSON(http.StatusOK, resp)
}
