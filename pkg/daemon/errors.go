package daemon

const (
	ErrOk = iota
	ErrInitLoggerFailed
	ErrInitHttpServeFailed
	ErrInitInterruptsFailed
)
