package errors

func NewUserError(msg string) (err error) {
	return NewError(ErrorPrefixUser, msg)
}

var (
	ErrorUserInvalidType           = NewUserError("invalid type")
	ErrorUserInvalidToken          = NewUserError("invalid token")
	ErrorUserNotExists             = NewUserError("not exists")
	ErrorUserNotExistsInContext    = NewUserError("not exists in context")
	ErrorUserAlreadyExists         = NewUserError("already exists")
	ErrorUserMismatch              = NewUserError("mismatch")
	ErrorUserMissingRequiredFields = NewUserError("missing required fields")
	ErrorUserUnauthorized          = NewUserError("unauthorized")
	ErrorUserInvalidPassword       = NewUserError("invalid password (length must be no less than 5)")
)
