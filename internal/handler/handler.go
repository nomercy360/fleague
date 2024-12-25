package handler

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/user/project/internal/contract"
	"github.com/user/project/internal/db"
	"github.com/user/project/internal/terrors"
	"net/http"
)

// servicer interface for database operations
type servicer interface {
	GetActiveMatches(ctx context.Context, uid string) ([]contract.MatchResponse, error)
	Health() (db.HealthStats, error)
	TelegramAuth(request contract.AuthTelegramRequest) (*contract.UserAuthResponse, error)
	SavePrediction(ctx context.Context, uid string, prediction contract.PredictionRequest) error
	GetUserPredictions(ctx context.Context, uid string) ([]contract.PredictionResponse, error)
	GetLeaderboard(ctx context.Context) ([]contract.LeaderboardEntry, error)
	GetUserInfo(ctx context.Context, username string) (*contract.UserInfoResponse, error)
	GetActiveSeason(ctx context.Context) (contract.SeasonResponse, error)
	GetUserReferrals(ctx context.Context, userID string) ([]contract.UserProfile, error)
}

// Handler struct for handling business logic
type Handler struct {
	service servicer
}

// New creates a new Handler instance
func New(service servicer) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Health(c echo.Context) error {
	stats, err := h.service.Health()
	if err != nil {
		return terrors.InternalServer(err, "cannot get health stats")
	}

	return c.JSON(http.StatusOK, stats)
}

func (h *Handler) AuthTelegram(c echo.Context) error {
	var req contract.AuthTelegramRequest
	if err := c.Bind(&req); err != nil {
		return terrors.BadRequest(err, "failed to bind request")
	}

	if err := req.Validate(); err != nil {
		return terrors.BadRequest(err, "failed to validate request")
	}

	user, err := h.service.TelegramAuth(req)
	if err != nil {
		return terrors.InternalServer(err, "failed to authenticate user")
	}

	return c.JSON(http.StatusOK, user)
}

func (h *Handler) GetActiveSeason(c echo.Context) error {
	resp, err := h.service.GetActiveSeason(c.Request().Context())
	if err != nil {
		return terrors.InternalServer(err, "cannot get active season")
	}

	return c.JSON(http.StatusOK, resp)
}
