package models

import "github.com/crawlab-team/crawlab/core/models/models/v2"

type TokenV2 struct {
	any                         `collection:"tokens"`
	models.BaseModelV2[TokenV2] `bson:",inline"`
	Name                        string `json:"name" bson:"name"`
	Token                       string `json:"token" bson:"token"`
}
