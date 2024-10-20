package errors

func NewNodeError(msg string) (err error) {
	return NewError(ErrorPrefixNode, msg)
}

var ErrorNodeUnregistered = NewNodeError("unregistered")
var ErrorNodeNotExists = NewNodeError("not exists")
