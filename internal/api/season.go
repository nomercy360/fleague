package api

import (
	"github.com/labstack/echo/v4"
	"github.com/user/project/internal/contract"
	"github.com/user/project/internal/terrors"
	"net/http"
)

func (a API) GetActiveSeason(c echo.Context) error {
	season, err := a.storage.GetActiveSeason(c.Request().Context())
	if err != nil {
		return terrors.InternalServer(err, "failed to get active season")
	}

	resp := contract.SeasonResponse{
		ID:        season.ID,
		Name:      season.Name,
		IsActive:  season.IsActive,
		StartDate: season.StartDate,
		EndDate:   season.EndDate,
	}

	return c.JSON(http.StatusOK, resp)
}
