package handler

import (
	"github.com/go-chi/render"
	"github.com/user/project/internal/handler/errrender"
	"net/http"
)

func (h *Handler) ListMatches(w http.ResponseWriter, r *http.Request) {
	resp, err := h.service.GetActiveMatches(r.Context(), nil)
	if err != nil {
		errrender.RenderError(w, r, err)
		return
	}

	render.JSON(w, r, resp)
}
