package constants

import (
	"crawlab/errors"
	"net/http"
)

var (
	ErrorMongoError = errors.NewSystemOPError(1001, "system error:[mongo]%s", http.StatusInternalServerError)
	//users
	ErrorUserNotFound              = errors.NewBusinessError(10001, "user not found.", http.StatusUnauthorized)
	ErrorUsernameOrPasswordInvalid = errors.NewBusinessError(11001, "username or password invalid", http.StatusUnauthorized)
)
