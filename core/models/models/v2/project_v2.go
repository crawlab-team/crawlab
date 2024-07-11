package models

import "github.com/crawlab-team/crawlab/core/models/models/v2"

type ProjectV2 struct {
	any                           `collection:"projects"`
	models.BaseModelV2[ProjectV2] `bson:",inline"`
	Name                          string `json:"name" bson:"name"`
	Description                   string `json:"description" bson:"description"`
	Spiders                       int    `json:"spiders" bson:"-"`
}
