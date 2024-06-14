package errors

import "errors"

func NewModelError(msg string) (err error) {
	return NewError(ErrorPrefixModel, msg)
}

var ErrorModelInvalidType = NewModelError("invalid type")
var ErrorModelInvalidModelId = NewModelError("invalid model id")
var ErrorModelNotImplemented = NewModelError("not implemented")
var ErrorModelNotFound = NewModelError("not found")
var ErrorModelAlreadyExists = NewModelError("already exists")
var ErrorModelNotExists = NewModelError("not exists")
var ErrorModelMissingRequiredData = NewModelError("missing required data")
var ErrorModelMissingId = errors.New("missing _id")
var ErrorModelNotAllowed = NewModelError("not allowed")
var ErrorModelDeleteListError = NewModelError("delete list error")
var ErrorModelNilPointer = NewModelError("nil pointer")
