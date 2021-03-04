package daemon

import "errors"

const (
	ErrOk = iota
	ErrInitLoggerFailed
	ErrInitHttpServeFailed
)

var (
	ErrInvalidRequest     = errors.New("invalid request")
	ErrAuthRequired       = errors.New("auth required")
	ErrUnknownAuthService = errors.New("unknown auth service")
)
