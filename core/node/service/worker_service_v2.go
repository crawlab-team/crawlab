package service

import (
	"context"
	"encoding/json"
	"github.com/apex/log"
	config2 "github.com/crawlab-team/crawlab/core/config"
	"github.com/crawlab-team/crawlab/core/container"
	"github.com/crawlab-team/crawlab/core/grpc/client"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/task/handler"
	"github.com/crawlab-team/crawlab/core/utils"
	grpc "github.com/crawlab-team/crawlab/grpc"
	"github.com/crawlab-team/crawlab/trace"
	"time"
)

type WorkerServiceV2 struct {
	// dependencies
	cfgSvc     interfaces.NodeConfigService
	client     *client.GrpcClientV2
	handlerSvc *handler.ServiceV2

	// settings
	cfgPath           string
	address           interfaces.Address
	heartbeatInterval time.Duration

	// internals
	n interfaces.Node
	s grpc.NodeService_SubscribeClient
}

func (svc *WorkerServiceV2) Init() (err error) {
	// do nothing
	return nil
}

func (svc *WorkerServiceV2) Start() {
	// start grpc client
	if err := svc.client.Start(); err != nil {
		panic(err)
	}

	// register to master
	svc.Register()

	// start receiving stream messages
	go svc.Recv()

	// start sending heartbeat to master
	go svc.ReportStatus()

	// start handler
	go svc.handlerSvc.Start()

	// wait for quit signal
	svc.Wait()

	// stop
	svc.Stop()
}

func (svc *WorkerServiceV2) Wait() {
	utils.DefaultWait()
}

func (svc *WorkerServiceV2) Stop() {
	_ = svc.client.Stop()
	log.Infof("worker[%s] service has stopped", svc.cfgSvc.GetNodeKey())
}

func (svc *WorkerServiceV2) Register() {
	ctx, cancel := svc.client.Context()
	defer cancel()
	req := svc.client.NewRequest(svc.GetConfigService().GetBasicNodeInfo())
	res, err := svc.client.NodeClient.Register(ctx, req)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(res.Data, svc.n); err != nil {
		panic(err)
	}
	log.Infof("worker[%s] registered to master. id: %s", svc.GetConfigService().GetNodeKey(), svc.n.GetId().Hex())
	return
}

func (svc *WorkerServiceV2) Recv() {
	msgCh := svc.client.GetMessageChannel()
	for {
		// return if client is closed
		if svc.client.IsClosed() {
			return
		}

		// receive message from channel
		msg := <-msgCh

		// handle message
		if err := svc.handleStreamMessage(msg); err != nil {
			continue
		}
	}
}

func (svc *WorkerServiceV2) handleStreamMessage(msg *grpc.StreamMessage) (err error) {
	log.Debugf("[WorkerServiceV2] handle msg: %v", msg)
	switch msg.Code {
	case grpc.StreamMessageCode_PING:
		if _, err := svc.client.NodeClient.SendHeartbeat(context.Background(), svc.client.NewRequest(svc.cfgSvc.GetBasicNodeInfo())); err != nil {
			return trace.TraceError(err)
		}
	case grpc.StreamMessageCode_RUN_TASK:
		var t models.Task
		if err := json.Unmarshal(msg.Data, &t); err != nil {
			return trace.TraceError(err)
		}
		if err := svc.handlerSvc.Run(t.Id); err != nil {
			return trace.TraceError(err)
		}
	case grpc.StreamMessageCode_CANCEL_TASK:
		var t models.Task
		if err := json.Unmarshal(msg.Data, &t); err != nil {
			return trace.TraceError(err)
		}
		if err := svc.handlerSvc.Cancel(t.Id); err != nil {
			return trace.TraceError(err)
		}
	}

	return nil
}

func (svc *WorkerServiceV2) ReportStatus() {
	for {
		// return if client is closed
		if svc.client.IsClosed() {
			return
		}

		// report status
		svc.reportStatus()

		// sleep
		time.Sleep(svc.heartbeatInterval)
	}
}

func (svc *WorkerServiceV2) GetConfigService() (cfgSvc interfaces.NodeConfigService) {
	return svc.cfgSvc
}

func (svc *WorkerServiceV2) GetConfigPath() (path string) {
	return svc.cfgPath
}

func (svc *WorkerServiceV2) SetConfigPath(path string) {
	svc.cfgPath = path
}

func (svc *WorkerServiceV2) GetAddress() (address interfaces.Address) {
	return svc.address
}

func (svc *WorkerServiceV2) SetAddress(address interfaces.Address) {
	svc.address = address
}

func (svc *WorkerServiceV2) SetHeartbeatInterval(duration time.Duration) {
	svc.heartbeatInterval = duration
}

func (svc *WorkerServiceV2) reportStatus() {
	ctx, cancel := context.WithTimeout(context.Background(), svc.heartbeatInterval)
	defer cancel()
	_, err := svc.client.NodeClient.SendHeartbeat(ctx, &grpc.Request{
		NodeKey: svc.cfgSvc.GetNodeKey(),
	})
	if err != nil {
		trace.PrintError(err)
	}
}

func NewWorkerServiceV2() (res *WorkerServiceV2, err error) {
	svc := &WorkerServiceV2{
		cfgPath:           config2.GetConfigPath(),
		heartbeatInterval: 15 * time.Second,
		n:                 &models.Node{},
	}

	// dependency options
	var clientOpts []client.Option
	if svc.address != nil {
		clientOpts = append(clientOpts, client.WithAddress(svc.address))
	}

	// dependency injection
	if err := container.GetContainer().Invoke(func(
		cfgSvc interfaces.NodeConfigService,
	) {
		svc.cfgSvc = cfgSvc
	}); err != nil {
		return nil, err
	}

	// grpc client
	svc.client, err = client.NewGrpcClientV2()
	if err != nil {
		return nil, err
	}

	// handler service
	svc.handlerSvc, err = handler.GetTaskHandlerServiceV2()
	if err != nil {
		return nil, err
	}

	// init
	if err := svc.Init(); err != nil {
		return nil, err
	}

	return svc, nil
}
