package errors

import (
	"errors"
	"fmt"
)

var (
	ErrorRedisInvalidType = NewRedisError("invalid type")
	ErrorRedisLocked      = NewRedisError("locked")
)

func NewRedisError(msg string) (err error) {
	return errors.New(fmt.Sprintf("%s: %s", errorPrefixRedis, msg))
}
