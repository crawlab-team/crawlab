package service

import (
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	models2 "github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func convertTypeTaskStat(d interface{}, err error) (res *models2.TaskStat, err2 error) {
	if err != nil {
		return nil, err
	}
	res, ok := d.(*models2.TaskStat)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return res, nil
}

func (svc *Service) GetTaskStatById(id primitive.ObjectID) (res *models2.TaskStat, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdTaskStat).GetById(id)
	return convertTypeTaskStat(d, err)
}

func (svc *Service) GetTaskStat(query bson.M, opts *mongo.FindOptions) (res *models2.TaskStat, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdTaskStat).Get(query, opts)
	return convertTypeTaskStat(d, err)
}

func (svc *Service) GetTaskStatList(query bson.M, opts *mongo.FindOptions) (res []models2.TaskStat, err error) {
	l, err := svc.GetBaseService(interfaces.ModelIdTaskStat).GetList(query, opts)
	for _, doc := range l.GetModels() {
		d := doc.(*models2.TaskStat)
		res = append(res, *d)
	}
	return res, nil
}
