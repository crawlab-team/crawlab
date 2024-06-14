package interfaces

import grpc "github.com/crawlab-team/crawlab/grpc"

type GrpcStream interface {
	Send(msg *grpc.StreamMessage) (err error)
}

type GrpcStreamBidirectional interface {
	GrpcStream
	Recv() (msg *grpc.StreamMessage, err error)
}
