package errors

func NewFilterError(msg string) (err error) {
	return NewError(ErrorPrefixFilter, msg)
}

var ErrorFilterInvalidOperation = NewFilterError("invalid operation")
var ErrorFilterUnableToParseQuery = NewFilterError("unable to parse query")
