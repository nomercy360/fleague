package handler

import (
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
