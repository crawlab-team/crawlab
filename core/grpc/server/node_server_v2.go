package server

import (
	"context"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/entity"
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

type NodeServerV2 struct {
	grpc.UnimplementedNodeServiceServer

	// dependencies
	cfgSvc interfaces.NodeConfigService

	// internals
	server *GrpcServerV2
}

// Register from handler/worker to master
func (svr NodeServerV2) Register(_ context.Context, req *grpc.NodeServiceRegisterRequest) (res *grpc.Response, err error) {
	// unmarshall data
	if req.IsMaster {
		// error: cannot register master node
		return HandleError(errors.ErrorGrpcNotAllowed)
	}

	// node key
	if req.Key == "" {
		return HandleError(errors.ErrorModelMissingRequiredData)
	}

	// find in db
	var node *models.NodeV2
	node, err = service.NewModelServiceV2[models.NodeV2]().GetOne(bson.M{"key": req.Key}, nil)
	if err == nil {
		// register existing
		node.Status = constants.NodeStatusRegistered
		node.Active = true
		node.ActiveAt = time.Now()
		err = service.NewModelServiceV2[models.NodeV2]().ReplaceById(node.Id, *node)
		if err != nil {
			return HandleError(err)
		}
		log.Infof("[NodeServerV2] updated worker[%s] in db. id: %s", req.Key, node.Id.Hex())
	} else if errors2.Is(err, mongo.ErrNoDocuments) {
		// register new
		node = &models.NodeV2{
			Key:        req.Key,
			Name:       req.Name,
			Status:     constants.NodeStatusRegistered,
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
		log.Infof("[NodeServerV2] added worker[%s] in db. id: %s", req.Key, node.Id.Hex())
	} else {
		// error
		return HandleError(err)
	}

	log.Infof("[NodeServerV2] master registered worker[%s]", req.Key)

	return HandleSuccessWithData(node)
}

// SendHeartbeat from worker to master
func (svr NodeServerV2) SendHeartbeat(_ context.Context, req *grpc.NodeServiceSendHeartbeatRequest) (res *grpc.Response, err error) {
	// find in db
	node, err := service.NewModelServiceV2[models.NodeV2]().GetOne(bson.M{"key": req.Key}, nil)
	if err != nil {
		if errors2.Is(err, mongo.ErrNoDocuments) {
			return HandleError(errors.ErrorNodeNotExists)
		}
		return HandleError(err)
	}
	oldStatus := node.Status

	// validate status
	if node.Status == constants.NodeStatusUnregistered {
		return HandleError(errors.ErrorNodeUnregistered)
	}

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

func (svr NodeServerV2) Subscribe(request *grpc.Request, stream grpc.NodeService_SubscribeServer) (err error) {
	log.Infof("[NodeServerV2] master received subscribe request from node[%s]", request.NodeKey)

	// finished channel
	finished := make(chan bool)

	// set subscribe
	svr.server.SetSubscribe("node:"+request.NodeKey, &entity.GrpcSubscribe{
		Stream:   stream,
		Finished: finished,
	})
	ctx := stream.Context()

	log.Infof("[NodeServerV2] master subscribed node[%s]", request.NodeKey)

	// Keep this scope alive because once this scope exits - the stream is closed
	for {
		select {
		case <-finished:
			log.Infof("[NodeServerV2] closing stream for node[%s]", request.NodeKey)
			return nil
		case <-ctx.Done():
			log.Infof("[NodeServerV2] node[%s] has disconnected", request.NodeKey)
			return nil
		}
	}
}

func (svr NodeServerV2) Unsubscribe(_ context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	sub, err := svr.server.GetSubscribe("node:" + req.NodeKey)
	if err != nil {
		return nil, errors.ErrorGrpcSubscribeNotExists
	}
	select {
	case sub.GetFinished() <- true:
		log.Infof("unsubscribed node[%s]", req.NodeKey)
	default:
		// Default case is to avoid blocking in case client has already unsubscribed
	}
	svr.server.DeleteSubscribe(req.NodeKey)
	return &grpc.Response{
		Code:    grpc.ResponseCode_OK,
		Message: "unsubscribed successfully",
	}, nil
}

var nodeSvrV2 *NodeServerV2
var nodeSvrV2Once = new(sync.Once)

func NewNodeServerV2() (res *NodeServerV2, err error) {
	if nodeSvrV2 != nil {
		return nodeSvrV2, nil
	}
	nodeSvrV2Once.Do(func() {
		nodeSvrV2 = &NodeServerV2{}
		nodeSvrV2.cfgSvc = nodeconfig.GetNodeConfigService()
		if err != nil {
			log.Errorf("[NodeServerV2] error: %s", err.Error())
		}
	})
	if err != nil {
		return nil, err
	}
	return nodeSvrV2, nil
}
