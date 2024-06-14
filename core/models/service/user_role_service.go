package service

import (
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	models2 "github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func convertTypeUserRole(d interface{}, err error) (res *models2.UserRole, err2 error) {
	if err != nil {
		return nil, err
	}
	res, ok := d.(*models2.UserRole)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return res, nil
}

func (svc *Service) GetUserRoleById(id primitive.ObjectID) (res *models2.UserRole, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdUserRole).GetById(id)
	return convertTypeUserRole(d, err)
}

func (svc *Service) GetUserRole(query bson.M, opts *mongo.FindOptions) (res *models2.UserRole, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdUserRole).Get(query, opts)
	return convertTypeUserRole(d, err)
}

func (svc *Service) GetUserRoleList(query bson.M, opts *mongo.FindOptions) (res []models2.UserRole, err error) {
	l, err := svc.GetBaseService(interfaces.ModelIdUserRole).GetList(query, opts)
	if err != nil {
		return nil, err
	}
	for _, doc := range l.GetModels() {
		d := doc.(*models2.UserRole)
		res = append(res, *d)
	}
	return res, nil
}

func (svc *Service) GetUserRoleListByUserId(id primitive.ObjectID, opts *mongo.FindOptions) (res []models2.UserRole, err error) {
	return svc.GetUserRoleList(bson.M{"user_id": id}, opts)
}

func (svc *Service) GetUserRoleListByRoleId(id primitive.ObjectID, opts *mongo.FindOptions) (res []models2.UserRole, err error) {
	return svc.GetUserRoleList(bson.M{"role_id": id}, opts)
}
