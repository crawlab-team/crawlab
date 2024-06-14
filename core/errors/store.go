package errors

func NewInjectError(msg string) (err error) {
	return NewError(ErrorPrefixInject, msg)
}

var ErrorInjectEmptyValue = NewInjectError("empty value")
var ErrorInjectNotExists = NewInjectError("not exists")
var ErrorInjectInvalidType = NewInjectError("invalid type")
