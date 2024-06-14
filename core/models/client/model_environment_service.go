package client

import (
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EnvironmentServiceDelegate struct {
	interfaces.GrpcClientModelBaseService
}

func (svc *EnvironmentServiceDelegate) GetEnvironmentById(id primitive.ObjectID) (e interfaces.Environment, err error) {
	res, err := svc.GetById(id)
	if err != nil {
		return nil, err
	}
	s, ok := res.(interfaces.Environment)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return s, nil
}

func (svc *EnvironmentServiceDelegate) GetEnvironment(query bson.M, opts *mongo.FindOptions) (e interfaces.Environment, err error) {
	res, err := svc.Get(query, opts)
	if err != nil {
		return nil, err
	}
	s, ok := res.(interfaces.Environment)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return s, nil
}

func (svc *EnvironmentServiceDelegate) GetEnvironmentList(query bson.M, opts *mongo.FindOptions) (res []interfaces.Environment, err error) {
	list, err := svc.GetList(query, opts)
	if err != nil {
		return nil, err
	}
	for _, item := range list.GetModels() {
		s, ok := item.(interfaces.Environment)
		if !ok {
			return nil, errors.ErrorModelInvalidType
		}
		res = append(res, s)
	}
	return res, nil
}

func NewEnvironmentServiceDelegate() (svc2 interfaces.GrpcClientModelEnvironmentService, err error) {
	var opts []ModelBaseServiceDelegateOption

	// apply options
	opts = append(opts, WithBaseServiceModelId(interfaces.ModelIdEnvironment))

	// base service
	baseSvc, err := NewBaseServiceDelegate(opts...)
	if err != nil {
		return nil, err
	}

	// service
	svc := &EnvironmentServiceDelegate{baseSvc}

	return svc, nil
}
