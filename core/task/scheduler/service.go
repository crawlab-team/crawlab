package scheduler

import (
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/container"
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/delegate"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/core/task"
	"github.com/crawlab-team/crawlab/db/mongo"
	grpc "github.com/crawlab-team/crawlab/grpc"
	"github.com/crawlab-team/crawlab/trace"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Service struct {
	// dependencies
	interfaces.TaskBaseService
	nodeCfgSvc interfaces.NodeConfigService
	modelSvc   service.ModelService
	svr        interfaces.GrpcServer
	handlerSvc interfaces.TaskHandlerService

	// settings
	interval time.Duration
}

func (svc *Service) Start() {
	go svc.initTaskStatus()
	go svc.cleanupTasks()
	svc.Wait()
	svc.Stop()
}

func (svc *Service) Enqueue(t interfaces.Task) (t2 interfaces.Task, err error) {
	// set task status
	t.SetStatus(constants.TaskStatusPending)

	// user
	var u *models.User
	if !t.GetUserId().IsZero() {
		u, _ = svc.modelSvc.GetUserById(t.GetUserId())
	}

	// add task
	if err = delegate.NewModelDelegate(t, u).Add(); err != nil {
		return nil, err
	}

	// task queue item
	tq := &models.TaskQueueItem{
		Id:       t.GetId(),
		Priority: t.GetPriority(),
		NodeId:   t.GetNodeId(),
	}

	// task stat
	ts := &models.TaskStat{
		Id:       t.GetId(),
		CreateTs: time.Now(),
	}

	// enqueue task
	_, err = mongo.GetMongoCol(interfaces.ModelColNameTaskQueue).Insert(tq)
	if err != nil {
		return nil, trace.TraceError(err)
	}

	// add task stat
	_, err = mongo.GetMongoCol(interfaces.ModelColNameTaskStat).Insert(ts)
	if err != nil {
		return nil, trace.TraceError(err)
	}

	// success
	return t, nil
}

func (svc *Service) Cancel(id primitive.ObjectID, args ...interface{}) (err error) {
	// task
	t, err := svc.modelSvc.GetTaskById(id)
	if err != nil {
		return trace.TraceError(err)
	}

	// initial status
	initialStatus := t.Status

	// set task status as "cancelled"
	_ = svc.SaveTask(t, constants.TaskStatusCancelled)

	// set status of pending tasks as "cancelled" and remove from task item queue
	if initialStatus == constants.TaskStatusPending {
		// remove from task item queue
		if err := mongo.GetMongoCol(interfaces.ModelColNameTaskQueue).DeleteId(t.GetId()); err != nil {
			return trace.TraceError(err)
		}
		return nil
	}

	// whether task is running on master node
	isMasterTask, err := svc.isMasterNode(t)
	if err != nil {
		// when error, force status being set as "cancelled"
		return svc.SaveTask(t, constants.TaskStatusCancelled)
	}

	// node
	n, err := svc.modelSvc.GetNodeById(t.GetNodeId())
	if err != nil {
		return trace.TraceError(err)
	}

	if isMasterTask {
		// cancel task on master
		if err := svc.handlerSvc.Cancel(id); err != nil {
			return trace.TraceError(err)
		}
		// cancel success
		return nil
	} else {
		// send to cancel task on worker nodes
		if err := svc.svr.SendStreamMessageWithData("node:"+n.GetKey(), grpc.StreamMessageCode_CANCEL_TASK, t); err != nil {
			return trace.TraceError(err)
		}
		// cancel success
		return nil
	}
}

func (svc *Service) SetInterval(interval time.Duration) {
	svc.interval = interval
}

// initTaskStatus initialize task status of existing tasks
func (svc *Service) initTaskStatus() {
	// set status of running tasks as TaskStatusAbnormal
	runningTasks, err := svc.modelSvc.GetTaskList(bson.M{
		"status": bson.M{
			"$in": []string{
				constants.TaskStatusPending,
				constants.TaskStatusRunning,
			},
		},
	}, nil)
	if err != nil {
		if err == mongo2.ErrNoDocuments {
			return
		}
		trace.PrintError(err)
	}
	for _, t := range runningTasks {
		go func(t *models.Task) {
			if err := svc.SaveTask(t, constants.TaskStatusAbnormal); err != nil {
				trace.PrintError(err)
			}
		}(&t)
	}
	if err := svc.modelSvc.GetBaseService(interfaces.ModelIdTaskQueue).DeleteList(nil); err != nil {
		return
	}
}

func (svc *Service) isMasterNode(t *models.Task) (ok bool, err error) {
	if t.GetNodeId().IsZero() {
		return false, trace.TraceError(errors.ErrorTaskNoNodeId)
	}
	n, err := svc.modelSvc.GetNodeById(t.GetNodeId())
	if err != nil {
		if err == mongo2.ErrNoDocuments {
			return false, trace.TraceError(errors.ErrorTaskNodeNotFound)
		}
		return false, trace.TraceError(err)
	}
	return n.IsMaster, nil
}

func (svc *Service) cleanupTasks() {
	for {
		// task stats over 30 days ago
		taskStats, err := svc.modelSvc.GetTaskStatList(bson.M{
			"create_ts": bson.M{
				"$lt": time.Now().Add(-30 * 24 * time.Hour),
			},
		}, nil)
		if err != nil {
			time.Sleep(30 * time.Minute)
			continue
		}

		// task ids
		var ids []primitive.ObjectID
		for _, ts := range taskStats {
			ids = append(ids, ts.Id)
		}

		if len(ids) > 0 {
			// remove tasks
			if err := svc.modelSvc.GetBaseService(interfaces.ModelIdTask).DeleteList(bson.M{
				"_id": bson.M{"$in": ids},
			}); err != nil {
				trace.PrintError(err)
			}

			// remove task stats
			if err := svc.modelSvc.GetBaseService(interfaces.ModelIdTaskStat).DeleteList(bson.M{
				"_id": bson.M{"$in": ids},
			}); err != nil {
				trace.PrintError(err)
			}
		}

		time.Sleep(30 * time.Minute)
	}
}

func NewTaskSchedulerService() (svc2 interfaces.TaskSchedulerService, err error) {
	// base service
	baseSvc, err := task.NewBaseService()
	if err != nil {
		return nil, trace.TraceError(err)
	}

	// service
	svc := &Service{
		TaskBaseService: baseSvc,
		interval:        5 * time.Second,
	}

	// dependency injection
	if err := container.GetContainer().Invoke(func(
		nodeCfgSvc interfaces.NodeConfigService,
		modelSvc service.ModelService,
		svr interfaces.GrpcServer,
		handlerSvc interfaces.TaskHandlerService,
	) {
		svc.nodeCfgSvc = nodeCfgSvc
		svc.modelSvc = modelSvc
		svc.svr = svr
		svc.handlerSvc = handlerSvc
	}); err != nil {
		return nil, err
	}

	return svc, nil
}

var svc interfaces.TaskSchedulerService

func GetTaskSchedulerService() (svr interfaces.TaskSchedulerService, err error) {
	if svc != nil {
		return svc, nil
	}
	svc, err = NewTaskSchedulerService()
	if err != nil {
		return nil, err
	}
	return svc, nil
}
