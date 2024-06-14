package errors

func NewFsError(msg string) (err error) {
	return NewError(ErrorPrefixFs, msg)
}

var ErrorFsForbidden = NewFsError("forbidden")
var ErrorFsEmptyWorkspacePath = NewFsError("empty workspace path")
var ErrorFsInvalidType = NewFsError("invalid type")
var ErrorFsAlreadyExists = NewFsError("already exists")
var ErrorFsInvalidContent = NewFsError("invalid content")
