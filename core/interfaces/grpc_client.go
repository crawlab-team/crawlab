package interfaces

import (
	"context"
	grpc "github.com/crawlab-team/crawlab/grpc"
	"time"
)

type GrpcClient interface {
	GrpcBase
	WithConfigPath
	GetModelDelegateClient() grpc.ModelDelegateClient
	GetModelBaseServiceClient() grpc.ModelBaseServiceClient
	GetNodeClient() grpc.NodeServiceClient
	GetTaskClient() grpc.TaskServiceClient
	GetMessageClient() grpc.MessageServiceClient
	SetAddress(Address)
	SetTimeout(time.Duration)
	SetSubscribeType(string)
	SetHandleMessage(bool)
	Context() (context.Context, context.CancelFunc)
	NewRequest(interface{}) *grpc.Request
	GetMessageChannel() chan *grpc.StreamMessage
	Restart() error
	NewModelBaseServiceRequest(ModelId, GrpcBaseServiceParams) (*grpc.Request, error)
	IsStarted() bool
	IsClosed() bool
	Err() error
	GetStream() grpc.NodeService_SubscribeClient
}
