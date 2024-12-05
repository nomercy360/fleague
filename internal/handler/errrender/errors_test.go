package errrender

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_renderError(t *testing.T) {
	tests := []struct {
		name  string
		error error
		msg   string
		code  int
	}{
		{
			name:  "decode json",
			error: contract.ErrDecodeJSON,
			msg:   contract.FailedDecodeJSON,
			code:  http.StatusBadRequest,
		},
		{
			name:  "invalid session token",
			error: contract.ErrInvalidSessionToken,
			msg:   contract.InvalidSessionToken,
			code:  http.StatusBadRequest,
		},
		{
			name:  "insufficient balance",
			error: contract.ErrInsufficientBalance,
			msg:   contract.InsufficientBalance,
			code:  http.StatusForbidden,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var (
				w    = httptest.NewRecorder()
				r    = httptest.NewRequest(http.MethodPost, "/", nil)
				resp = fmt.Sprintf("{\"message\":\"%s\"}", test.msg)
			)

			RenderError(w, r, test.error)

			assert.JSONEq(t, resp, w.Body.String())
			assert.Equal(t, test.code, w.Code)
		})
	}
}
