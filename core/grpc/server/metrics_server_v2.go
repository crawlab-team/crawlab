package server

import (
	"github.com/crawlab-team/crawlab/grpc"
	"sync"
)

type MetricsServerV2 struct {
	grpc.UnimplementedMetricsServiceV2Server
	mu      *sync.Mutex
	streams map[string]*grpc.MetricsServiceV2_ConnectServer
}
