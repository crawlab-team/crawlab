package server

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/entity"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/models/service"
	nodeconfig "github.com/crawlab-team/crawlab/core/node/config"
	"github.com/crawlab-team/crawlab/core/notification"
	"github.com/crawlab-team/crawlab/core/task/stats"
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

type TaskServerV2 struct {
	grpc.UnimplementedTaskServiceServer

	// dependencies
	cfgSvc   interfaces.NodeConfigService
	statsSvc *stats.ServiceV2

	// internals
	server interfaces.GrpcServer
}

// Subscribe to task stream when a task runner in a node starts
func (svr TaskServerV2) Subscribe(stream grpc.TaskService_SubscribeServer) (err error) {
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
			err = errors.New("invalid stream message code")
			log.Errorf("invalid stream message code: %d", msg.Code)
			continue
		}
		if err != nil {
			log.Errorf("grpc error[%d]: %v", msg.Code, err)
		}
	}
}

// Fetch tasks to be executed by a task handler
func (svr TaskServerV2) Fetch(ctx context.Context, request *grpc.Request) (response *grpc.Response, err error) {
	nodeKey := request.GetNodeKey()
	if nodeKey == "" {
		return nil, errors.New("invalid node key")
	}
	n, err := service.NewModelServiceV2[models.NodeV2]().GetOne(bson.M{"key": nodeKey}, nil)
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

func (svr TaskServerV2) SendNotification(ctx context.Context, request *grpc.Request) (response *grpc.Response, err error) {
	svc := notification.GetServiceV2()
	var t = new(models.TaskV2)
	if err := json.Unmarshal(request.Data, t); err != nil {
		return nil, trace.TraceError(err)
	}
	t, err = service.NewModelServiceV2[models.TaskV2]().GetById(t.Id)
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
	ts, err := service.NewModelServiceV2[models.TaskStatV2]().GetById(t.Id)
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
				_ = svc.Send(&s, e)
			}
		case constants.NotificationTriggerTaskError:
			if t.Status == constants.TaskStatusError || t.Status == constants.TaskStatusAbnormal {
				_ = svc.Send(&s, e)
			}
		case constants.NotificationTriggerTaskEmptyResults:
			if t.Status != constants.TaskStatusPending && t.Status != constants.TaskStatusRunning {
				if ts.ResultCount == 0 {
					_ = svc.Send(&s, e)
				}
			}
		case constants.NotificationTriggerTaskNever:
		}
	}
	return nil, nil
}

func (svr TaskServerV2) handleInsertData(msg *grpc.StreamMessage) (err error) {
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

func (svr TaskServerV2) handleInsertLogs(msg *grpc.StreamMessage) (err error) {
	data, err := svr.deserialize(msg)
	if err != nil {
		return err
	}
	return svr.statsSvc.InsertLogs(data.TaskId, data.Logs...)
}

func (svr TaskServerV2) getTaskQueueItemIdAndDequeue(query bson.M, opts *mongo.FindOptions, nid primitive.ObjectID) (tid primitive.ObjectID, err error) {
	tq, err := service.NewModelServiceV2[models.TaskQueueItemV2]().GetOne(query, opts)
	if err != nil {
		if errors.Is(err, mongo2.ErrNoDocuments) {
			return tid, nil
		}
		return tid, trace.TraceError(err)
	}
	t, err := service.NewModelServiceV2[models.TaskV2]().GetById(tq.Id)
	if err == nil {
		t.NodeId = nid
		err = service.NewModelServiceV2[models.TaskV2]().ReplaceById(t.Id, *t)
		if err != nil {
			return tid, trace.TraceError(err)
		}
	}
	err = service.NewModelServiceV2[models.TaskQueueItemV2]().DeleteById(tq.Id)
	if err != nil {
		return tid, trace.TraceError(err)
	}
	return tq.Id, nil
}

func (svr TaskServerV2) deserialize(msg *grpc.StreamMessage) (data entity.StreamMessageTaskData, err error) {
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		return data, trace.TraceError(err)
	}
	if data.TaskId.IsZero() {
		return data, errors.New("invalid task id")
	}
	return data, nil
}

func NewTaskServerV2() (res *TaskServerV2, err error) {
	// task server
	svr := &TaskServerV2{}

	svr.cfgSvc = nodeconfig.GetNodeConfigService()

	svr.statsSvc, err = stats.GetTaskStatsServiceV2()
	if err != nil {
		return nil, err
	}

	return svr, nil
}
