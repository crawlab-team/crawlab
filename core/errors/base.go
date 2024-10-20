package errors

import (
	"errors"
	"fmt"
)

const (
	ErrorPrefixModel = "model"
	ErrorPrefixGrpc  = "grpc"
	ErrorPrefixNode  = "node"
	ErrorPrefixTask  = "task"
	ErrorPrefixUser  = "user"
)

type ErrorPrefix string

func NewError(prefix ErrorPrefix, msg string) (err error) {
	return errors.New(fmt.Sprintf("%s error: %s", prefix, msg))
}
