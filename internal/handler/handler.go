package handler

import (
	"context"
	"github.com/go-chi/render"
	"github.com/user/project/internal/contract"
	"github.com/user/project/internal/db"
	"github.com/user/project/internal/handler/errrender"
	"net/http"
)

// servicer interface for database operations
type servicer interface {
	GetActiveMatches(ctx context.Context) ([]contract.MatchResponse, error)
	Health() (db.HealthStats, error)
	TelegramAuth(query string) (*contract.UserAuthResponse, error)
	SavePrediction(ctx context.Context, prediction contract.PredictionRequest) error
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

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	stats, err := h.service.Health()
	if err != nil {
		errrender.RenderError(w, r, err, "failed to get health stats")
	}

	render.JSON(w, r, stats)
}

func (h *Handler) AuthTelegram(w http.ResponseWriter, r *http.Request) {
	query := r.URL.RawQuery
	user, err := h.service.TelegramAuth(query)
	if err != nil {
		errrender.RenderError(w, r, err, "failed to authenticate user")
		return
	}

	render.JSON(w, r, user)
}
