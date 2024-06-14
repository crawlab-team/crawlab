package errors

func NewHttpError(msg string) (err error) {
	return NewError(ErrorPrefixHttp, msg)
}

var ErrorHttpBadRequest = NewHttpError("bad request")
var ErrorHttpUnauthorized = NewHttpError("unauthorized")
var ErrorHttpNotFound = NewHttpError("not found")
