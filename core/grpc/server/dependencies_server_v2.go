package server

import (
	"context"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/grpc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"sync"
)

type DependenciesServerV2 struct {
	grpc.UnimplementedDependenciesServiceV2Server
	mu      *sync.Mutex
	streams map[string]grpc.DependenciesServiceV2_ConnectServer
}

func (svr DependenciesServerV2) Connect(stream grpc.DependenciesServiceV2_ConnectServer) (err error) {
	svr.mu.Lock()
	defer svr.mu.Unlock()
	req, err := stream.Recv()
	if err != nil {
		return err
	}
	svr.streams[req.NodeKey] = stream
	return nil
}

func (svr DependenciesServerV2) Sync(ctx context.Context, request *grpc.DependenciesServiceV2SyncRequest) (response *grpc.Response, err error) {
	n, err := service.NewModelServiceV2[models.NodeV2]().GetOne(bson.M{"key": request.NodeKey}, nil)
	if err != nil {
		return nil, err
	}
	var deps []models.DependencyV2
	for _, dep := range request.Dependencies {
		deps = append(deps, models.DependencyV2{
			Name:    dep.Name,
			NodeId:  n.Id,
			Type:    request.Lang,
			Version: dep.Version,
		})
	}
	_, err = service.NewModelServiceV2[models.DependencyV2]().InsertMany(deps)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (svr DependenciesServerV2) UpdateTaskLog(stream grpc.DependenciesServiceV2_UpdateTaskLogServer) (err error) {
	var t *models.DependencyTaskV2
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// all messages have been received
			return stream.SendAndClose(&grpc.Response{Message: "update task log finished"})
		}
		if err != nil {
			return err
		}
		taskId, err := primitive.ObjectIDFromHex(req.TaskId)
		if err != nil {
			return err
		}
		if t == nil {
			t, err = service.NewModelServiceV2[models.DependencyTaskV2]().GetById(taskId)
			if err != nil {
				return err
			}
		}
		l := models.DependencyLogV2{
			TaskId:  taskId,
			Content: req.Content,
		}
		l.SetCreated(t.CreatedBy)
		_, err = service.NewModelServiceV2[models.DependencyLogV2]().InsertOne(l)
		if err != nil {
			return err
		}
	}
}

func NewDependenciesServerV2() *DependenciesServerV2 {
	return &DependenciesServerV2{}
}
