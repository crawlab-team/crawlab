package errors

import "errors"

var (
	ErrInvalidType   = errors.New("invalid type")
	ErrMissingValue  = errors.New("missing value")
	ErrNoCursor      = errors.New("no cursor")
	ErrAlreadyLocked = errors.New("already locked")
)
