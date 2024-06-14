package service

import (
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	models2 "github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func convertTypeDependencySetting(d interface{}, err error) (res *models2.DependencySetting, err2 error) {
	if err != nil {
		return nil, err
	}
	res, ok := d.(*models2.DependencySetting)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return res, nil
}

func (svc *Service) GetDependencySettingById(id primitive.ObjectID) (res *models2.DependencySetting, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdDependencySetting).GetById(id)
	return convertTypeDependencySetting(d, err)
}

func (svc *Service) GetDependencySetting(query bson.M, opts *mongo.FindOptions) (res *models2.DependencySetting, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdDependencySetting).Get(query, opts)
	return convertTypeDependencySetting(d, err)
}

func (svc *Service) GetDependencySettingList(query bson.M, opts *mongo.FindOptions) (res []models2.DependencySetting, err error) {
	l, err := svc.GetBaseService(interfaces.ModelIdDependencySetting).GetList(query, opts)
	for _, doc := range l.GetModels() {
		d := doc.(*models2.DependencySetting)
		res = append(res, *d)
	}
	return res, nil
}
