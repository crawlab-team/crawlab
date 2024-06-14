package service

import (
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	models2 "github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func convertTypeRole(d interface{}, err error) (res *models2.Role, err2 error) {
	if err != nil {
		return nil, err
	}
	res, ok := d.(*models2.Role)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return res, nil
}

func (svc *Service) GetRoleById(id primitive.ObjectID) (res *models2.Role, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdRole).GetById(id)
	return convertTypeRole(d, err)
}

func (svc *Service) GetRole(query bson.M, opts *mongo.FindOptions) (res *models2.Role, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdRole).Get(query, opts)
	return convertTypeRole(d, err)
}

func (svc *Service) GetRoleList(query bson.M, opts *mongo.FindOptions) (res []models2.Role, err error) {
	l, err := svc.GetBaseService(interfaces.ModelIdRole).GetList(query, opts)
	if err != nil {
		return nil, err
	}
	for _, doc := range l.GetModels() {
		d := doc.(*models2.Role)
		res = append(res, *d)
	}
	return res, nil
}

func (svc *Service) GetRoleByName(name string, opts *mongo.FindOptions) (res *models2.Role, err error) {
	query := bson.M{"name": name}
	return svc.GetRole(query, opts)
}

func (svc *Service) GetRoleByKey(key string, opts *mongo.FindOptions) (res *models2.Role, err error) {
	query := bson.M{"key": key}
	return svc.GetRole(query, opts)
}
