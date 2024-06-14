package client

import (
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NodeServiceDelegate struct {
	interfaces.GrpcClientModelBaseService
}

func (svc *NodeServiceDelegate) GetNodeById(id primitive.ObjectID) (n interfaces.Node, err error) {
	res, err := svc.GetById(id)
	if err != nil {
		return nil, err
	}
	s, ok := res.(interfaces.Node)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return s, nil
}

func (svc *NodeServiceDelegate) GetNode(query bson.M, opts *mongo.FindOptions) (n interfaces.Node, err error) {
	res, err := svc.Get(query, opts)
	if err != nil {
		return nil, err
	}
	s, ok := res.(interfaces.Node)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return s, nil
}

func (svc *NodeServiceDelegate) GetNodeByKey(key string) (n interfaces.Node, err error) {
	return svc.GetNode(bson.M{"key": key}, nil)
}

func (svc *NodeServiceDelegate) GetNodeList(query bson.M, opts *mongo.FindOptions) (res []interfaces.Node, err error) {
	list, err := svc.GetList(query, opts)
	if err != nil {
		return nil, err
	}
	for _, item := range list.GetModels() {
		s, ok := item.(interfaces.Node)
		if !ok {
			return nil, errors.ErrorModelInvalidType
		}
		res = append(res, s)
	}
	return res, nil
}

func NewNodeServiceDelegate() (svc2 interfaces.GrpcClientModelNodeService, err error) {
	var opts []ModelBaseServiceDelegateOption

	// apply options
	opts = append(opts, WithBaseServiceModelId(interfaces.ModelIdNode))

	// base service
	baseSvc, err := NewBaseServiceDelegate(opts...)
	if err != nil {
		return nil, err
	}

	// service
	svc := &NodeServiceDelegate{baseSvc}

	return svc, nil
}
