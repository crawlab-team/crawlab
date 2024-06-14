package service

import (
	"errors"
	"github.com/apex/log"
	"github.com/cenkalti/backoff/v4"
	config2 "github.com/crawlab-team/crawlab/core/config"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/container"
	"github.com/crawlab-team/crawlab/core/grpc/server"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/common"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/core/node/config"
	"github.com/crawlab-team/crawlab/core/notification"
	"github.com/crawlab-team/crawlab/core/schedule"
	"github.com/crawlab-team/crawlab/core/spider/admin"
	"github.com/crawlab-team/crawlab/core/system"
	"github.com/crawlab-team/crawlab/core/task/handler"
	"github.com/crawlab-team/crawlab/core/task/scheduler"
	"github.com/crawlab-team/crawlab/core/utils"
	grpc "github.com/crawlab-team/crawlab/grpc"
	"github.com/crawlab-team/crawlab/trace"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"sync"
	"time"
)

type MasterServiceV2 struct {
	// dependencies
	cfgSvc          interfaces.NodeConfigService
	server          *server.GrpcServerV2
	schedulerSvc    *scheduler.ServiceV2
	handlerSvc      *handler.ServiceV2
	scheduleSvc     *schedule.ServiceV2
	notificationSvc *notification.Service
	spiderAdminSvc  *admin.ServiceV2
	systemSvc       *system.Service

	// settings
	cfgPath         string
	address         interfaces.Address
	monitorInterval time.Duration
	stopOnError     bool
}

func (svc *MasterServiceV2) Init() (err error) {
	// do nothing
	return nil
}

func (svc *MasterServiceV2) Start() {
	// create indexes
	common.CreateIndexes()

	// start grpc server
	if err := svc.server.Start(); err != nil {
		panic(err)
	}

	// register to db
	if err := svc.Register(); err != nil {
		panic(err)
	}

	// start monitoring worker nodes
	go svc.Monitor()

	// start task handler
	go svc.handlerSvc.Start()

	// start task scheduler
	go svc.schedulerSvc.Start()

	// start schedule service
	go svc.scheduleSvc.Start()

	// start notification service
	go svc.notificationSvc.Start()

	// start spider admin service
	go svc.spiderAdminSvc.Start()

	// wait for quit signal
	svc.Wait()

	// stop
	svc.Stop()
}

func (svc *MasterServiceV2) Wait() {
	utils.DefaultWait()
}

func (svc *MasterServiceV2) Stop() {
	_ = svc.server.Stop()
	log.Infof("master[%s] service has stopped", svc.GetConfigService().GetNodeKey())
}

func (svc *MasterServiceV2) Monitor() {
	log.Infof("master[%s] monitoring started", svc.GetConfigService().GetNodeKey())
	for {
		if err := svc.monitor(); err != nil {
			trace.PrintError(err)
			if svc.stopOnError {
				log.Errorf("master[%s] monitor error, now stopping...", svc.GetConfigService().GetNodeKey())
				svc.Stop()
				return
			}
		}

		time.Sleep(svc.monitorInterval)
	}
}

func (svc *MasterServiceV2) GetConfigService() (cfgSvc interfaces.NodeConfigService) {
	return svc.cfgSvc
}

func (svc *MasterServiceV2) GetConfigPath() (path string) {
	return svc.cfgPath
}

func (svc *MasterServiceV2) SetConfigPath(path string) {
	svc.cfgPath = path
}

func (svc *MasterServiceV2) GetAddress() (address interfaces.Address) {
	return svc.address
}

func (svc *MasterServiceV2) SetAddress(address interfaces.Address) {
	svc.address = address
}

func (svc *MasterServiceV2) SetMonitorInterval(duration time.Duration) {
	svc.monitorInterval = duration
}

func (svc *MasterServiceV2) Register() (err error) {
	nodeKey := svc.GetConfigService().GetNodeKey()
	nodeName := svc.GetConfigService().GetNodeName()
	node, err := service.NewModelServiceV2[models.NodeV2]().GetOne(bson.M{"key": nodeKey}, nil)
	if err != nil && err.Error() == mongo2.ErrNoDocuments.Error() {
		// not exists
		log.Infof("master[%s] does not exist in db", nodeKey)
		node := models.NodeV2{
			Key:        nodeKey,
			Name:       nodeName,
			MaxRunners: config.DefaultConfigOptions.MaxRunners,
			IsMaster:   true,
			Status:     constants.NodeStatusOnline,
			Enabled:    true,
			Active:     true,
			ActiveAt:   time.Now(),
		}
		node.SetCreated(primitive.NilObjectID)
		node.SetUpdated(primitive.NilObjectID)
		id, err := service.NewModelServiceV2[models.NodeV2]().InsertOne(node)
		if err != nil {
			return err
		}
		log.Infof("added master[%s] in db. id: %s", nodeKey, id.Hex())
		return nil
	} else if err == nil {
		// exists
		log.Infof("master[%s] exists in db", nodeKey)
		node.Status = constants.NodeStatusOnline
		node.Active = true
		node.ActiveAt = time.Now()
		err = service.NewModelServiceV2[models.NodeV2]().ReplaceById(node.Id, *node)
		if err != nil {
			return err
		}
		log.Infof("updated master[%s] in db. id: %s", nodeKey, node.Id.Hex())
		return nil
	} else {
		// error
		return err
	}
}

func (svc *MasterServiceV2) StopOnError() {
	svc.stopOnError = true
}

func (svc *MasterServiceV2) GetServer() (svr interfaces.GrpcServer) {
	return svc.server
}

func (svc *MasterServiceV2) monitor() (err error) {
	// update master node status in db
	if err := svc.updateMasterNodeStatus(); err != nil {
		if err.Error() == mongo2.ErrNoDocuments.Error() {
			return nil
		}
		return err
	}

	// all worker nodes
	workerNodes, err := svc.getAllWorkerNodes()
	if err != nil {
		return err
	}

	// iterate all worker nodes
	wg := sync.WaitGroup{}
	wg.Add(len(workerNodes))
	for _, n := range workerNodes {
		go func(n *models.NodeV2) {
			// subscribe
			ok := svc.subscribeNode(n)
			if !ok {
				go svc.setWorkerNodeOffline(n)
				wg.Done()
				return
			}

			// ping client
			ok = svc.pingNodeClient(n)
			if !ok {
				go svc.setWorkerNodeOffline(n)
				wg.Done()
				return
			}

			// update node available runners
			if err := svc.updateNodeAvailableRunners(n); err != nil {
				trace.PrintError(err)
				wg.Done()
				return
			}

			// done
			wg.Done()
		}(&n)
	}

	wg.Wait()

	return nil
}

func (svc *MasterServiceV2) getAllWorkerNodes() (nodes []models.NodeV2, err error) {
	query := bson.M{
		"key":    bson.M{"$ne": svc.cfgSvc.GetNodeKey()}, // not self
		"active": true,                                   // active
	}
	nodes, err = service.NewModelServiceV2[models.NodeV2]().GetMany(query, nil)
	if err != nil {
		if errors.Is(err, mongo2.ErrNoDocuments) {
			return nil, nil
		}
		return nil, trace.TraceError(err)
	}
	return nodes, nil
}

func (svc *MasterServiceV2) updateMasterNodeStatus() (err error) {
	nodeKey := svc.GetConfigService().GetNodeKey()
	node, err := service.NewModelServiceV2[models.NodeV2]().GetOne(bson.M{"key": nodeKey}, nil)
	if err != nil {
		return err
	}
	node.Status = constants.NodeStatusOnline
	node.Active = true
	node.ActiveAt = time.Now()
	err = service.NewModelServiceV2[models.NodeV2]().ReplaceById(node.Id, *node)
	if err != nil {
		return err
	}
	return nil
}

func (svc *MasterServiceV2) setWorkerNodeOffline(node *models.NodeV2) {
	node.Status = constants.NodeStatusOffline
	node.Active = false
	err := backoff.Retry(func() error {
		return service.NewModelServiceV2[models.NodeV2]().ReplaceById(node.Id, *node)
	}, backoff.WithMaxRetries(backoff.NewConstantBackOff(1*time.Second), 3))
	if err != nil {
		trace.PrintError(err)
	}
}

func (svc *MasterServiceV2) subscribeNode(n *models.NodeV2) (ok bool) {
	_, err := svc.server.GetSubscribe("node:" + n.Key)
	if err != nil {
		log.Errorf("cannot subscribe worker node[%s]: %v", n.Key, err)
		return false
	}
	return true
}

func (svc *MasterServiceV2) pingNodeClient(n *models.NodeV2) (ok bool) {
	if err := svc.server.SendStreamMessage("node:"+n.Key, grpc.StreamMessageCode_PING); err != nil {
		log.Errorf("cannot ping worker node client[%s]: %v", n.Key, err)
		return false
	}
	return true
}

func (svc *MasterServiceV2) updateNodeAvailableRunners(node *models.NodeV2) (err error) {
	query := bson.M{
		"node_id": node.Id,
		"status":  constants.TaskStatusRunning,
	}
	runningTasksCount, err := service.NewModelServiceV2[models.TaskV2]().Count(query)
	if err != nil {
		return trace.TraceError(err)
	}
	node.AvailableRunners = node.MaxRunners - runningTasksCount
	err = service.NewModelServiceV2[models.NodeV2]().ReplaceById(node.Id, *node)
	if err != nil {
		return err
	}
	return nil
}

func NewMasterServiceV2() (res interfaces.NodeMasterService, err error) {
	// master service
	svc := &MasterServiceV2{
		cfgPath:         config2.GetConfigPath(),
		monitorInterval: 15 * time.Second,
		stopOnError:     false,
	}

	// server options
	var serverOpts []server.Option
	if svc.address != nil {
		serverOpts = append(serverOpts, server.WithAddress(svc.address))
	}

	// dependency injection
	if err := container.GetContainer().Invoke(func(
		cfgSvc interfaces.NodeConfigService,
	) {
		svc.cfgSvc = cfgSvc
	}); err != nil {
		return nil, err
	}

	// grpc server
	svc.server, err = server.GetGrpcServerV2()
	if err != nil {
		return nil, err
	}

	// scheduler service
	svc.schedulerSvc, err = scheduler.GetTaskSchedulerServiceV2()
	if err != nil {
		return nil, err
	}

	// handler service
	svc.handlerSvc, err = handler.GetTaskHandlerServiceV2()
	if err != nil {
		return nil, err
	}

	// schedule service
	svc.scheduleSvc, err = schedule.GetScheduleServiceV2()
	if err != nil {
		return nil, err
	}

	// notification service
	svc.notificationSvc = notification.GetService()

	// spider admin service
	svc.spiderAdminSvc, err = admin.GetSpiderAdminServiceV2()
	if err != nil {
		return nil, err
	}

	// system service
	svc.systemSvc = system.GetService()

	// init
	if err := svc.Init(); err != nil {
		return nil, err
	}

	return svc, nil
}
