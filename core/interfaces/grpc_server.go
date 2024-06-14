package interfaces

import (
	grpc "github.com/crawlab-team/crawlab/grpc"
)

type GrpcServer interface {
	GrpcBase
	SetAddress(Address)
	GetSubscribe(key string) (sub GrpcSubscribe, err error)
	SetSubscribe(key string, sub GrpcSubscribe)
	DeleteSubscribe(key string)
	SendStreamMessage(key string, code grpc.StreamMessageCode) (err error)
	SendStreamMessageWithData(nodeKey string, code grpc.StreamMessageCode, d interface{}) (err error)
	IsStopped() (res bool)
}
