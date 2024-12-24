package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/user/project/internal/terrors"
	"net/http"
)

func (h *Handler) GetLeaderboard(c echo.Context) error {
	resp, err := h.service.GetLeaderboard(c.Request().Context())

	if err != nil {
		return terrors.InternalServer(err, "cannot get leaderboard")
	}

	return c.JSON(http.StatusOK, resp)
}
