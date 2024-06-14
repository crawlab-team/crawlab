package errors

func NewGrpcError(msg string) (err error) {
	return NewError(ErrorPrefixGrpc, msg)
}

var (
	ErrorGrpcClientFailedToStart  = NewGrpcError("client failed to start")
	ErrorGrpcServerFailedToListen = NewGrpcError("server failed to listen")
	ErrorGrpcServerFailedToServe  = NewGrpcError("server failed to serve")
	ErrorGrpcClientNotExists      = NewGrpcError("client not exists")
	ErrorGrpcClientAlreadyExists  = NewGrpcError("client already exists")
	ErrorGrpcInvalidType          = NewGrpcError("invalid type")
	ErrorGrpcNotAllowed           = NewGrpcError("not allowed")
	ErrorGrpcSubscribeNotExists   = NewGrpcError("subscribe not exists")
	ErrorGrpcStreamNotFound       = NewGrpcError("stream not found")
	ErrorGrpcInvalidCode          = NewGrpcError("invalid code")
	ErrorGrpcUnauthorized         = NewGrpcError("unauthorized")
	ErrorGrpcInvalidNodeKey       = NewGrpcError("invalid node key")
)
