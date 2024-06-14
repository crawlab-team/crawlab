package service

import (
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	models2 "github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func convertTypePermission(d interface{}, err error) (res *models2.Permission, err2 error) {
	if err != nil {
		return nil, err
	}
	res, ok := d.(*models2.Permission)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return res, nil
}

func (svc *Service) GetPermissionById(id primitive.ObjectID) (res *models2.Permission, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdPermission).GetById(id)
	return convertTypePermission(d, err)
}

func (svc *Service) GetPermission(query bson.M, opts *mongo.FindOptions) (res *models2.Permission, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdPermission).Get(query, opts)
	return convertTypePermission(d, err)
}

func (svc *Service) GetPermissionList(query bson.M, opts *mongo.FindOptions) (res []models2.Permission, err error) {
	l, err := svc.GetBaseService(interfaces.ModelIdPermission).GetList(query, opts)
	if err != nil {
		return nil, err
	}
	for _, doc := range l.GetModels() {
		d := doc.(*models2.Permission)
		res = append(res, *d)
	}
	return res, nil
}

func (svc *Service) GetPermissionByKey(key string, opts *mongo.FindOptions) (res *models2.Permission, err error) {
	query := bson.M{"key": key}
	return svc.GetPermission(query, opts)
}
