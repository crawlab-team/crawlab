package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab/core/constants"
	errors2 "github.com/crawlab-team/crawlab/core/errors"
	grpcclient "github.com/crawlab-team/crawlab/core/grpc/client"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/client"
	models2 "github.com/crawlab-team/crawlab/core/models/models/v2"
	"github.com/crawlab-team/crawlab/core/models/service"
	nodeconfig "github.com/crawlab-team/crawlab/core/node/config"
	"github.com/crawlab-team/crawlab/grpc"
	"github.com/crawlab-team/crawlab/trace"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"sync"
	"time"
)

type Service struct {
	// dependencies
	cfgSvc interfaces.NodeConfigService
	c      *grpcclient.GrpcClient // grpc client

	// settings
	//maxRunners        int
	exitWatchDuration time.Duration
	reportInterval    time.Duration
	fetchInterval     time.Duration
	fetchTimeout      time.Duration
	cancelTimeout     time.Duration

	// internals variables
	stopped   bool
	mu        sync.Mutex
	runners   sync.Map // pool of task runners started
	syncLocks sync.Map // files sync locks map of task runners
}

func (svc *Service) Start() {
	// Initialize gRPC if not started
	if !svc.c.IsStarted() {
		err := svc.c.Start()
		if err != nil {
			return
		}
	}

	go svc.reportStatus()
	go svc.fetchAndRunTasks()
}

func (svc *Service) Stop() {
	svc.stopped = true
}

func (svc *Service) Run(taskId primitive.ObjectID) (err error) {
	return svc.runTask(taskId)
}

func (svc *Service) Cancel(taskId primitive.ObjectID, force bool) (err error) {
	return svc.cancelTask(taskId, force)
}

func (svc *Service) fetchAndRunTasks() {
	ticker := time.NewTicker(svc.fetchInterval)
	for {
		if svc.stopped {
			return
		}

		select {
		case <-ticker.C:
			// current node
			n, err := svc.GetCurrentNode()
			if err != nil {
				continue
			}

			// skip if node is not active or enabled
			if !n.Active || !n.Enabled {
				continue
			}

			// validate if there are available runners
			if svc.getRunnerCount() >= n.MaxRunners {
				continue
			}

			// fetch task id
			tid, err := svc.fetchTask()
			if err != nil {
				continue
			}

			// skip if no task id
			if tid.IsZero() {
				continue
			}

			// run task
			if err := svc.runTask(tid); err != nil {
				t, err := svc.GetTaskById(tid)
				if err != nil && t.Status != constants.TaskStatusCancelled {
					t.Error = err.Error()
					t.Status = constants.TaskStatusError
					t.SetUpdated(t.CreatedBy)
					_ = client.NewModelServiceV2[models2.TaskV2]().ReplaceById(t.Id, *t)
					continue
				}
				continue
			}
		}
	}
}

func (svc *Service) reportStatus() {
	ticker := time.NewTicker(svc.reportInterval)
	for {
		if svc.stopped {
			return
		}

		select {
		case <-ticker.C:
			// update node status
			if err := svc.updateNodeStatus(); err != nil {
				log.Errorf("failed to report status: %v", err)
			}
		}
	}
}

func (svc *Service) GetExitWatchDuration() (duration time.Duration) {
	return svc.exitWatchDuration
}

func (svc *Service) GetCancelTimeout() (timeout time.Duration) {
	return svc.cancelTimeout
}

func (svc *Service) GetNodeConfigService() (cfgSvc interfaces.NodeConfigService) {
	return svc.cfgSvc
}

func (svc *Service) GetCurrentNode() (n *models2.NodeV2, err error) {
	// node key
	nodeKey := svc.cfgSvc.GetNodeKey()

	// current node
	if svc.cfgSvc.IsMaster() {
		n, err = service.NewModelServiceV2[models2.NodeV2]().GetOne(bson.M{"key": nodeKey}, nil)
	} else {
		n, err = client.NewModelServiceV2[models2.NodeV2]().GetOne(bson.M{"key": nodeKey}, nil)
	}
	if err != nil {
		return nil, err
	}

	return n, nil
}

func (svc *Service) GetTaskById(id primitive.ObjectID) (t *models2.TaskV2, err error) {
	if svc.cfgSvc.IsMaster() {
		t, err = service.NewModelServiceV2[models2.TaskV2]().GetById(id)
	} else {
		t, err = client.NewModelServiceV2[models2.TaskV2]().GetById(id)
	}
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (svc *Service) UpdateTask(t *models2.TaskV2) (err error) {
	t.SetUpdated(t.CreatedBy)
	if svc.cfgSvc.IsMaster() {
		err = service.NewModelServiceV2[models2.TaskV2]().ReplaceById(t.Id, *t)
	} else {
		err = client.NewModelServiceV2[models2.TaskV2]().ReplaceById(t.Id, *t)
	}
	if err != nil {
		return err
	}
	return nil
}

func (svc *Service) GetSpiderById(id primitive.ObjectID) (s *models2.SpiderV2, err error) {
	if svc.cfgSvc.IsMaster() {
		s, err = service.NewModelServiceV2[models2.SpiderV2]().GetById(id)
	} else {
		s, err = client.NewModelServiceV2[models2.SpiderV2]().GetById(id)
	}
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (svc *Service) getRunnerCount() (count int) {
	n, err := svc.GetCurrentNode()
	if err != nil {
		log.Errorf("failed to get current node: %v", err)
		return
	}
	query := bson.M{
		"node_id": n.Id,
		"status": bson.M{
			"$in": []string{constants.TaskStatusAssigned, constants.TaskStatusRunning},
		},
	}
	if svc.cfgSvc.IsMaster() {
		count, err = service.NewModelServiceV2[models2.TaskV2]().Count(query)
		if err != nil {
			trace.PrintError(err)
			return
		}
	} else {
		count, err = client.NewModelServiceV2[models2.TaskV2]().Count(query)
		if err != nil {
			trace.PrintError(err)
			return
		}
	}
	return count
}

func (svc *Service) getRunner(taskId primitive.ObjectID) (r interfaces.TaskRunner, err error) {
	log.Debugf("[TaskHandlerService] getRunner: taskId[%v]", taskId)
	v, ok := svc.runners.Load(taskId)
	if !ok {
		return nil, trace.TraceError(errors2.ErrorTaskNotExists)
	}
	switch v.(type) {
	case interfaces.TaskRunner:
		r = v.(interfaces.TaskRunner)
	default:
		return nil, trace.TraceError(errors2.ErrorModelInvalidType)
	}
	return r, nil
}

func (svc *Service) addRunner(taskId primitive.ObjectID, r interfaces.TaskRunner) {
	log.Debugf("[TaskHandlerService] addRunner: taskId[%v]", taskId)
	svc.runners.Store(taskId, r)
}

func (svc *Service) deleteRunner(taskId primitive.ObjectID) {
	log.Debugf("[TaskHandlerService] deleteRunner: taskId[%v]", taskId)
	svc.runners.Delete(taskId)
}

func (svc *Service) updateNodeStatus() (err error) {
	// current node
	n, err := svc.GetCurrentNode()
	if err != nil {
		return err
	}

	// available runners of handler
	ar := n.MaxRunners - svc.getRunnerCount()

	// set available runners
	n.AvailableRunners = ar

	// save node
	n.SetUpdated(n.CreatedBy)
	if svc.cfgSvc.IsMaster() {
		err = service.NewModelServiceV2[models2.NodeV2]().ReplaceById(n.Id, *n)
	} else {
		err = client.NewModelServiceV2[models2.NodeV2]().ReplaceById(n.Id, *n)
	}
	if err != nil {
		return err
	}

	return nil
}

func (svc *Service) fetchTask() (tid primitive.ObjectID, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), svc.fetchTimeout)
	defer cancel()
	res, err := svc.c.TaskClient.FetchTask(ctx, &grpc.TaskServiceFetchTaskRequest{
		NodeKey: svc.cfgSvc.GetNodeKey(),
	})
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("fetchTask task error: %v", err)
	}
	// validate task id
	tid, err = primitive.ObjectIDFromHex(res.GetTaskId())
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("invalid task id: %s", res.GetTaskId())
	}
	return tid, nil
}

func (svc *Service) runTask(taskId primitive.ObjectID) (err error) {
	// attempt to get runner from pool
	_, ok := svc.runners.Load(taskId)
	if ok {
		err = fmt.Errorf("task[%s] already exists", taskId.Hex())
		log.Errorf("run task error: %v", err)
		return err
	}

	// create a new task runner
	r, err := NewTaskRunnerV2(taskId, svc)
	if err != nil {
		err = fmt.Errorf("failed to create task runner: %v", err)
		log.Errorf("run task error: %v", err)
		return err
	}

	// add runner to pool
	svc.addRunner(taskId, r)

	// create a goroutine to run task
	go func() {
		// get subscription stream
		stopCh := make(chan struct{})
		stream, err := svc.subscribeTask(r.GetTaskId())
		if err == nil {
			// create a goroutine to handle stream messages
			go svc.handleStreamMessages(r.GetTaskId(), stream, stopCh)
		} else {
			log.Errorf("failed to subscribe task[%s]: %v", r.GetTaskId().Hex(), err)
			log.Warnf("task[%s] will not be able to receive stream messages", r.GetTaskId().Hex())
		}

		// run task process (blocking) error or finish after task runner ends
		if err := r.Run(); err != nil {
			switch {
			case errors.Is(err, constants.ErrTaskError):
				log.Errorf("task[%s] finished with error: %v", r.GetTaskId().Hex(), err)
			case errors.Is(err, constants.ErrTaskCancelled):
				log.Errorf("task[%s] cancelled", r.GetTaskId().Hex())
			default:
				log.Errorf("task[%s] finished with unknown error: %v", r.GetTaskId().Hex(), err)
			}
		}
		log.Infof("task[%s] finished", r.GetTaskId().Hex())

		// send stopCh signal to stream message handler
		stopCh <- struct{}{}

		// delete runner from pool
		svc.deleteRunner(r.GetTaskId())
	}()

	return nil
}

func (svc *Service) subscribeTask(taskId primitive.ObjectID) (stream grpc.TaskService_SubscribeClient, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req := &grpc.TaskServiceSubscribeRequest{
		TaskId: taskId.Hex(),
	}
	stream, err = svc.c.TaskClient.Subscribe(ctx, req)
	if err != nil {
		log.Errorf("failed to subscribe task[%s]: %v", taskId.Hex(), err)
		return nil, err
	}
	return stream, nil
}

func (svc *Service) handleStreamMessages(taskId primitive.ObjectID, stream grpc.TaskService_SubscribeClient, stopCh chan struct{}) {
	for {
		select {
		case <-stopCh:
			err := stream.CloseSend()
			if err != nil {
				log.Errorf("task[%s] failed to close stream: %v", taskId.Hex(), err)
				return
			}
			return
		default:
			msg, err := stream.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					log.Infof("task[%s] received EOF, stream closed", taskId.Hex())
					return
				}
				log.Errorf("task[%s] stream error: %v", taskId.Hex(), err)
				continue
			}
			switch msg.Code {
			case grpc.TaskServiceSubscribeCode_CANCEL:
				log.Infof("task[%s] received cancel signal", taskId.Hex())
				go svc.handleCancel(msg, taskId)
			}
		}
	}
}

func (svc *Service) handleCancel(msg *grpc.TaskServiceSubscribeResponse, taskId primitive.ObjectID) {
	// validate task id
	if msg.TaskId != taskId.Hex() {
		log.Errorf("task[%s] received cancel signal for another task[%s]", taskId.Hex(), msg.TaskId)
		return
	}

	// cancel task
	err := svc.cancelTask(taskId, msg.Force)
	if err != nil {
		log.Errorf("task[%s] failed to cancel: %v", taskId.Hex(), err)
		return
	}
	log.Infof("task[%s] cancelled", taskId.Hex())

	// set task status as "cancelled"
	t, err := svc.GetTaskById(taskId)
	if err != nil {
		log.Errorf("task[%s] failed to get task: %v", taskId.Hex(), err)
		return
	}
	t.Status = constants.TaskStatusCancelled
	err = svc.UpdateTask(t)
	if err != nil {
		log.Errorf("task[%s] failed to update task: %v", taskId.Hex(), err)
	}
}

func (svc *Service) cancelTask(taskId primitive.ObjectID, force bool) (err error) {
	r, err := svc.getRunner(taskId)
	if err != nil {
		return err
	}
	if err := r.Cancel(force); err != nil {
		return err
	}
	return nil
}

func newTaskHandlerService() (svc2 *Service, err error) {
	// service
	svc := &Service{
		exitWatchDuration: 60 * time.Second,
		fetchInterval:     1 * time.Second,
		fetchTimeout:      15 * time.Second,
		reportInterval:    5 * time.Second,
		cancelTimeout:     5 * time.Second,
		mu:                sync.Mutex{},
		runners:           sync.Map{},
	}

	// dependency injection
	svc.cfgSvc = nodeconfig.GetNodeConfigService()

	// grpc client
	svc.c = grpcclient.GetGrpcClient()

	log.Debugf("[NewTaskHandlerService] svc[cfgPath: %s]", svc.cfgSvc.GetConfigPath())

	return svc, nil
}

var _serviceV2 *Service
var _serviceV2Once = new(sync.Once)

func GetTaskHandlerService() (svr *Service, err error) {
	if _serviceV2 != nil {
		return _serviceV2, nil
	}
	_serviceV2Once.Do(func() {
		_serviceV2, err = newTaskHandlerService()
		if err != nil {
			log.Errorf("failed to create task handler service: %v", err)
		}
	})
	if err != nil {
		return nil, err
	}
	return _serviceV2, nil
}
