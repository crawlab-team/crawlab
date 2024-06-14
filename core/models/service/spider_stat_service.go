package service

import (
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	models2 "github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func convertTypeSpiderStat(d interface{}, err error) (res *models2.SpiderStat, err2 error) {
	if err != nil {
		return nil, err
	}
	res, ok := d.(*models2.SpiderStat)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return res, nil
}

func (svc *Service) GetSpiderStatById(id primitive.ObjectID) (res *models2.SpiderStat, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdSpiderStat).GetById(id)
	return convertTypeSpiderStat(d, err)
}

func (svc *Service) GetSpiderStat(query bson.M, opts *mongo.FindOptions) (res *models2.SpiderStat, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdSpiderStat).Get(query, opts)
	return convertTypeSpiderStat(d, err)
}

func (svc *Service) GetSpiderStatList(query bson.M, opts *mongo.FindOptions) (res []models2.SpiderStat, err error) {
	l, err := svc.GetBaseService(interfaces.ModelIdSpiderStat).GetList(query, opts)
	for _, doc := range l.GetModels() {
		d := doc.(*models2.SpiderStat)
		res = append(res, *d)
	}
	return res, nil
}
