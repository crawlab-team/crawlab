package server

import (
	"context"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/models/v2"
	"github.com/crawlab-team/crawlab/core/models/service"
	nodeconfig "github.com/crawlab-team/crawlab/core/node/config"
	"github.com/crawlab-team/crawlab/core/notification"
	"github.com/crawlab-team/crawlab/core/utils"
	"github.com/crawlab-team/crawlab/grpc"
	errors2 "github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
	"time"
)

var nodeServiceMutex = sync.Mutex{}

type NodeServiceServer struct {
	grpc.UnimplementedNodeServiceServer

	// dependencies
	cfgSvc interfaces.NodeConfigService

	// internals
	subs map[primitive.ObjectID]grpc.NodeService_SubscribeServer
}

// Register from handler/worker to master
func (svr NodeServiceServer) Register(_ context.Context, req *grpc.NodeServiceRegisterRequest) (res *grpc.Response, err error) {
	// node key
	if req.NodeKey == "" {
		return HandleError(errors.ErrorModelMissingRequiredData)
	}

	// find in db
	var node *models.NodeV2
	node, err = service.NewModelServiceV2[models.NodeV2]().GetOne(bson.M{"key": req.NodeKey}, nil)
	if err == nil {
		// register existing
		node.Status = constants.NodeStatusOnline
		node.Active = true
		node.ActiveAt = time.Now()
		err = service.NewModelServiceV2[models.NodeV2]().ReplaceById(node.Id, *node)
		if err != nil {
			return HandleError(err)
		}
		log.Infof("[NodeServiceServer] updated worker[%s] in db. id: %s", req.NodeKey, node.Id.Hex())
	} else if errors2.Is(err, mongo.ErrNoDocuments) {
		// register new
		node = &models.NodeV2{
			Key:        req.NodeKey,
			Name:       req.NodeName,
			Status:     constants.NodeStatusOnline,
			Active:     true,
			ActiveAt:   time.Now(),
			Enabled:    true,
			MaxRunners: int(req.MaxRunners),
		}
		node.SetCreated(primitive.NilObjectID)
		node.SetUpdated(primitive.NilObjectID)
		node.Id, err = service.NewModelServiceV2[models.NodeV2]().InsertOne(*node)
		if err != nil {
			return HandleError(err)
		}
		log.Infof("[NodeServiceServer] added worker[%s] in db. id: %s", req.NodeKey, node.Id.Hex())
	} else {
		// error
		return HandleError(err)
	}

	log.Infof("[NodeServiceServer] master registered worker[%s]", req.NodeKey)

	return HandleSuccessWithData(node)
}

// SendHeartbeat from worker to master
func (svr NodeServiceServer) SendHeartbeat(_ context.Context, req *grpc.NodeServiceSendHeartbeatRequest) (res *grpc.Response, err error) {
	// find in db
	node, err := service.NewModelServiceV2[models.NodeV2]().GetOne(bson.M{"key": req.NodeKey}, nil)
	if err != nil {
		if errors2.Is(err, mongo.ErrNoDocuments) {
			return HandleError(errors.ErrorNodeNotExists)
		}
		return HandleError(err)
	}
	oldStatus := node.Status

	// update status
	node.Status = constants.NodeStatusOnline
	node.Active = true
	node.ActiveAt = time.Now()
	err = service.NewModelServiceV2[models.NodeV2]().ReplaceById(node.Id, *node)
	if err != nil {
		return HandleError(err)
	}
	newStatus := node.Status

	// send notification if status changed
	if utils.IsPro() {
		if oldStatus != newStatus {
			go notification.GetNotificationServiceV2().SendNodeNotification(node)
		}
	}

	return HandleSuccessWithData(node)
}

func (svr NodeServiceServer) Subscribe(request *grpc.NodeServiceSubscribeRequest, stream grpc.NodeService_SubscribeServer) (err error) {
	log.Infof("[NodeServiceServer] master received subscribe request from node[%s]", request.NodeKey)

	// find in db
	node, err := service.NewModelServiceV2[models.NodeV2]().GetOne(bson.M{"key": request.NodeKey}, nil)
	if err != nil {
		log.Errorf("[NodeServiceServer] error getting node: %v", err)
		return err
	}

	// subscribe
	nodeServiceMutex.Lock()
	svr.subs[node.Id] = stream
	nodeServiceMutex.Unlock()

	// TODO: send notification

	// wait for stream to close
	<-stream.Context().Done()

	// unsubscribe
	nodeServiceMutex.Lock()
	delete(svr.subs, node.Id)
	nodeServiceMutex.Unlock()
	log.Infof("[NodeServiceServer] master unsubscribed from node[%s]", request.NodeKey)

	return nil
}

func (svr NodeServiceServer) GetSubscribeStream(nodeId primitive.ObjectID) (stream grpc.NodeService_SubscribeServer, ok bool) {
	nodeServiceMutex.Lock()
	defer nodeServiceMutex.Unlock()
	stream, ok = svr.subs[nodeId]
	return stream, ok
}

var nodeSvr *NodeServiceServer
var nodeSvrOnce = new(sync.Once)

func NewNodeServiceServer() (res *NodeServiceServer, err error) {
	if nodeSvr != nil {
		return nodeSvr, nil
	}
	nodeSvrOnce.Do(func() {
		nodeSvr = &NodeServiceServer{
			subs: make(map[primitive.ObjectID]grpc.NodeService_SubscribeServer),
		}
		nodeSvr.cfgSvc = nodeconfig.GetNodeConfigService()
		if err != nil {
			log.Errorf("[NodeServiceServer] error: %s", err.Error())
		}
	})
	if err != nil {
		return nil, err
	}
	return nodeSvr, nil
}
