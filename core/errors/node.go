package errors

func NewNodeError(msg string) (err error) {
	return NewError(ErrorPrefixNode, msg)
}

var ErrorNodeUnregistered = NewNodeError("unregistered")
var ErrorNodeServiceNotExists = NewNodeError("service not exists")
var ErrorNodeInvalidType = NewNodeError("invalid type")
var ErrorNodeInvalidStatus = NewNodeError("invalid status")
var ErrorNodeInvalidCode = NewNodeError("invalid code")
var ErrorNodeInvalidNodeKey = NewNodeError("invalid node key")
var ErrorNodeMonitorError = NewNodeError("monitor error")
var ErrorNodeNotExists = NewNodeError("not exists")
