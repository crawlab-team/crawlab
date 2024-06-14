package service

import (
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	models2 "github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func convertTypeTaskQueueItem(d interface{}, err error) (res *models2.TaskQueueItem, err2 error) {
	if err != nil {
		return nil, err
	}
	res, ok := d.(*models2.TaskQueueItem)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return res, nil
}

func (svc *Service) GetTaskQueueItemById(id primitive.ObjectID) (res *models2.TaskQueueItem, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdTaskQueue).GetById(id)
	return convertTypeTaskQueueItem(d, err)
}

func (svc *Service) GetTaskQueueItem(query bson.M, opts *mongo.FindOptions) (res *models2.TaskQueueItem, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdTaskQueue).Get(query, opts)
	return convertTypeTaskQueueItem(d, err)
}

func (svc *Service) GetTaskQueueItemList(query bson.M, opts *mongo.FindOptions) (res []models2.TaskQueueItem, err error) {
	l, err := svc.GetBaseService(interfaces.ModelIdTaskQueue).GetList(query, opts)
	for _, doc := range l.GetModels() {
		d := doc.(*models2.TaskQueueItem)
		res = append(res, *d)
	}
	return res, nil
}
