package server

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/interfaces"
	models2 "github.com/crawlab-team/crawlab/core/models/models/v2"
	"github.com/crawlab-team/crawlab/core/models/service"
	nodeconfig "github.com/crawlab-team/crawlab/core/node/config"
	"github.com/crawlab-team/crawlab/core/notification"
	"github.com/crawlab-team/crawlab/core/task/stats"
	"github.com/crawlab-team/crawlab/core/utils"
	"github.com/crawlab-team/crawlab/db/mongo"
	"github.com/crawlab-team/crawlab/grpc"
	"github.com/crawlab-team/crawlab/trace"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"io"
	"strings"
	"sync"
)

var taskServiceMutex = sync.Mutex{}

type TaskServiceServer struct {
	grpc.UnimplementedTaskServiceServer

	// dependencies
	cfgSvc   interfaces.NodeConfigService
	statsSvc *stats.Service

	// internals
	subs map[primitive.ObjectID]grpc.TaskService_SubscribeServer
}

func (svr TaskServiceServer) Subscribe(req *grpc.TaskServiceSubscribeRequest, stream grpc.TaskService_SubscribeServer) (err error) {
	// task id
	taskId, err := primitive.ObjectIDFromHex(req.TaskId)
	if err != nil {
		return errors.New("invalid task id")
	}

	// validate stream
	if stream == nil {
		return errors.New("invalid stream")
	}

	// add stream
	taskServiceMutex.Lock()
	svr.subs[taskId] = stream
	taskServiceMutex.Unlock()

	// wait for stream to close
	<-stream.Context().Done()

	// remove stream
	taskServiceMutex.Lock()
	delete(svr.subs, taskId)
	taskServiceMutex.Unlock()
	log.Infof("[TaskServiceServer] task stream closed: %s", taskId.Hex())

	return nil
}

// Connect to task stream when a task runner in a node starts
func (svr TaskServiceServer) Connect(stream grpc.TaskService_ConnectServer) (err error) {
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			if strings.HasSuffix(err.Error(), "context canceled") {
				return nil
			}
			trace.PrintError(err)
			continue
		}

		// validate task id
		taskId, err := primitive.ObjectIDFromHex(msg.TaskId)
		if err != nil {
			log.Errorf("invalid task id: %s", msg.TaskId)
			continue
		}

		switch msg.Code {
		case grpc.TaskServiceConnectCode_INSERT_DATA:
			err = svr.handleInsertData(taskId, msg)
		case grpc.TaskServiceConnectCode_INSERT_LOGS:
			err = svr.handleInsertLogs(taskId, msg)
		default:
			err = errors.New("invalid stream message code")
			log.Errorf("invalid stream message code: %d", msg.Code)
			continue
		}
		if err != nil {
			log.Errorf("grpc error[%d]: %v", msg.Code, err)
		}
	}
}

// FetchTask tasks to be executed by a task handler
func (svr TaskServiceServer) FetchTask(ctx context.Context, request *grpc.TaskServiceFetchTaskRequest) (response *grpc.TaskServiceFetchTaskResponse, err error) {
	nodeKey := request.GetNodeKey()
	if nodeKey == "" {
		return nil, errors.New("invalid node key")
	}
	n, err := service.NewModelServiceV2[models2.NodeV2]().GetOne(bson.M{"key": nodeKey}, nil)
	if err != nil {
		return nil, trace.TraceError(err)
	}
	var tid primitive.ObjectID
	opts := &mongo.FindOptions{
		Sort: bson.D{
			{"priority", 1},
			{"_id", 1},
		},
		Limit: 1,
	}
	if err := mongo.RunTransactionWithContext(ctx, func(sc mongo2.SessionContext) (err error) {
		// fetch task for the given node
		t, err := service.NewModelServiceV2[models2.TaskV2]().GetOne(bson.M{
			"node_id": n.Id,
			"status":  constants.TaskStatusPending,
		}, opts)
		if err == nil {
			tid = t.Id
			t.Status = constants.TaskStatusAssigned
			return svr.saveTask(t)
		} else if !errors.Is(err, mongo2.ErrNoDocuments) {
			log.Errorf("error fetching task for node[%s]: %v", nodeKey, err)
			return err
		}

		// fetch task for any node
		t, err = service.NewModelServiceV2[models2.TaskV2]().GetOne(bson.M{
			"node_id": primitive.NilObjectID,
			"status":  constants.TaskStatusPending,
		}, opts)
		if err == nil {
			tid = t.Id
			t.NodeId = n.Id
			t.Status = constants.TaskStatusAssigned
			return svr.saveTask(t)
		} else if !errors.Is(err, mongo2.ErrNoDocuments) {
			log.Errorf("error fetching task for any node: %v", err)
			return err
		}

		// no task found
		return nil
	}); err != nil {
		return nil, err
	}

	return &grpc.TaskServiceFetchTaskResponse{TaskId: tid.Hex()}, nil
}

func (svr TaskServiceServer) SendNotification(_ context.Context, request *grpc.TaskServiceSendNotificationRequest) (response *grpc.Response, err error) {
	if !utils.IsPro() {
		return nil, nil
	}

	// task id
	taskId, err := primitive.ObjectIDFromHex(request.TaskId)
	if err != nil {
		log.Errorf("invalid task id: %s", request.TaskId)
		return nil, trace.TraceError(err)
	}

	// arguments
	var args []any

	// task
	task, err := service.NewModelServiceV2[models2.TaskV2]().GetById(taskId)
	if err != nil {
		log.Errorf("task not found: %s", request.TaskId)
		return nil, trace.TraceError(err)
	}
	args = append(args, task)

	// task stat
	taskStat, err := service.NewModelServiceV2[models2.TaskStatV2]().GetById(task.Id)
	if err != nil {
		log.Errorf("task stat not found for task: %s", request.TaskId)
		return nil, trace.TraceError(err)
	}
	args = append(args, taskStat)

	// spider
	spider, err := service.NewModelServiceV2[models2.SpiderV2]().GetById(task.SpiderId)
	if err != nil {
		log.Errorf("spider not found for task: %s", request.TaskId)
		return nil, trace.TraceError(err)
	}
	args = append(args, spider)

	// node
	node, err := service.NewModelServiceV2[models2.NodeV2]().GetById(task.NodeId)
	if err != nil {
		return nil, trace.TraceError(err)
	}
	args = append(args, node)

	// schedule
	var schedule *models2.ScheduleV2
	if !task.ScheduleId.IsZero() {
		schedule, err = service.NewModelServiceV2[models2.ScheduleV2]().GetById(task.ScheduleId)
		if err != nil {
			log.Errorf("schedule not found for task: %s", request.TaskId)
			return nil, trace.TraceError(err)
		}
		args = append(args, schedule)
	}

	// settings
	settings, err := service.NewModelServiceV2[models2.NotificationSettingV2]().GetMany(bson.M{
		"enabled": true,
		"trigger": bson.M{
			"$regex": constants.NotificationTriggerPatternTask,
		},
	}, nil)
	if err != nil {
		return nil, trace.TraceError(err)
	}

	// notification service
	svc := notification.GetNotificationServiceV2()

	for _, s := range settings {
		// compatible with old settings
		trigger := s.Trigger
		if trigger == "" {
			trigger = s.TaskTrigger
		}

		// send notification
		switch trigger {
		case constants.NotificationTriggerTaskFinish:
			if task.Status != constants.TaskStatusPending && task.Status != constants.TaskStatusRunning {
				go svc.Send(&s, args...)
			}
		case constants.NotificationTriggerTaskError:
			if task.Status == constants.TaskStatusError || task.Status == constants.TaskStatusAbnormal {
				go svc.Send(&s, args...)
			}
		case constants.NotificationTriggerTaskEmptyResults:
			if task.Status != constants.TaskStatusPending && task.Status != constants.TaskStatusRunning {
				if taskStat.ResultCount == 0 {
					go svc.Send(&s, args...)
				}
			}
		}
	}

	return nil, nil
}

func (svr TaskServiceServer) GetSubscribeStream(taskId primitive.ObjectID) (stream grpc.TaskService_SubscribeServer, ok bool) {
	taskServiceMutex.Lock()
	defer taskServiceMutex.Unlock()
	stream, ok = svr.subs[taskId]
	return stream, ok
}

func (svr TaskServiceServer) handleInsertData(taskId primitive.ObjectID, msg *grpc.TaskServiceConnectRequest) (err error) {
	var records []map[string]interface{}
	err = json.Unmarshal(msg.Data, &records)
	if err != nil {
		return trace.TraceError(err)
	}
	return svr.statsSvc.InsertData(taskId, records...)
}

func (svr TaskServiceServer) handleInsertLogs(taskId primitive.ObjectID, msg *grpc.TaskServiceConnectRequest) (err error) {
	var logs []string
	err = json.Unmarshal(msg.Data, &logs)
	if err != nil {
		return trace.TraceError(err)
	}
	return svr.statsSvc.InsertLogs(taskId, logs...)
}

func (svr TaskServiceServer) saveTask(t *models2.TaskV2) (err error) {
	t.SetUpdated(t.CreatedBy)
	return service.NewModelServiceV2[models2.TaskV2]().ReplaceById(t.Id, *t)
}

func NewTaskServiceServer() (res *TaskServiceServer, err error) {
	// task server
	svr := &TaskServiceServer{
		subs: make(map[primitive.ObjectID]grpc.TaskService_SubscribeServer),
	}

	svr.cfgSvc = nodeconfig.GetNodeConfigService()

	svr.statsSvc, err = stats.GetTaskStatsServiceV2()
	if err != nil {
		return nil, err
	}

	return svr, nil
}
