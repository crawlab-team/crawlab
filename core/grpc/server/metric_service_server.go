package server

import (
	"context"
	"github.com/apex/log"
	models2 "github.com/crawlab-team/crawlab/core/models/models/v2"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/grpc"
	"go.mongodb.org/mongo-driver/bson"
	"sync"
	"time"
)

type MetricServiceServer struct {
	grpc.UnimplementedMetricServiceServer
}

func (svr MetricServiceServer) Send(_ context.Context, req *grpc.MetricServiceSendRequest) (res *grpc.Response, err error) {
	log.Info("[MetricServiceServer] received metric from node: " + req.NodeKey)
	n, err := service.NewModelService[models2.NodeV2]().GetOne(bson.M{"key": req.NodeKey}, nil)
	if err != nil {
		log.Errorf("[MetricServiceServer] error getting node: %v", err)
		return HandleError(err)
	}
	metric := models2.MetricV2{
		Type:                 req.Type,
		NodeId:               n.Id,
		CpuUsagePercent:      req.CpuUsagePercent,
		TotalMemory:          req.TotalMemory,
		AvailableMemory:      req.AvailableMemory,
		UsedMemory:           req.UsedMemory,
		UsedMemoryPercent:    req.UsedMemoryPercent,
		TotalDisk:            req.TotalDisk,
		AvailableDisk:        req.AvailableDisk,
		UsedDisk:             req.UsedDisk,
		UsedDiskPercent:      req.UsedDiskPercent,
		DiskReadBytesRate:    req.DiskReadBytesRate,
		DiskWriteBytesRate:   req.DiskWriteBytesRate,
		NetworkBytesSentRate: req.NetworkBytesSentRate,
		NetworkBytesRecvRate: req.NetworkBytesRecvRate,
	}
	metric.CreatedAt = time.Unix(req.Timestamp, 0)
	_, err = service.NewModelService[models2.MetricV2]().InsertOne(metric)
	if err != nil {
		log.Errorf("[MetricServiceServer] error inserting metric: %v", err)
		return HandleError(err)
	}
	return HandleSuccess()
}

func newMetricsServerV2() *MetricServiceServer {
	return &MetricServiceServer{}
}

var metricsServerV2 *MetricServiceServer
var metricsServerV2Once = &sync.Once{}

func GetMetricsServerV2() *MetricServiceServer {
	if metricsServerV2 != nil {
		return metricsServerV2
	}
	metricsServerV2Once.Do(func() {
		metricsServerV2 = newMetricsServerV2()
	})
	return metricsServerV2
}
