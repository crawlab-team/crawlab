package server

import (
	"context"
	"encoding/json"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/container"
	"github.com/crawlab-team/crawlab/core/entity"
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/delegate"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/core/notification"
	"github.com/crawlab-team/crawlab/core/utils"
	"github.com/crawlab-team/crawlab/db/mongo"
	grpc "github.com/crawlab-team/crawlab/grpc"
	"github.com/crawlab-team/crawlab/trace"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"io"
	"strings"
)

type TaskServer struct {
	grpc.UnimplementedTaskServiceServer

	// dependencies
	modelSvc service.ModelService
	cfgSvc   interfaces.NodeConfigService
	statsSvc interfaces.TaskStatsService

	// internals
	server interfaces.GrpcServer
}

// Subscribe to task stream when a task runner in a node starts
func (svr TaskServer) Subscribe(stream grpc.TaskService_SubscribeServer) (err error) {
	for {
		msg, err := stream.Recv()
		utils.LogDebug(msg.String())
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
		switch msg.Code {
		case grpc.StreamMessageCode_INSERT_DATA:
			err = svr.handleInsertData(msg)
		case grpc.StreamMessageCode_INSERT_LOGS:
			err = svr.handleInsertLogs(msg)
		default:
			err = errors.ErrorGrpcInvalidCode
			log.Errorf("invalid stream message code: %d", msg.Code)
			continue
		}
		if err != nil {
			log.Errorf("grpc error[%d]: %v", msg.Code, err)
		}
	}
}

// Fetch tasks to be executed by a task handler
func (svr TaskServer) Fetch(ctx context.Context, request *grpc.Request) (response *grpc.Response, err error) {
	nodeKey := request.GetNodeKey()
	if nodeKey == "" {
		return nil, trace.TraceError(errors.ErrorGrpcInvalidNodeKey)
	}
	n, err := svr.modelSvc.GetNodeByKey(nodeKey, nil)
	if err != nil {
		return nil, trace.TraceError(err)
	}
	var tid primitive.ObjectID
	opts := &mongo.FindOptions{
		Sort: bson.D{
			{"p", 1},
			{"_id", 1},
		},
		Limit: 1,
	}
	if err := mongo.RunTransactionWithContext(ctx, func(sc mongo2.SessionContext) (err error) {
		// get task queue item assigned to this node
		tid, err = svr.getTaskQueueItemIdAndDequeue(bson.M{"nid": n.Id}, opts, n.Id)
		if err != nil {
			return err
		}
		if !tid.IsZero() {
			return nil
		}

		// get task queue item assigned to any node (random mode)
		tid, err = svr.getTaskQueueItemIdAndDequeue(bson.M{"nid": nil}, opts, n.Id)
		if !tid.IsZero() {
			return nil
		}
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return HandleSuccessWithData(tid)
}

func (svr TaskServer) SendNotification(ctx context.Context, request *grpc.Request) (response *grpc.Response, err error) {
	svc := notification.GetService()
	var t = new(models.Task)
	if err := json.Unmarshal(request.Data, t); err != nil {
		return nil, trace.TraceError(err)
	}
	t, err = svr.modelSvc.GetTaskById(t.Id)
	if err != nil {
		return nil, trace.TraceError(err)
	}
	td, err := json.Marshal(t)
	if err != nil {
		return nil, trace.TraceError(err)
	}
	var e bson.M
	if err := json.Unmarshal(td, &e); err != nil {
		return nil, trace.TraceError(err)
	}
	ts, err := svr.modelSvc.GetTaskStatById(t.Id)
	if err != nil {
		return nil, trace.TraceError(err)
	}
	settings, _, err := svc.GetSettingList(bson.M{
		"enabled": true,
	}, nil, nil)
	if err != nil {
		return nil, trace.TraceError(err)
	}
	for _, s := range settings {
		switch s.TaskTrigger {
		case constants.NotificationTriggerTaskFinish:
			if t.Status != constants.TaskStatusPending && t.Status != constants.TaskStatusRunning {
				_ = svc.Send(s, e)
			}
		case constants.NotificationTriggerTaskError:
			if t.Status == constants.TaskStatusError || t.Status == constants.TaskStatusAbnormal {
				_ = svc.Send(s, e)
			}
		case constants.NotificationTriggerTaskEmptyResults:
			if t.Status != constants.TaskStatusPending && t.Status != constants.TaskStatusRunning {
				if ts.ResultCount == 0 {
					_ = svc.Send(s, e)
				}
			}
		case constants.NotificationTriggerTaskNever:
		}
	}
	return nil, nil
}

func (svr TaskServer) handleInsertData(msg *grpc.StreamMessage) (err error) {
	data, err := svr.deserialize(msg)
	if err != nil {
		return err
	}
	var records []interface{}
	for _, d := range data.Records {
		res, ok := d[constants.TaskKey]
		if ok {
			switch res.(type) {
			case string:
				id, err := primitive.ObjectIDFromHex(res.(string))
				if err == nil {
					d[constants.TaskKey] = id
				}
			}
		}
		records = append(records, d)
	}
	return svr.statsSvc.InsertData(data.TaskId, records...)
}

func (svr TaskServer) handleInsertLogs(msg *grpc.StreamMessage) (err error) {
	data, err := svr.deserialize(msg)
	if err != nil {
		return err
	}
	return svr.statsSvc.InsertLogs(data.TaskId, data.Logs...)
}

func (svr TaskServer) getTaskQueueItemIdAndDequeue(query bson.M, opts *mongo.FindOptions, nid primitive.ObjectID) (tid primitive.ObjectID, err error) {
	var tq models.TaskQueueItem
	if err := mongo.GetMongoCol(interfaces.ModelColNameTaskQueue).Find(query, opts).One(&tq); err != nil {
		if err == mongo2.ErrNoDocuments {
			return tid, nil
		}
		return tid, trace.TraceError(err)
	}
	t, err := svr.modelSvc.GetTaskById(tq.Id)
	if err == nil {
		t.NodeId = nid
		_ = delegate.NewModelDelegate(t).Save()
	}
	_ = delegate.NewModelDelegate(&tq).Delete()
	return tq.Id, nil
}

func (svr TaskServer) deserialize(msg *grpc.StreamMessage) (data entity.StreamMessageTaskData, err error) {
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		return data, trace.TraceError(err)
	}
	if data.TaskId.IsZero() {
		return data, trace.TraceError(errors.ErrorGrpcInvalidType)
	}
	return data, nil
}

func NewTaskServer() (res *TaskServer, err error) {
	// task server
	svr := &TaskServer{}

	// dependency injection
	if err := container.GetContainer().Invoke(func(
		modelSvc service.ModelService,
		statsSvc interfaces.TaskStatsService,
		cfgSvc interfaces.NodeConfigService,
	) {
		svr.modelSvc = modelSvc
		svr.statsSvc = statsSvc
		svr.cfgSvc = cfgSvc
	}); err != nil {
		return nil, err
	}

	return svr, nil
}
