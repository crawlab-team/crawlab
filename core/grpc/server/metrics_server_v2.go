package server

import (
	"errors"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab/grpc"
	"io"
	"sync"
)

type MetricsServerV2 struct {
	grpc.UnimplementedMetricsServiceV2Server
	mu       *sync.Mutex
	streams  map[string]*grpc.MetricsServiceV2_ConnectServer
	channels map[string]chan []*grpc.Metric
}

func (svr MetricsServerV2) Connect(stream grpc.MetricsServiceV2_ConnectServer) (err error) {
	// receive first message
	req, err := stream.Recv()
	if err != nil {
		log.Errorf("[MetricsServerV2] receive error: %v", err)
		return err
	}

	// save stream and channel
	svr.mu.Lock()
	svr.streams[req.NodeKey] = &stream
	svr.channels[req.NodeKey] = make(chan []*grpc.Metric)
	svr.mu.Unlock()

	log.Info("[MetricsServerV2] connected: " + req.NodeKey)

	for {
		// receive metrics
		req, err = stream.Recv()
		if errors.Is(err, io.EOF) {
			log.Errorf("[MetricsServerV2] receive EOF: %v", err)
			return
		}

		// send metrics to channel
		svr.channels[req.NodeKey] <- req.Metrics

		// keep this scope alive because once this scope exits - the stream is closed
		select {
		case <-stream.Context().Done():
			log.Info("[MetricsServerV2] disconnected: " + req.NodeKey)
			delete(svr.streams, req.NodeKey)
			delete(svr.channels, req.NodeKey)
			return nil
		}
	}
}

func (svr MetricsServerV2) GetStream(nodeKey string) (stream *grpc.MetricsServiceV2_ConnectServer, ok bool) {
	svr.mu.Lock()
	defer svr.mu.Unlock()
	stream, ok = svr.streams[nodeKey]
	return stream, ok
}

func (svr MetricsServerV2) GetChannel(nodeKey string) (ch chan []*grpc.Metric, ok bool) {
	svr.mu.Lock()
	defer svr.mu.Unlock()
	ch, ok = svr.channels[nodeKey]
	return ch, ok
}

func NewMetricsServerV2() *MetricsServerV2 {
	return &MetricsServerV2{
		mu:       new(sync.Mutex),
		streams:  make(map[string]*grpc.MetricsServiceV2_ConnectServer),
		channels: make(map[string]chan []*grpc.Metric),
	}
}

var metricsServerV2 *MetricsServerV2

func GetMetricsServerV2() *MetricsServerV2 {
	if metricsServerV2 != nil {
		return metricsServerV2
	}
	metricsServerV2 = NewMetricsServerV2()
	return metricsServerV2
}
