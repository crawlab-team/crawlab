package constants

import (
	"errors"
)

var (
	//ErrorMongoError                = e.NewSystemOPError(1001, "system error:[mongo]%s", http.StatusInternalServerError)
	//ErrorUserNotFound              = e.NewBusinessError(10001, "user not found.", http.StatusUnauthorized)
	//ErrorUsernameOrPasswordInvalid = e.NewBusinessError(11001, "username or password invalid", http.StatusUnauthorized)
	ErrAlreadyExists    = errors.New("already exists")
	ErrNotExists        = errors.New("not exists")
	ErrForbidden        = errors.New("forbidden")
	ErrInvalidOperation = errors.New("invalid operation")
	ErrInvalidOptions   = errors.New("invalid options")
	ErrNoTasksAvailable = errors.New("no tasks available")
	ErrInvalidType      = errors.New("invalid type")
	ErrInvalidSignal    = errors.New("invalid signal")
	ErrEmptyValue       = errors.New("empty value")
	ErrTaskError        = errors.New("task error")
	ErrTaskLost         = errors.New("task lost")
	ErrTaskCancelled    = errors.New("task cancelled")
	ErrUnableToCancel   = errors.New("unable to cancel")
	ErrUnableToDispose  = errors.New("unable to dispose")
	ErrAlreadyDisposed  = errors.New("already disposed")
	ErrStopped          = errors.New("stopped")
	ErrMissingCol       = errors.New("missing col")
	ErrInvalidValue     = errors.New("invalid value")
	ErrInvalidCronSpec  = errors.New("invalid cron spec")
)
