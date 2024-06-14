package models

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RolePermission struct {
	Id           primitive.ObjectID `json:"_id" bson:"_id"`
	RoleId       primitive.ObjectID `json:"role_id" bson:"role_id"`
	PermissionId primitive.ObjectID `json:"permission_id" bson:"permission_id"`
}

func (ur *RolePermission) GetId() (id primitive.ObjectID) {
	return ur.Id
}

func (ur *RolePermission) SetId(id primitive.ObjectID) {
	ur.Id = id
}

type RolePermissionList []RolePermission

func (l *RolePermissionList) GetModels() (res []interfaces.Model) {
	for i := range *l {
		d := (*l)[i]
		res = append(res, &d)
	}
	return res
}
