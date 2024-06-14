package models

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRole struct {
	Id     primitive.ObjectID `json:"_id" bson:"_id"`
	RoleId primitive.ObjectID `json:"role_id" bson:"role_id"`
	UserId primitive.ObjectID `json:"user_id" bson:"user_id"`
}

func (ur *UserRole) GetId() (id primitive.ObjectID) {
	return ur.Id
}

func (ur *UserRole) SetId(id primitive.ObjectID) {
	ur.Id = id
}

type UserRoleList []UserRole

func (l *UserRoleList) GetModels() (res []interfaces.Model) {
	for i := range *l {
		d := (*l)[i]
		res = append(res, &d)
	}
	return res
}
