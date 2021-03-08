package usecases

import "errors"

var (
	ErrInvalidRequest     = errors.New("invalid request")
	ErrAuthRequired       = errors.New("auth required")
	ErrUnknownAuthService = errors.New("unknown auth service")
)
