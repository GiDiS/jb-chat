package utils

import (
	"fmt"
)

const ErrUnknown = 0
const ErrBeginFailed = 1
const ErrCommitFailed = 2
const ErrRollbackFailed = 3
const ErrContext = 4
const ErrCallback = 5

type Error struct {
	DbErr error
	Code  int
}

func newDbErr(err error) error {
	dbErr := Error{DbErr: err, Code: ErrBeginFailed}
	return fmt.Errorf("db error: %w", dbErr)
}

func newBeginErr(err error) error {
	dbErr := Error{DbErr: err, Code: ErrBeginFailed}
	return fmt.Errorf("transaction begin failed: %w", dbErr)
}

func newCommitErr(err error) error {
	dbErr := Error{DbErr: err, Code: ErrCommitFailed}
	return fmt.Errorf("transaction commit failed: %w", dbErr)
}

func newRollbackErr(err error) error {
	dbErr := Error{DbErr: err, Code: ErrRollbackFailed}
	return fmt.Errorf("transaction rollback failed: %w", dbErr)
}

func newContextErr(err error) error {
	dbErr := Error{DbErr: err, Code: ErrContext}
	return fmt.Errorf("transaction context failed: %w", dbErr)
}

func newCallbackErr(err error) error {
	dbErr := Error{DbErr: err, Code: ErrCallback}
	return fmt.Errorf("callback failed: %w", dbErr)
}

func (ce Error) Error() string {
	return fmt.Sprint(ce.DbErr.Error())
}

func (ce Error) Unwrap() error {
	return ce.DbErr
}
