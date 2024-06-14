package service

import (
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	models2 "github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func convertTypePassword(d interface{}, err error) (res *models2.Password, err2 error) {
	if err != nil {
		return nil, err
	}
	res, ok := d.(*models2.Password)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return res, nil
}

func (svc *Service) GetPasswordById(id primitive.ObjectID) (res *models2.Password, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdPassword).GetById(id)
	return convertTypePassword(d, err)
}

func (svc *Service) GetPassword(query bson.M, opts *mongo.FindOptions) (res *models2.Password, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdPassword).Get(query, opts)
	return convertTypePassword(d, err)
}

func (svc *Service) GetPasswordList(query bson.M, opts *mongo.FindOptions) (res []models2.Password, err error) {
	l, err := svc.GetBaseService(interfaces.ModelIdPassword).GetList(query, opts)
	for _, doc := range l.GetModels() {
		d := doc.(*models2.Password)
		res = append(res, *d)
	}
	return res, nil
}
