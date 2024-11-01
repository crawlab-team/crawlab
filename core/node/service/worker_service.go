package service

import (
	"context"
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/cenkalti/backoff/v4"
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
)

type WorkerService struct {
	// dependencies
	cfgSvc     interfaces.NodeConfigService
	client     *client.GrpcClient
	handlerSvc *handler.Service

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
	svc.register()

	// subscribe
	go svc.subscribe()

	// start sending heartbeat to master
	go svc.reportStatus()

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

func (svc *WorkerService) register() {
	ctx, cancel := svc.client.Context()
	defer cancel()
	_, err := svc.client.NodeClient.Register(ctx, &grpc.NodeServiceRegisterRequest{
		NodeKey:    svc.cfgSvc.GetNodeKey(),
		NodeName:   svc.cfgSvc.GetNodeName(),
		MaxRunners: int32(svc.cfgSvc.GetMaxRunners()),
	})
	if err != nil {
		log.Fatalf("failed to register worker[%s] to master: %v", svc.cfgSvc.GetNodeKey(), err)
		panic(err)
	}
	svc.n, err = client2.NewModelService[models.NodeV2]().GetOne(bson.M{"key": svc.GetConfigService().GetNodeKey()}, nil)
	if err != nil {
		log.Fatalf("failed to get node: %v", err)
		panic(err)
	}
	log.Infof("worker[%s] registered to master. id: %s", svc.GetConfigService().GetNodeKey(), svc.n.Id.Hex())
	return
}

func (svc *WorkerService) reportStatus() {
	ticker := time.NewTicker(svc.heartbeatInterval)
	for {
		// return if client is closed
		if svc.client.IsClosed() {
			ticker.Stop()
			return
		}

		// send heartbeat
		svc.sendHeartbeat()

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

func (svc *WorkerService) subscribe() {
	// Configure exponential backoff
	b := backoff.NewExponentialBackOff()
	b.InitialInterval = 1 * time.Second
	b.MaxInterval = 1 * time.Minute
	b.MaxElapsedTime = 10 * time.Minute
	b.Multiplier = 2.0

	for {
		if svc.stopped {
			return
		}

		// Use backoff for connection attempts
		operation := func() error {
			stream, err := svc.client.NodeClient.Subscribe(context.Background(), &grpc.NodeServiceSubscribeRequest{
				NodeKey: svc.cfgSvc.GetNodeKey(),
			})
			if err != nil {
				log.Errorf("failed to subscribe to master: %v", err)
				return err
			}

			// Handle messages
			for {
				if svc.stopped {
					return nil
				}

				msg, err := stream.Recv()
				if err != nil {
					if svc.client.IsClosed() {
						log.Errorf("connection to master is closed: %v", err)
						return err
					}
					log.Errorf("failed to receive message from master: %v", err)
					return err
				}

				switch msg.Code {
				case grpc.NodeServiceSubscribeCode_PING:
					// do nothing
				}
			}
		}

		// Execute with backoff
		err := backoff.Retry(operation, b)
		if err != nil {
			log.Errorf("subscription failed after max retries: %v", err)
			return
		}

		// Wait before attempting to reconnect
		time.Sleep(time.Second)
	}
}

func (svc *WorkerService) sendHeartbeat() {
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

	// node config service
	svc.cfgSvc = nodeconfig.GetNodeConfigService()

	// grpc client
	svc.client = client.GetGrpcClient()

	// handler service
	svc.handlerSvc, err = handler.GetTaskHandlerService()
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
