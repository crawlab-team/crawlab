package constants

import "crawlab/errors"

var (
	//users
	ErrorUserNotFound = errors.NewBusinessError(10001, "user not found.")
)
