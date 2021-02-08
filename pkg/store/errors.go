package store

import "errors"

var ErrChanNotFound = errors.New("channel not found")
var ErrEmptyChanId = errors.New("empty channel id")
var ErrMsgNotFound = errors.New("message not found")
var ErrUserNotFound = errors.New("user not found")
var ErrUserAlreadyRegistered = errors.New("user already registered")
var ErrUserAlreadyJoined = errors.New("user already joined to channel")
var ErrUserAlreadyLeave = errors.New("user already leave channel")
var ErrMessageNotFound = errors.New("message not found")
