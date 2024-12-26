package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) GetLeaderboard(c echo.Context) error {
	resp, err := h.service.GetLeaderboard(c.Request().Context())

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}
