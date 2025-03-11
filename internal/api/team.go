package api

import (
	"github.com/labstack/echo/v4"
	"github.com/user/project/internal/terrors"
	"net/http"
)

func (a *API) ListTeams(c echo.Context) error {
	teams, err := a.storage.ListTeams(c.Request().Context())
	if err != nil {
		return terrors.InternalServer(err, "failed to list teams")
	}

	return c.JSON(http.StatusOK, teams)
}
