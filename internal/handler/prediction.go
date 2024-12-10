package handler

import (
	"encoding/json"
	"github.com/go-chi/render"
	"github.com/user/project/internal/contract"
	"github.com/user/project/internal/handler/errrender"
	"net/http"
)

func (h *Handler) SavePrediction(w http.ResponseWriter, r *http.Request) {
	var prediction contract.PredictionRequest
	err := json.NewDecoder(r.Body).Decode(&prediction)
	if err != nil {
		errrender.RenderError(w, r, err, contract.FailedDecodeJSON)
		return
	}

	if err := prediction.Validate(); err != nil {
		errrender.RenderError(w, r, err, contract.InvalidRequest)
		return
	}

	err = h.service.SavePrediction(r.Context(), prediction)
	if err != nil {
		errrender.RenderError(w, r, err, "failed to save prediction")
		return
	}

	render.NoContent(w, r)
}
