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

type MetricsServerV2 struct {
	grpc.UnimplementedMetricsServiceV2Server
}

func (svr MetricsServerV2) Send(_ context.Context, req *grpc.MetricsServiceV2SendRequest) (res *grpc.Response, err error) {
	log.Info("[MetricsServerV2] received metric from node: " + req.NodeKey)
	n, err := service.NewModelServiceV2[models2.NodeV2]().GetOne(bson.M{"key": req.NodeKey}, nil)
	if err != nil {
		log.Errorf("[MetricsServerV2] error getting node: %v", err)
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
	_, err = service.NewModelServiceV2[models2.MetricV2]().InsertOne(metric)
	if err != nil {
		log.Errorf("[MetricsServerV2] error inserting metric: %v", err)
		return HandleError(err)
	}
	return HandleSuccess()
}

func newMetricsServerV2() *MetricsServerV2 {
	return &MetricsServerV2{}
}

var metricsServerV2 *MetricsServerV2
var metricsServerV2Once = &sync.Once{}

func GetMetricsServerV2() *MetricsServerV2 {
	if metricsServerV2 != nil {
		return metricsServerV2
	}
	metricsServerV2Once.Do(func() {
		metricsServerV2 = newMetricsServerV2()
	})
	return metricsServerV2
}
