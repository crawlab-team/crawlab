package interfaces

type GrpcSubscribe interface {
	GetStream() GrpcStream
	GetStreamBidirectional() GrpcStreamBidirectional
	GetFinished() chan bool
}
