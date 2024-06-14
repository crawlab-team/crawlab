package service

import (
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	models2 "github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func convertTypeSetting(d interface{}, err error) (res *models2.Setting, err2 error) {
	if err != nil {
		return nil, err
	}
	res, ok := d.(*models2.Setting)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return res, nil
}

func (svc *Service) GetSettingById(id primitive.ObjectID) (res *models2.Setting, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdSetting).GetById(id)
	return convertTypeSetting(d, err)
}

func (svc *Service) GetSetting(query bson.M, opts *mongo.FindOptions) (res *models2.Setting, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdSetting).Get(query, opts)
	return convertTypeSetting(d, err)
}

func (svc *Service) GetSettingList(query bson.M, opts *mongo.FindOptions) (res []models2.Setting, err error) {
	l, err := svc.GetBaseService(interfaces.ModelIdSetting).GetList(query, opts)
	for _, doc := range l.GetModels() {
		d := doc.(*models2.Setting)
		res = append(res, *d)
	}
	return res, nil
}

func (svc *Service) GetSettingByKey(key string, opts *mongo.FindOptions) (res *models2.Setting, err error) {
	query := bson.M{"key": key}
	return svc.GetSetting(query, opts)
}
