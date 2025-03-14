package api

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/user/project/internal/contract"
	"github.com/user/project/internal/db"
	"github.com/user/project/internal/terrors"
	"net/http"
)

var ErrInsufficientTokens = terrors.Forbidden(errors.New("insufficient tokens"), "insufficient tokens")

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

	match, err := a.storage.GetMatchByID(ctx, req.MatchID)
	if err != nil && errors.Is(err, db.ErrNotFound) {
		return terrors.BadRequest(nil, "match not found")
	} else if err != nil {
		return err
	}

	if match.Status != db.MatchStatusScheduled {
		return terrors.BadRequest(nil, "match is not scheduled")
	}

	// Check if this is an update
	existing, err := a.storage.GetUserPredictionByMatchID(ctx, uid, req.MatchID)
	if err != nil && !errors.Is(err, db.ErrNotFound) {
		return terrors.InternalServer(err, "failed to check existing prediction")
	}

	isUpdate := false
	if err == nil && existing.UserID == uid {
		isUpdate = true
	}

	newTokenCost := 0
	if req.PredictedHomeScore != nil && req.PredictedAwayScore != nil {
		newTokenCost = 10 // Exact score
	} else if req.PredictedOutcome != nil {
		newTokenCost = 20 // Outcome
	}

	var resultBalance int

	user, err := a.storage.GetUserByID(uid)
	if err != nil {
		return terrors.InternalServer(err, "failed to get user")
	}

	// For new predictions, set token cost and deduct
	if !isUpdate {
		if newTokenCost > 0 {
			if user.PredictionTokens < newTokenCost {
				return ErrInsufficientTokens
			}

			resultBalance, err = a.storage.UpdateUserTokens(ctx, uid, -newTokenCost, db.TokenTransactionTypePrediction)
			if err != nil {
				return terrors.InternalServer(err, "failed to update user tokens")
			}
		}
	} else {
		// For updates, calculate token difference
		oldTokenCost := existing.TokenCost
		tokenDiff := newTokenCost - oldTokenCost

		switch {
		case tokenDiff > 0: // Additional tokens required
			if user.PredictionTokens < tokenDiff {
				return ErrInsufficientTokens
			}

			resultBalance, err = a.storage.UpdateUserTokens(ctx, uid, -tokenDiff, db.TokenTransactionTypePrediction)
			if err != nil {
				return terrors.InternalServer(err, "failed to update user tokens")
			}

		case tokenDiff < 0: // Refund excess tokens
			resultBalance, err = a.storage.UpdateUserTokens(ctx, uid, -tokenDiff, db.TokenTransactionTypePredictionRefund)
			if err != nil {
				return terrors.InternalServer(err, "failed to update user tokens")
			}

		default: // No change
			resultBalance = user.PredictionTokens
		}
	}

	// Save prediction
	prediction := db.Prediction{
		UserID:             uid,
		MatchID:            req.MatchID,
		PredictedOutcome:   req.PredictedOutcome,
		PredictedHomeScore: req.PredictedHomeScore,
		PredictedAwayScore: req.PredictedAwayScore,
		TokenCost:          newTokenCost,
	}

	if err := a.storage.SavePrediction(ctx, prediction); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"status": "ok", "balance": resultBalance})
}

func (a *API) CancelPrediction(c echo.Context) error {
	ctx := c.Request().Context()
	uid := GetContextUserID(c)

	matchID := c.Param("id")
	if matchID == "" {
		return terrors.BadRequest(nil, "match_id is required")
	}

	// Check match status
	match, err := a.storage.GetMatchByID(ctx, matchID)
	if err != nil && errors.Is(err, db.ErrNotFound) {
		return terrors.BadRequest(nil, "match not found")
	} else if err != nil {
		return err
	}
	if match.Status != db.MatchStatusScheduled {
		return terrors.BadRequest(nil, "cannot cancel prediction for a match that has started or completed")
	}

	// Check if prediction exists
	prediction, err := a.storage.GetUserPredictionByMatchID(ctx, uid, matchID)
	if err != nil && errors.Is(err, db.ErrNotFound) {
		return terrors.BadRequest(nil, "no prediction found for this match")
	} else if err != nil {
		return err
	}

	// Delete prediction and refund tokens
	if err := a.storage.DeletePrediction(ctx, uid, matchID); err != nil {
		return err
	}

	balance, err := a.storage.UpdateUserTokens(
		ctx,
		uid,
		prediction.TokenCost,
		db.TokenTransactionTypePredictionRefund,
	)

	if err != nil {
		return terrors.InternalServer(err, "failed to update user tokens")
	}

	return c.JSON(http.StatusOK, echo.Map{"status": "ok", "balance": balance})
}

func (a *API) GetUserPredictions(c echo.Context) error {
	ctx := c.Request().Context()
	uid := GetContextUserID(c)

	resp, err := a.predictionsByUserID(ctx, uid, false)

	if err != nil {
		return terrors.InternalServer(err, "failed to get user predictions")
	}

	return c.JSON(http.StatusOK, resp)
}
