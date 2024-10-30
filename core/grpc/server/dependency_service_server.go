package server

import (
	"context"
	"errors"
	"github.com/apex/log"
	models2 "github.com/crawlab-team/crawlab/core/models/models/v2"
	"github.com/crawlab-team/crawlab/core/models/service"
	mongo2 "github.com/crawlab-team/crawlab/db/mongo"
	"github.com/crawlab-team/crawlab/grpc"
	"github.com/crawlab-team/crawlab/trace"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"sync"
	"time"
)

type DependencyServiceServer struct {
	grpc.UnimplementedDependencyServiceV2Server
	mu      *sync.Mutex
	streams map[string]*grpc.DependencyServiceV2_ConnectServer
}

func (svr DependencyServiceServer) Connect(req *grpc.DependencyServiceV2ConnectRequest, stream grpc.DependencyServiceV2_ConnectServer) (err error) {
	svr.mu.Lock()
	svr.streams[req.NodeKey] = &stream
	svr.mu.Unlock()
	log.Info("[DependencyServiceServer] connected: " + req.NodeKey)

	// Keep this scope alive because once this scope exits - the stream is closed
	for {
		select {
		case <-stream.Context().Done():
			log.Info("[DependencyServiceServer] disconnected: " + req.NodeKey)
			return nil
		}
	}
}

func (svr DependencyServiceServer) Sync(ctx context.Context, request *grpc.DependencyServiceV2SyncRequest) (response *grpc.Response, err error) {
	n, err := service.NewModelServiceV2[models2.NodeV2]().GetOne(bson.M{"key": request.NodeKey}, nil)
	if err != nil {
		return nil, err
	}

	depsDb, err := service.NewModelServiceV2[models2.DependencyV2]().GetMany(bson.M{
		"node_id": n.Id,
		"type":    request.Lang,
	}, nil)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			log.Errorf("[DependencyServiceV2] get dependencies from db error: %v", err)
			return nil, err
		}
	}

	depsDbMap := make(map[string]*models2.DependencyV2)
	for _, d := range depsDb {
		depsDbMap[d.Name] = &d
	}

	var depsToInsert []models2.DependencyV2
	depsMap := make(map[string]*models2.DependencyV2)
	for _, dep := range request.Dependencies {
		d := models2.DependencyV2{
			Name:    dep.Name,
			NodeId:  n.Id,
			Type:    request.Lang,
			Version: dep.Version,
		}
		d.SetCreatedAt(time.Now())

		depsMap[d.Name] = &d

		_, ok := depsDbMap[d.Name]
		if !ok {
			depsToInsert = append(depsToInsert, d)
		}
	}

	var depIdsToDelete []primitive.ObjectID
	for _, d := range depsDb {
		_, ok := depsMap[d.Name]
		if !ok {
			depIdsToDelete = append(depIdsToDelete, d.Id)
		}
	}

	err = mongo2.RunTransaction(func(ctx mongo.SessionContext) (err error) {
		if len(depIdsToDelete) > 0 {
			err = service.NewModelServiceV2[models2.DependencyV2]().DeleteMany(bson.M{
				"_id": bson.M{"$in": depIdsToDelete},
			})
			if err != nil {
				log.Errorf("[DependencyServiceServer] delete dependencies in db error: %v", err)
				trace.PrintError(err)
				return err
			}
		}

		if len(depsToInsert) > 0 {
			_, err = service.NewModelServiceV2[models2.DependencyV2]().InsertMany(depsToInsert)
			if err != nil {
				log.Errorf("[DependencyServiceServer] insert dependencies in db error: %v", err)
				trace.PrintError(err)
				return err
			}
		}

		return nil
	})

	return nil, err
}

func (svr DependencyServiceServer) UpdateTaskLog(stream grpc.DependencyServiceV2_UpdateTaskLogServer) (err error) {
	var t *models2.DependencyTaskV2
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
			t, err = service.NewModelServiceV2[models2.DependencyTaskV2]().GetById(taskId)
			if err != nil {
				return err
			}
		}
		var logs []models2.DependencyLogV2
		for _, line := range req.LogLines {
			l := models2.DependencyLogV2{
				TaskId:  taskId,
				Content: line,
			}
			l.SetCreated(t.CreatedBy)
			logs = append(logs, l)
		}
		_, err = service.NewModelServiceV2[models2.DependencyLogV2]().InsertMany(logs)
		if err != nil {
			return err
		}
	}
}

func (svr DependencyServiceServer) GetStream(key string) (stream *grpc.DependencyServiceV2_ConnectServer, err error) {
	svr.mu.Lock()
	defer svr.mu.Unlock()
	stream, ok := svr.streams[key]
	if !ok {
		return nil, errors.New("stream not found")
	}
	return stream, nil
}

func NewDependencyServerV2() *DependencyServiceServer {
	return &DependencyServiceServer{
		mu:      new(sync.Mutex),
		streams: make(map[string]*grpc.DependencyServiceV2_ConnectServer),
	}
}

var depSvc *DependencyServiceServer

func GetDependencyServerV2() *DependencyServiceServer {
	if depSvc != nil {
		return depSvc
	}
	depSvc = NewDependencyServerV2()
	return depSvc
}
