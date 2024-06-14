package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/container"
	errors2 "github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/client"
	"github.com/crawlab-team/crawlab/core/models/delegate"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/core/task"
	"github.com/crawlab-team/crawlab/trace"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
	"time"
)

type Service struct {
	// dependencies
	interfaces.TaskBaseService
	cfgSvc                    interfaces.NodeConfigService
	modelSvc                  service.ModelService
	clientModelSvc            interfaces.GrpcClientModelService
	clientModelNodeSvc        interfaces.GrpcClientModelNodeService
	clientModelSpiderSvc      interfaces.GrpcClientModelSpiderService
	clientModelTaskSvc        interfaces.GrpcClientModelTaskService
	clientModelTaskStatSvc    interfaces.GrpcClientModelTaskStatService
	clientModelEnvironmentSvc interfaces.GrpcClientModelEnvironmentService
	c                         interfaces.GrpcClient // grpc client

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
		svc.c.Start()
	}

	go svc.ReportStatus()
	go svc.Fetch()
}

func (svc *Service) Run(taskId primitive.ObjectID) (err error) {
	return svc.run(taskId)
}

func (svc *Service) Reset() {
	svc.mu.Lock()
	defer svc.mu.Unlock()
}

func (svc *Service) Cancel(taskId primitive.ObjectID) (err error) {
	r, err := svc.getRunner(taskId)
	if err != nil {
		return err
	}
	if err := r.Cancel(); err != nil {
		return err
	}
	return nil
}

func (svc *Service) Fetch() {
	for {
		// wait
		time.Sleep(svc.fetchInterval)

		// current node
		n, err := svc.GetCurrentNode()
		if err != nil {
			continue
		}

		// skip if node is not active or enabled
		if !n.GetActive() || !n.GetEnabled() {
			continue
		}

		// validate if there are available runners
		if svc.getRunnerCount() >= n.GetMaxRunners() {
			continue
		}

		// stop
		if svc.stopped {
			return
		}

		// fetch task
		tid, err := svc.fetch()
		if err != nil {
			trace.PrintError(err)
			continue
		}

		// skip if no task id
		if tid.IsZero() {
			continue
		}

		// run task
		if err := svc.run(tid); err != nil {
			trace.PrintError(err)
			t, err := svc.GetTaskById(tid)
			if err == nil && t.GetStatus() != constants.TaskStatusCancelled {
				t.SetError(err.Error())
				_ = svc.SaveTask(t, constants.TaskStatusError)
				continue
			}
			continue
		}
	}
}

func (svc *Service) ReportStatus() {
	for {
		if svc.stopped {
			return
		}

		// report handler status
		if err := svc.reportStatus(); err != nil {
			trace.PrintError(err)
		}

		// wait
		time.Sleep(svc.reportInterval)
	}
}

func (svc *Service) IsSyncLocked(path string) (ok bool) {
	_, ok = svc.syncLocks.Load(path)
	return ok
}

func (svc *Service) LockSync(path string) {
	svc.syncLocks.Store(path, true)
}

func (svc *Service) UnlockSync(path string) {
	svc.syncLocks.Delete(path)
}

//func (svc *Service) GetMaxRunners() (maxRunners int) {
//	return svc.maxRunners
//}
//
//func (svc *Service) SetMaxRunners(maxRunners int) {
//	svc.maxRunners = maxRunners
//}

func (svc *Service) GetExitWatchDuration() (duration time.Duration) {
	return svc.exitWatchDuration
}

func (svc *Service) SetExitWatchDuration(duration time.Duration) {
	svc.exitWatchDuration = duration
}

func (svc *Service) GetFetchInterval() (interval time.Duration) {
	return svc.fetchInterval
}

func (svc *Service) SetFetchInterval(interval time.Duration) {
	svc.fetchInterval = interval
}

func (svc *Service) GetReportInterval() (interval time.Duration) {
	return svc.reportInterval
}

func (svc *Service) SetReportInterval(interval time.Duration) {
	svc.reportInterval = interval
}

func (svc *Service) GetCancelTimeout() (timeout time.Duration) {
	return svc.cancelTimeout
}

func (svc *Service) SetCancelTimeout(timeout time.Duration) {
	svc.cancelTimeout = timeout
}

func (svc *Service) GetModelService() (modelSvc interfaces.GrpcClientModelService) {
	return svc.clientModelSvc
}

func (svc *Service) GetModelSpiderService() (modelSpiderSvc interfaces.GrpcClientModelSpiderService) {
	return svc.clientModelSpiderSvc
}

func (svc *Service) GetModelTaskService() (modelTaskSvc interfaces.GrpcClientModelTaskService) {
	return svc.clientModelTaskSvc
}

func (svc *Service) GetModelTaskStatService() (modelTaskSvc interfaces.GrpcClientModelTaskStatService) {
	return svc.clientModelTaskStatSvc
}

func (svc *Service) GetModelEnvironmentService() (modelTaskSvc interfaces.GrpcClientModelEnvironmentService) {
	return svc.clientModelEnvironmentSvc
}

func (svc *Service) GetNodeConfigService() (cfgSvc interfaces.NodeConfigService) {
	return svc.cfgSvc
}

func (svc *Service) GetCurrentNode() (n interfaces.Node, err error) {
	// node key
	nodeKey := svc.cfgSvc.GetNodeKey()

	// current node
	if svc.cfgSvc.IsMaster() {
		n, err = svc.modelSvc.GetNodeByKey(nodeKey, nil)
	} else {
		n, err = svc.clientModelNodeSvc.GetNodeByKey(nodeKey)
	}
	if err != nil {
		return nil, err
	}

	return n, nil
}

func (svc *Service) GetTaskById(id primitive.ObjectID) (t interfaces.Task, err error) {
	if svc.cfgSvc.IsMaster() {
		t, err = svc.modelSvc.GetTaskById(id)
	} else {
		t, err = svc.clientModelTaskSvc.GetTaskById(id)
	}
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (svc *Service) GetSpiderById(id primitive.ObjectID) (s interfaces.Spider, err error) {
	if svc.cfgSvc.IsMaster() {
		s, err = svc.modelSvc.GetSpiderById(id)
	} else {
		s, err = svc.clientModelSpiderSvc.GetSpiderById(id)
	}
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (svc *Service) getRunners() (runners []*Runner) {
	svc.mu.Lock()
	defer svc.mu.Unlock()
	svc.runners.Range(func(key, value interface{}) bool {
		r := value.(Runner)
		runners = append(runners, &r)
		return true
	})
	return runners
}

func (svc *Service) getRunnerCount() (count int) {
	n, err := svc.GetCurrentNode()
	if err != nil {
		trace.PrintError(err)
		return
	}
	query := bson.M{
		"node_id": n.GetId(),
		"status":  constants.TaskStatusRunning,
	}
	if svc.cfgSvc.IsMaster() {
		count, err = svc.modelSvc.GetBaseService(interfaces.ModelIdTask).Count(query)
		if err != nil {
			trace.PrintError(err)
			return
		}
	} else {
		count, err = svc.clientModelTaskSvc.Count(query)
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

func (svc *Service) saveTask(t interfaces.Task, status string) (err error) {
	// normalize status
	if status == "" {
		status = constants.TaskStatusPending
	}

	// set task status
	t.SetStatus(status)

	// attempt to get task from database
	_, err = svc.clientModelTaskSvc.GetTaskById(t.GetId())
	if err != nil {
		// if task does not exist, add to database
		if err == mongo.ErrNoDocuments {
			if err := client.NewModelDelegate(t, client.WithDelegateConfigPath(svc.GetConfigPath())).Add(); err != nil {
				return err
			}
			return nil
		} else {
			return err
		}
	} else {
		// otherwise, update
		if err := client.NewModelDelegate(t, client.WithDelegateConfigPath(svc.GetConfigPath())).Save(); err != nil {
			return err
		}
		return nil
	}
}

func (svc *Service) reportStatus() (err error) {
	// current node
	n, err := svc.GetCurrentNode()
	if err != nil {
		return err
	}

	// available runners of handler
	ar := n.GetMaxRunners() - svc.getRunnerCount()

	// set available runners
	n.SetAvailableRunners(ar)

	// save node
	if svc.cfgSvc.IsMaster() {
		err = delegate.NewModelDelegate(n).Save()
	} else {
		err = client.NewModelDelegate(n, client.WithDelegateConfigPath(svc.GetConfigPath())).Save()
	}
	if err != nil {
		return err
	}

	return nil
}

func (svc *Service) fetch() (tid primitive.ObjectID, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), svc.fetchTimeout)
	defer cancel()
	res, err := svc.c.GetTaskClient().Fetch(ctx, svc.c.NewRequest(nil))
	if err != nil {
		return tid, trace.TraceError(err)
	}
	if err := json.Unmarshal(res.Data, &tid); err != nil {
		return tid, trace.TraceError(err)
	}
	return tid, nil
}

func (svc *Service) run(taskId primitive.ObjectID) (err error) {
	// attempt to get runner from pool
	_, ok := svc.runners.Load(taskId)
	if ok {
		return trace.TraceError(errors2.ErrorTaskAlreadyExists)
	}

	// create a new task runner
	r, err := NewTaskRunner(taskId, svc)
	if err != nil {
		return trace.TraceError(err)
	}

	// add runner to pool
	svc.addRunner(taskId, r)

	// create a goroutine to run task
	go func() {
		// delete runner from pool
		defer svc.deleteRunner(r.GetTaskId())
		defer func(r interfaces.TaskRunner) {
			err := r.CleanUp()
			if err != nil {
				log.Errorf("task[%s] clean up error: %v", r.GetTaskId().Hex(), err)
			}
		}(r)
		// run task process (blocking)
		// error or finish after task runner ends
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
	}()

	return nil
}

func NewTaskHandlerService() (svc2 interfaces.TaskHandlerService, err error) {
	// base service
	baseSvc, err := task.NewBaseService()
	if err != nil {
		return nil, trace.TraceError(err)
	}

	// service
	svc := &Service{
		TaskBaseService:   baseSvc,
		exitWatchDuration: 60 * time.Second,
		fetchInterval:     1 * time.Second,
		fetchTimeout:      15 * time.Second,
		reportInterval:    5 * time.Second,
		cancelTimeout:     5 * time.Second,
		mu:                sync.Mutex{},
		runners:           sync.Map{},
		syncLocks:         sync.Map{},
	}

	// dependency injection
	if err := container.GetContainer().Invoke(func(
		cfgSvc interfaces.NodeConfigService,
		modelSvc service.ModelService,
		clientModelSvc interfaces.GrpcClientModelService,
		clientModelNodeSvc interfaces.GrpcClientModelNodeService,
		clientModelSpiderSvc interfaces.GrpcClientModelSpiderService,
		clientModelTaskSvc interfaces.GrpcClientModelTaskService,
		clientModelTaskStatSvc interfaces.GrpcClientModelTaskStatService,
		clientModelEnvironmentSvc interfaces.GrpcClientModelEnvironmentService,
		c interfaces.GrpcClient,
	) {
		svc.cfgSvc = cfgSvc
		svc.modelSvc = modelSvc
		svc.clientModelSvc = clientModelSvc
		svc.clientModelNodeSvc = clientModelNodeSvc
		svc.clientModelSpiderSvc = clientModelSpiderSvc
		svc.clientModelTaskSvc = clientModelTaskSvc
		svc.clientModelTaskStatSvc = clientModelTaskStatSvc
		svc.clientModelEnvironmentSvc = clientModelEnvironmentSvc
		svc.c = c
	}); err != nil {
		return nil, trace.TraceError(err)
	}

	log.Debugf("[NewTaskHandlerService] svc[cfgPath: %s]", svc.cfgSvc.GetConfigPath())

	return svc, nil
}

var _service interfaces.TaskHandlerService

func GetTaskHandlerService() (svr interfaces.TaskHandlerService, err error) {
	if _service != nil {
		return _service, nil
	}
	_service, err = NewTaskHandlerService()
	if err != nil {
		return nil, err
	}
	return _service, nil
}
