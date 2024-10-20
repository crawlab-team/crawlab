package constants

import (
	"errors"
)

var (
	ErrNotExists      = errors.New("not exists")
	ErrInvalidOptions = errors.New("invalid options")
	ErrInvalidSignal  = errors.New("invalid signal")
	ErrTaskError      = errors.New("task error")
	ErrTaskLost       = errors.New("task lost")
	ErrTaskCancelled  = errors.New("task cancelled")
)
