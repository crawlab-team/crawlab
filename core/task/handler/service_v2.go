package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab/core/constants"
	errors2 "github.com/crawlab-team/crawlab/core/errors"
	grpcclient "github.com/crawlab-team/crawlab/core/grpc/client"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/client"
	models2 "github.com/crawlab-team/crawlab/core/models/models/v2"
	"github.com/crawlab-team/crawlab/core/models/service"
	nodeconfig "github.com/crawlab-team/crawlab/core/node/config"
	"github.com/crawlab-team/crawlab/trace"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
	"time"
)

type ServiceV2 struct {
	// dependencies
	cfgSvc interfaces.NodeConfigService
	c      *grpcclient.GrpcClientV2 // grpc client

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

func (svc *ServiceV2) Start() {
	// Initialize gRPC if not started
	if !svc.c.IsStarted() {
		err := svc.c.Start()
		if err != nil {
			return
		}
	}

	go svc.ReportStatus()
	go svc.Fetch()
}

func (svc *ServiceV2) Run(taskId primitive.ObjectID) (err error) {
	return svc.run(taskId)
}

func (svc *ServiceV2) Reset() {
	svc.mu.Lock()
	defer svc.mu.Unlock()
}

func (svc *ServiceV2) Cancel(taskId primitive.ObjectID) (err error) {
	r, err := svc.getRunner(taskId)
	if err != nil {
		return err
	}
	if err := r.Cancel(); err != nil {
		return err
	}
	return nil
}

func (svc *ServiceV2) Fetch() {
	ticker := time.NewTicker(svc.fetchInterval)
	for {
		// wait
		<-ticker.C

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

		// stop
		if svc.stopped {
			ticker.Stop()
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

func (svc *ServiceV2) ReportStatus() {
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

func (svc *ServiceV2) IsSyncLocked(path string) (ok bool) {
	_, ok = svc.syncLocks.Load(path)
	return ok
}

func (svc *ServiceV2) LockSync(path string) {
	svc.syncLocks.Store(path, true)
}

func (svc *ServiceV2) UnlockSync(path string) {
	svc.syncLocks.Delete(path)
}

//func (svc *ServiceV2) GetMaxRunners() (maxRunners int) {
//	return svc.maxRunners
//}
//
//func (svc *ServiceV2) SetMaxRunners(maxRunners int) {
//	svc.maxRunners = maxRunners
//}

func (svc *ServiceV2) GetExitWatchDuration() (duration time.Duration) {
	return svc.exitWatchDuration
}

func (svc *ServiceV2) SetExitWatchDuration(duration time.Duration) {
	svc.exitWatchDuration = duration
}

func (svc *ServiceV2) GetFetchInterval() (interval time.Duration) {
	return svc.fetchInterval
}

func (svc *ServiceV2) SetFetchInterval(interval time.Duration) {
	svc.fetchInterval = interval
}

func (svc *ServiceV2) GetReportInterval() (interval time.Duration) {
	return svc.reportInterval
}

func (svc *ServiceV2) SetReportInterval(interval time.Duration) {
	svc.reportInterval = interval
}

func (svc *ServiceV2) GetCancelTimeout() (timeout time.Duration) {
	return svc.cancelTimeout
}

func (svc *ServiceV2) SetCancelTimeout(timeout time.Duration) {
	svc.cancelTimeout = timeout
}

func (svc *ServiceV2) GetNodeConfigService() (cfgSvc interfaces.NodeConfigService) {
	return svc.cfgSvc
}

func (svc *ServiceV2) GetCurrentNode() (n *models2.NodeV2, err error) {
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

func (svc *ServiceV2) GetTaskById(id primitive.ObjectID) (t *models2.TaskV2, err error) {
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

func (svc *ServiceV2) GetSpiderById(id primitive.ObjectID) (s *models2.SpiderV2, err error) {
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

func (svc *ServiceV2) getRunners() (runners []*RunnerV2) {
	svc.mu.Lock()
	defer svc.mu.Unlock()
	svc.runners.Range(func(key, value interface{}) bool {
		r := value.(RunnerV2)
		runners = append(runners, &r)
		return true
	})
	return runners
}

func (svc *ServiceV2) getRunnerCount() (count int) {
	n, err := svc.GetCurrentNode()
	if err != nil {
		trace.PrintError(err)
		return
	}
	query := bson.M{
		"node_id": n.Id,
		"status":  constants.TaskStatusRunning,
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

func (svc *ServiceV2) getRunner(taskId primitive.ObjectID) (r interfaces.TaskRunner, err error) {
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

func (svc *ServiceV2) addRunner(taskId primitive.ObjectID, r interfaces.TaskRunner) {
	log.Debugf("[TaskHandlerService] addRunner: taskId[%v]", taskId)
	svc.runners.Store(taskId, r)
}

func (svc *ServiceV2) deleteRunner(taskId primitive.ObjectID) {
	log.Debugf("[TaskHandlerService] deleteRunner: taskId[%v]", taskId)
	svc.runners.Delete(taskId)
}

func (svc *ServiceV2) reportStatus() (err error) {
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

func (svc *ServiceV2) fetch() (tid primitive.ObjectID, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), svc.fetchTimeout)
	defer cancel()
	res, err := svc.c.TaskClient.Fetch(ctx, svc.c.NewRequest(nil))
	if err != nil {
		return tid, trace.TraceError(err)
	}
	if err := json.Unmarshal(res.Data, &tid); err != nil {
		return tid, trace.TraceError(err)
	}
	return tid, nil
}

func (svc *ServiceV2) run(taskId primitive.ObjectID) (err error) {
	// attempt to get runner from pool
	_, ok := svc.runners.Load(taskId)
	if ok {
		return trace.TraceError(errors2.ErrorTaskAlreadyExists)
	}

	// create a new task runner
	r, err := NewTaskRunnerV2(taskId, svc)
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

func newTaskHandlerServiceV2() (svc2 *ServiceV2, err error) {
	// service
	svc := &ServiceV2{
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
	svc.cfgSvc = nodeconfig.GetNodeConfigService()

	// grpc client
	svc.c = grpcclient.GetGrpcClientV2()

	log.Debugf("[NewTaskHandlerService] svc[cfgPath: %s]", svc.cfgSvc.GetConfigPath())

	return svc, nil
}

var _serviceV2 *ServiceV2
var _serviceV2Once = new(sync.Once)

func GetTaskHandlerServiceV2() (svr *ServiceV2, err error) {
	if _serviceV2 != nil {
		return _serviceV2, nil
	}
	_serviceV2Once.Do(func() {
		_serviceV2, err = newTaskHandlerServiceV2()
		if err != nil {
			log.Errorf("failed to create task handler service: %v", err)
		}
	})
	if err != nil {
		return nil, err
	}
	return _serviceV2, nil
}
