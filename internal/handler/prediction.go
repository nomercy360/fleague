package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/user/project/internal/contract"
	"github.com/user/project/internal/terrors"
	"net/http"
)

func (h *Handler) SavePrediction(c echo.Context) error {
	var prediction contract.PredictionRequest
	if err := c.Bind(&prediction); err != nil {
		return terrors.BadRequest(err, "failed to decode request")
	}

	if err := prediction.Validate(); err != nil {
		return terrors.BadRequest(err, "failed to validate request")
	}

	err := h.service.SavePrediction(c.Request().Context(), getUserID(c), prediction)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) GetUserPredictions(c echo.Context) error {
	resp, err := h.service.GetUserPredictions(c.Request().Context(), getUserID(c))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}
