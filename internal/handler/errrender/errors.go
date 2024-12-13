package errrender

import (
	"errors"
	"github.com/go-chi/render"
	"github.com/user/project/internal/contract"
	"log"
	"net/http"
)

func RenderError(w http.ResponseWriter, r *http.Request, err error, msg string) {
	if errors.Is(err, contract.ErrDecodeJSON) {
		RenderJSON(w, r, http.StatusBadRequest, err, msg)
		return
	}

	if errors.Is(err, contract.ErrInvalidRequest) {
		RenderJSON(w, r, http.StatusBadRequest, err, msg)
		return
	}

	if errors.Is(err, contract.ErrMatchNotFound) {
		RenderJSON(w, r, http.StatusNotFound, err, msg)
		return
	}

	if errors.Is(err, contract.MatchStatusInvalid) {
		RenderJSON(w, r, http.StatusConflict, err, msg)
		return
	}

	if errors.Is(err, contract.ErrUserNotFound) {
		RenderJSON(w, r, http.StatusNotFound, err, msg)
		return
	}

	RenderJSON(w, r, http.StatusInternalServerError, err, msg)
}

func RenderJSON(w http.ResponseWriter, r *http.Request, status int, err error, msg string) {
	if err != nil {
		log.Printf("error: %s: %v", msg, err)
	}

	render.Status(r, status)
	render.JSON(w, r, contract.Error{Message: err.Error()})
}
