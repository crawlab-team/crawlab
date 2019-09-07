package constants

import (
	"crawlab/errors"
	"net/http"
)

var (
	ErrorMongoError = errors.NewSystemOPError(1001, "system error:[mongo]%s", http.StatusInternalServerError)
	//users
	ErrorBadRequest                = errors.NewBusinessError(400, "bad request, please confirm parameters.", http.StatusBadRequest)
	ErrorUserNotFound              = errors.NewBusinessError(10001, "user not found.", http.StatusUnauthorized)
	ErrorUsernameOrPasswordInvalid = errors.NewBusinessError(11001, "username or password invalid.", http.StatusUnauthorized)
	ErrorAccountHasExisted         = errors.NewBusinessError(11002, "account has been existed.", http.StatusBadRequest)
	ErrorAccountDisabled           = errors.NewBusinessError(11003, "This account has been deactivated.", http.StatusUnauthorized)
	ErrorAccountNoPermission       = errors.NewBusinessError(11004, "This account no  permission.", http.StatusForbidden)
	ErrorTokenExpired              = errors.NewBusinessError(11005, "Token has been expired.", http.StatusUnauthorized)
	ErrorNeedResetPassword         = errors.NewBusinessError(11006, "Should be reset password", http.StatusUnauthorized)
)
