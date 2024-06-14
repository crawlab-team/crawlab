package errors

func NewEventError(msg string) (err error) {
	return NewError(ErrorPrefixEvent, msg)
}

var ErrorEventNotFound = NewEventError("not found")
var ErrorEventInvalidType = NewEventError("invalid type")
var ErrorEventAlreadyExists = NewEventError("already exists")
var ErrorEventUnknownAction = NewEventError("unknown action")
