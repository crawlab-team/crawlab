package service

import (
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	models2 "github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func convertTypeRolePermission(d interface{}, err error) (res *models2.RolePermission, err2 error) {
	if err != nil {
		return nil, err
	}
	res, ok := d.(*models2.RolePermission)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return res, nil
}

func (svc *Service) GetRolePermissionById(id primitive.ObjectID) (res *models2.RolePermission, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdRolePermission).GetById(id)
	return convertTypeRolePermission(d, err)
}

func (svc *Service) GetRolePermission(query bson.M, opts *mongo.FindOptions) (res *models2.RolePermission, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdRolePermission).Get(query, opts)
	return convertTypeRolePermission(d, err)
}

func (svc *Service) GetRolePermissionList(query bson.M, opts *mongo.FindOptions) (res []models2.RolePermission, err error) {
	l, err := svc.GetBaseService(interfaces.ModelIdRolePermission).GetList(query, opts)
	if err != nil {
		return nil, err
	}
	for _, doc := range l.GetModels() {
		d := doc.(*models2.RolePermission)
		res = append(res, *d)
	}
	return res, nil
}

func (svc *Service) GetRolePermissionListByRoleId(id primitive.ObjectID, opts *mongo.FindOptions) (res []models2.RolePermission, err error) {
	return svc.GetRolePermissionList(bson.M{"role_id": id}, opts)
}

func (svc *Service) GetRolePermissionListByPermissionId(id primitive.ObjectID, opts *mongo.FindOptions) (res []models2.RolePermission, err error) {
	return svc.GetRolePermissionList(bson.M{"permission_id": id}, opts)
}
