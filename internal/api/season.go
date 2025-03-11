package api

import (
	"github.com/labstack/echo/v4"
	"github.com/user/project/internal/contract"
	"github.com/user/project/internal/terrors"
	"net/http"
)

func (a *API) GetActiveSeasons(c echo.Context) error {
	seasons, err := a.storage.GetActiveSeasons(c.Request().Context())
	if err != nil {
		return terrors.InternalServer(err, "failed to get active season")
	}

	var resp []contract.SeasonResponse
	for _, season := range seasons {
		resp = append(resp, contract.SeasonResponse{
			ID:        season.ID,
			Name:      season.Name,
			StartDate: season.StartDate,
			EndDate:   season.EndDate,
			Type:      season.Type,
		})
	}

	return c.JSON(http.StatusOK, resp)
}
