package store

import (
	"errors"
	"fmt"
)

var ErrChanNotFound = errors.New("channel not found")
var ErrEmptyChanId = errors.New("empty channel id")
var ErrMsgNotFound = errors.New("message not found")
var ErrUserNotFound = errors.New("user not found")
var ErrUserAlreadyRegistered = errors.New("user already registered")
var ErrUserAlreadyJoined = errors.New("user already joined to channel")
var ErrUserAlreadyLeave = errors.New("user already leave channel")
var ErrMessageNotFound = errors.New("message not found")

type Error struct {
	Err        error
	StorageErr error
}

func (e Error) Unwrap() error {
	return e.Err
}

func (e Error) StorageError() error {
	return e.StorageErr
}

func (e Error) Error() string {
	return fmt.Sprintf("Store err: %v, storage err: %v", e.Err, e.StorageErr)
}

func StorageError(storeErr, storageErr error) Error {
	return Error{Err: storeErr, StorageErr: storageErr}
}
