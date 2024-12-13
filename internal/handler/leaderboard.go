package handler

import (
	"github.com/go-chi/render"
	"github.com/user/project/internal/handler/errrender"
	"net/http"
)

func (h *Handler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	resp, err := h.service.GetLeaderboard(r.Context())

	if err != nil {
		errrender.RenderError(w, r, err, "failed to get leaderboard")
		return
	}

	render.JSON(w, r, resp)
}
