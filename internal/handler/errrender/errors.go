package errrender

import (
	"errors"
	"github.com/go-chi/render"
	"github.com/user/project/internal/contract"
	"log"
	"net/http"
)

func RenderError(w http.ResponseWriter, r *http.Request, err error) {
	if errors.Is(err, contract.ErrDecodeJSON) {
		RenderJSON(w, r, http.StatusBadRequest, err, contract.FailedDecodeJSON)
		return
	}

	if errors.Is(err, contract.ErrInvalidSessionToken) {
		RenderJSON(w, r, http.StatusBadRequest, err, contract.InvalidSessionToken)
		return
	}

	if errors.Is(err, contract.ErrInvalidRequest) {
		RenderJSON(w, r, http.StatusBadRequest, err, contract.InvalidRequest)
		return
	}

	if errors.Is(err, contract.ErrUnauthorized) {
		RenderJSON(w, r, http.StatusUnauthorized, err, contract.Unauthorized)
		return
	}

	RenderJSON(w, r, http.StatusInternalServerError, err, "internal error")
}

func RenderJSON(w http.ResponseWriter, r *http.Request, status int, err error, msg string) {
	if err != nil {
		log.Println(err)
	}

	render.Status(r, status)
	render.JSON(w, r, contract.Error{Message: msg})
}
