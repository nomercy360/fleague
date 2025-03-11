package api

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/user/project/internal/contract"
	"github.com/user/project/internal/db"
	"github.com/user/project/internal/nanoid"
	"github.com/user/project/internal/terrors"
	"net/http"
)

func (a *API) SaveSurvey(c echo.Context) error {
	var req contract.SurveyRequest
	if err := c.Bind(&req); err != nil {
		return terrors.BadRequest(err, "failed to decode request")
	}

	if req.Feature == "" || (req.Preference != "yes" && req.Preference != "no") {
		return terrors.BadRequest(nil, "invalid feature or preference")
	}

	ctx := c.Request().Context()
	uid := GetContextUserID(c)

	_, err := a.storage.GetSurveyByUserAndFeature(ctx, uid, req.Feature)
	if err == nil {
		return terrors.BadRequest(nil, "user already responded to this survey")
	} else if !errors.Is(err, db.ErrNotFound) {
		return err
	}

	survey := db.Survey{
		ID:         nanoid.Must(),
		UserID:     uid,
		Feature:    req.Feature,
		Preference: req.Preference,
	}

	if err := a.storage.SaveSurvey(ctx, survey); err != nil {
		return terrors.InternalServer(err, "failed to save survey response")
	}

	return c.JSON(http.StatusOK, echo.Map{"status": "ok"})
}

func (a *API) GetSurveyStats(c echo.Context) error {
	feature := c.QueryParam("feature")
	if feature == "" {
		return terrors.BadRequest(nil, "feature parameter is required")
	}

	ctx := c.Request().Context()
	stats, err := a.storage.GetSurveyStats(ctx, feature)
	if err != nil {
		return terrors.InternalServer(err, "failed to get survey stats")
	}

	return c.JSON(http.StatusOK, stats)
}
