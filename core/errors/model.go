package errors

func NewModelError(msg string) (err error) {
	return NewError(ErrorPrefixModel, msg)
}

var ErrorModelInvalidType = NewModelError("invalid type")
var ErrorModelMissingRequiredData = NewModelError("missing required data")
