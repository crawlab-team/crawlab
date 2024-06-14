package service

import (
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	models2 "github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func convertTypeEnvironment(d interface{}, err error) (res *models2.Environment, err2 error) {
	if err != nil {
		return nil, err
	}
	res, ok := d.(*models2.Environment)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return res, nil
}

func (svc *Service) GetEnvironmentById(id primitive.ObjectID) (res *models2.Environment, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdEnvironment).GetById(id)
	return convertTypeEnvironment(d, err)
}

func (svc *Service) GetEnvironment(query bson.M, opts *mongo.FindOptions) (res *models2.Environment, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdEnvironment).Get(query, opts)
	return convertTypeEnvironment(d, err)
}

func (svc *Service) GetEnvironmentList(query bson.M, opts *mongo.FindOptions) (res []models2.Environment, err error) {
	l, err := svc.GetBaseService(interfaces.ModelIdEnvironment).GetList(query, opts)
	for _, doc := range l.GetModels() {
		d := doc.(*models2.Environment)
		res = append(res, *d)
	}
	return res, nil
}
