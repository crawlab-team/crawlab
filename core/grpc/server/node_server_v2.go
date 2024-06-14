package server

import (
	"context"
	"encoding/json"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/entity"
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/models/service"
	nodeconfig "github.com/crawlab-team/crawlab/core/node/config"
	"github.com/crawlab-team/crawlab/grpc"
	errors2 "github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
func (svr NodeServerV2) Register(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	// unmarshall data
	var node models.NodeV2
	if req.Data != nil {
		if err := json.Unmarshal(req.Data, &node); err != nil {
			return HandleError(err)
		}

		if node.IsMaster {
			// error: cannot register master node
			return HandleError(errors.ErrorGrpcNotAllowed)
		}
	}

	// node key
	var nodeKey string
	if req.NodeKey != "" {
		nodeKey = req.NodeKey
	} else {
		nodeKey = node.Key
	}
	if nodeKey == "" {
		return HandleError(errors.ErrorModelMissingRequiredData)
	}

	// find in db
	nodeDb, err := service.NewModelServiceV2[models.NodeV2]().GetOne(bson.M{"key": nodeKey}, nil)
	if err == nil {
		if node.IsMaster {
			// error: cannot register master node
			return HandleError(errors.ErrorGrpcNotAllowed)
		} else {
			// register existing
			nodeDb.Status = constants.NodeStatusRegistered
			nodeDb.Active = true
			err = service.NewModelServiceV2[models.NodeV2]().ReplaceById(nodeDb.Id, *nodeDb)
			if err != nil {
				return HandleError(err)
			}
			log.Infof("[NodeServerV2] updated worker[%s] in db. id: %s", nodeKey, node.Id.Hex())
		}
	} else if errors2.Is(err, mongo.ErrNoDocuments) {
		// register new
		node.Key = nodeKey
		node.Status = constants.NodeStatusRegistered
		node.Active = true
		node.ActiveAt = time.Now()
		node.Enabled = true
		if node.Name == "" {
			node.Name = nodeKey
		}
		node.SetCreated(primitive.NilObjectID)
		node.SetUpdated(primitive.NilObjectID)
		_, err = service.NewModelServiceV2[models.NodeV2]().InsertOne(*nodeDb)
		if err != nil {
			return HandleError(err)
		}
		log.Infof("[NodeServerV2] added worker[%s] in db. id: %s", nodeKey, node.Id.Hex())
	} else {
		// error
		return HandleError(err)
	}

	log.Infof("[NodeServerV2] master registered worker[%s]", req.GetNodeKey())

	return HandleSuccessWithData(node)
}

// SendHeartbeat from worker to master
func (svr NodeServerV2) SendHeartbeat(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	// find in db
	node, err := service.NewModelServiceV2[models.NodeV2]().GetOne(bson.M{"key": req.NodeKey}, nil)
	if err != nil {
		if errors2.Is(err, mongo.ErrNoDocuments) {
			return HandleError(errors.ErrorNodeNotExists)
		}
		return HandleError(err)
	}

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

	return HandleSuccessWithData(node)
}

// Ping from worker to master
func (svr NodeServerV2) Ping(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
	return HandleSuccess()
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

func (svr NodeServerV2) Unsubscribe(ctx context.Context, req *grpc.Request) (res *grpc.Response, err error) {
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

func NewNodeServerV2() (res *NodeServerV2, err error) {
	// node server
	svr := &NodeServerV2{}
	svr.cfgSvc, err = nodeconfig.NewNodeConfigService()
	if err != nil {
		return nil, err
	}

	return svr, nil
}
