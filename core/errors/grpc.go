package errors

func NewGrpcError(msg string) (err error) {
	return NewError(ErrorPrefixGrpc, msg)
}

var (
	ErrorGrpcClientFailedToStart  = NewGrpcError("client failed to start")
	ErrorGrpcServerFailedToListen = NewGrpcError("server failed to listen")
	ErrorGrpcServerFailedToServe  = NewGrpcError("server failed to serve")
	ErrorGrpcUnauthorized         = NewGrpcError("unauthorized")
)
