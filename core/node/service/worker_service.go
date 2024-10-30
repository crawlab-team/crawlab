package service

import (
	"context"
	"errors"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab/core/config"
	"github.com/crawlab-team/crawlab/core/grpc/client"
	"github.com/crawlab-team/crawlab/core/interfaces"
	client2 "github.com/crawlab-team/crawlab/core/models/client"
	"github.com/crawlab-team/crawlab/core/models/models/v2"
	nodeconfig "github.com/crawlab-team/crawlab/core/node/config"
	"github.com/crawlab-team/crawlab/core/task/handler"
	"github.com/crawlab-team/crawlab/core/utils"
	"github.com/crawlab-team/crawlab/grpc"
	"github.com/crawlab-team/crawlab/trace"
	"go.mongodb.org/mongo-driver/bson"
	"io"
	"sync"
	"time"
)

type WorkerService struct {
	// dependencies
	cfgSvc     interfaces.NodeConfigService
	client     *client.GrpcClientV2
	handlerSvc *handler.ServiceV2

	// settings
	cfgPath           string
	address           interfaces.Address
	heartbeatInterval time.Duration

	// internals
	stopped bool
	n       *models.NodeV2
	s       grpc.NodeService_SubscribeClient
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

	// subscribe
	svc.Subscribe()

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
	svc.stopped = true
	_ = svc.client.Stop()
	svc.handlerSvc.Stop()
	log.Infof("worker[%s] service has stopped", svc.cfgSvc.GetNodeKey())
}

func (svc *WorkerService) Register() {
	ctx, cancel := svc.client.Context()
	defer cancel()
	_, err := svc.client.NodeClient.Register(ctx, &grpc.NodeServiceRegisterRequest{
		NodeKey:    svc.cfgSvc.GetNodeKey(),
		NodeName:   svc.cfgSvc.GetNodeName(),
		MaxRunners: int32(svc.cfgSvc.GetMaxRunners()),
	})
	if err != nil {
		panic(err)
	}
	svc.n, err = client2.NewModelServiceV2[models.NodeV2]().GetOne(bson.M{"key": svc.GetConfigService().GetNodeKey()}, nil)
	if err != nil {
		panic(err)
	}
	log.Infof("worker[%s] registered to master. id: %s", svc.GetConfigService().GetNodeKey(), svc.n.Id.Hex())
	return
}

func (svc *WorkerService) ReportStatus() {
	ticker := time.NewTicker(svc.heartbeatInterval)
	for {
		// return if client is closed
		if svc.client.IsClosed() {
			ticker.Stop()
			return
		}

		// report status
		svc.reportStatus()

		// sleep
		<-ticker.C
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

func (svc *WorkerService) Subscribe() {
	stream, err := svc.client.NodeClient.Subscribe(context.Background(), &grpc.NodeServiceSubscribeRequest{
		NodeKey: svc.cfgSvc.GetNodeKey(),
	})
	if err != nil {
		log.Errorf("failed to subscribe to master: %v", err)
		return
	}
	for {
		if svc.stopped {
			return
		}
		select {
		default:
			msg, err := stream.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					return
				}
				log.Errorf("failed to receive message from master: %v", err)
				continue
			}
			switch msg.Code {
			case grpc.NodeServiceSubscribeCode_PING:
				// do nothing
			}
		}
	}
}

func (svc *WorkerService) reportStatus() {
	ctx, cancel := context.WithTimeout(context.Background(), svc.heartbeatInterval)
	defer cancel()
	_, err := svc.client.NodeClient.SendHeartbeat(ctx, &grpc.NodeServiceSendHeartbeatRequest{
		NodeKey: svc.cfgSvc.GetNodeKey(),
	})
	if err != nil {
		trace.PrintError(err)
	}
}

var workerServiceV2 *WorkerService
var workerServiceV2Once = new(sync.Once)

func newWorkerService() (res *WorkerService, err error) {
	svc := &WorkerService{
		cfgPath:           config.GetConfigPath(),
		heartbeatInterval: 15 * time.Second,
	}

	// dependency options
	var clientOpts []client.Option
	if svc.address != nil {
		clientOpts = append(clientOpts, client.WithAddress(svc.address))
	}

	// node config service
	svc.cfgSvc = nodeconfig.GetNodeConfigService()

	// grpc client
	svc.client = client.GetGrpcClientV2()

	// handler service
	svc.handlerSvc, err = handler.GetTaskHandlerServiceV2()
	if err != nil {
		return nil, err
	}

	// init
	err = svc.Init()
	if err != nil {
		return nil, err
	}

	return svc, nil
}

func GetWorkerService() (res *WorkerService, err error) {
	workerServiceV2Once.Do(func() {
		workerServiceV2, err = newWorkerService()
		if err != nil {
			log.Errorf("failed to get worker service: %v", err)
		}
	})
	return workerServiceV2, err
}
