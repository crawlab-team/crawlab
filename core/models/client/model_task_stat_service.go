package client

import (
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskStatServiceDelegate struct {
	interfaces.GrpcClientModelBaseService
}

func (svc *TaskStatServiceDelegate) GetTaskStatById(id primitive.ObjectID) (t interfaces.TaskStat, err error) {
	res, err := svc.GetById(id)
	if err != nil {
		return nil, err
	}
	s, ok := res.(interfaces.TaskStat)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return s, nil
}

func (svc *TaskStatServiceDelegate) GetTaskStat(query bson.M, opts *mongo.FindOptions) (t interfaces.TaskStat, err error) {
	res, err := svc.Get(query, opts)
	if err != nil {
		return nil, err
	}
	s, ok := res.(interfaces.TaskStat)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return s, nil
}

func (svc *TaskStatServiceDelegate) GetTaskStatList(query bson.M, opts *mongo.FindOptions) (res []interfaces.TaskStat, err error) {
	list, err := svc.GetList(query, opts)
	if err != nil {
		return nil, err
	}
	for _, item := range list.GetModels() {
		s, ok := item.(interfaces.TaskStat)
		if !ok {
			return nil, errors.ErrorModelInvalidType
		}
		res = append(res, s)
	}
	return res, nil
}

func NewTaskStatServiceDelegate() (svc2 interfaces.GrpcClientModelTaskStatService, err error) {
	var opts []ModelBaseServiceDelegateOption

	// apply options
	opts = append(opts, WithBaseServiceModelId(interfaces.ModelIdTaskStat))

	// base service
	baseSvc, err := NewBaseServiceDelegate(opts...)
	if err != nil {
		return nil, err
	}

	// service
	svc := &TaskStatServiceDelegate{baseSvc}

	return svc, nil
}
