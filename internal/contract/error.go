package contract

import "errors"

var (
	ErrDecodeJSON          = errors.New("decode json")
	ErrInvalidSessionToken = errors.New("invalid session token")
	ErrInvalidRequest      = errors.New("invalid request")
	ErrUnauthorized        = errors.New("unauthorized")
)

const (
	FailedDecodeJSON    = "Failed to decode json"
	InvalidSessionToken = "Invalid session token"
	InvalidRequest      = "Invalid request"
	Unauthorized        = "Unauthorized"
)
