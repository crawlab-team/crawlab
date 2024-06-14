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
	"github.com/crawlab-team/crawlab/core/utils"
	grpc "github.com/crawlab-team/crawlab/grpc"
	"github.com/crawlab-team/crawlab/trace"
	"github.com/spf13/viper"
	"time"
)

type WorkerService struct {
	// dependencies
	cfgSvc     interfaces.NodeConfigService
	client     interfaces.GrpcClient
	handlerSvc interfaces.TaskHandlerService

	// settings
	cfgPath           string
	address           interfaces.Address
	heartbeatInterval time.Duration

	// internals
	n interfaces.Node
	s grpc.NodeService_SubscribeClient
}

func (svc *WorkerService) Init() (err error) {
	// do nothing
	return nil
}

func (svc *WorkerService) Start() {
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

func (svc *WorkerService) Wait() {
	utils.DefaultWait()
}

func (svc *WorkerService) Stop() {
	_ = svc.client.Stop()
	log.Infof("worker[%s] service has stopped", svc.cfgSvc.GetNodeKey())
}

func (svc *WorkerService) Register() {
	ctx, cancel := svc.client.Context()
	defer cancel()
	req := svc.client.NewRequest(svc.GetConfigService().GetBasicNodeInfo())
	res, err := svc.client.GetNodeClient().Register(ctx, req)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(res.Data, svc.n); err != nil {
		panic(err)
	}
	log.Infof("worker[%s] registered to master. id: %s", svc.GetConfigService().GetNodeKey(), svc.n.GetId().Hex())
	return
}

func (svc *WorkerService) Recv() {
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

func (svc *WorkerService) handleStreamMessage(msg *grpc.StreamMessage) (err error) {
	log.Debugf("[WorkerService] handle msg: %v", msg)
	switch msg.Code {
	case grpc.StreamMessageCode_PING:
		if _, err := svc.client.GetNodeClient().SendHeartbeat(context.Background(), svc.client.NewRequest(svc.cfgSvc.GetBasicNodeInfo())); err != nil {
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

func (svc *WorkerService) ReportStatus() {
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

func (svc *WorkerService) GetConfigService() (cfgSvc interfaces.NodeConfigService) {
	return svc.cfgSvc
}

func (svc *WorkerService) GetConfigPath() (path string) {
	return svc.cfgPath
}

func (svc *WorkerService) SetConfigPath(path string) {
	svc.cfgPath = path
}

func (svc *WorkerService) GetAddress() (address interfaces.Address) {
	return svc.address
}

func (svc *WorkerService) SetAddress(address interfaces.Address) {
	svc.address = address
}

func (svc *WorkerService) SetHeartbeatInterval(duration time.Duration) {
	svc.heartbeatInterval = duration
}

func (svc *WorkerService) reportStatus() {
	ctx, cancel := context.WithTimeout(context.Background(), svc.heartbeatInterval)
	defer cancel()
	_, err := svc.client.GetNodeClient().SendHeartbeat(ctx, &grpc.Request{
		NodeKey: svc.cfgSvc.GetNodeKey(),
	})
	if err != nil {
		trace.PrintError(err)
	}
}

func NewWorkerService(opts ...Option) (res *WorkerService, err error) {
	svc := &WorkerService{
		cfgPath:           config2.GetConfigPath(),
		heartbeatInterval: 15 * time.Second,
		n:                 &models.Node{},
	}

	// apply options
	for _, opt := range opts {
		opt(svc)
	}

	// dependency options
	var clientOpts []client.Option
	if svc.address != nil {
		clientOpts = append(clientOpts, client.WithAddress(svc.address))
	}

	// dependency injection
	if err := container.GetContainer().Invoke(func(
		cfgSvc interfaces.NodeConfigService,
		client interfaces.GrpcClient,
		taskHandlerSvc interfaces.TaskHandlerService,
	) {
		svc.cfgSvc = cfgSvc
		svc.client = client
		svc.handlerSvc = taskHandlerSvc
	}); err != nil {
		return nil, err
	}

	// init
	if err := svc.Init(); err != nil {
		return nil, err
	}

	return svc, nil
}

func ProvideWorkerService(path string, opts ...Option) func() (interfaces.NodeWorkerService, error) {
	// path
	if path == "" || path == config2.GetConfigPath() {
		if viper.GetString("config.path") != "" {
			path = viper.GetString("config.path")
		} else {
			path = config2.GetConfigPath()
		}
	}
	opts = append(opts, WithConfigPath(path))

	return func() (interfaces.NodeWorkerService, error) {
		return NewWorkerService(opts...)
	}
}
