package errors

func NewResultError(msg string) (err error) {
	return NewError(ErrorPrefixResult, msg)
}
