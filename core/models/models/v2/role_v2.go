package models

import "github.com/crawlab-team/crawlab/core/models/models/v2"

type RoleV2 struct {
	any                        `collection:"roles"`
	models.BaseModelV2[RoleV2] `bson:",inline"`
	Key                        string `json:"key" bson:"key"`
	Name                       string `json:"name" bson:"name"`
	Description                string `json:"description" bson:"description"`
}
