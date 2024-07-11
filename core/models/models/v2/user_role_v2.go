package models

import (
	"github.com/crawlab-team/crawlab/core/models/models/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRoleV2 struct {
	any                            `collection:"user_roles"`
	models.BaseModelV2[UserRoleV2] `bson:",inline"`
	RoleId                         primitive.ObjectID `json:"role_id" bson:"role_id"`
	UserId                         primitive.ObjectID `json:"user_id" bson:"user_id"`
}
