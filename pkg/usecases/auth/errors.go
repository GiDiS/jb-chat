package auth

import "errors"

var (
	ErrAuthRequired       = errors.New("auth required")
	ErrUnknownAuthService = errors.New("unknown auth service")
)
