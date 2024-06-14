package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRoleV2 struct {
	any                     `collection:"user_roles"`
	BaseModelV2[UserRoleV2] `bson:",inline"`
	RoleId                  primitive.ObjectID `json:"role_id" bson:"role_id"`
	UserId                  primitive.ObjectID `json:"user_id" bson:"user_id"`
}
