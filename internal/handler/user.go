package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/user/project/internal/handler/errrender"
	"net/http"
)

func (h *Handler) ListMatches(w http.ResponseWriter, r *http.Request) {
	resp, err := h.service.GetActiveMatches(r.Context())
	if err != nil {
		errrender.RenderError(w, r, err, "failed to get active matches")
		return
	}

	render.JSON(w, r, resp)
}

func (h *Handler) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	resp, err := h.service.GetUserInfo(r.Context(), username)
	if err != nil {
		errrender.RenderError(w, r, err, "failed to get user info")
		return
	}

	render.JSON(w, r, resp)
}
