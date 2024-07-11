package models

import (
	"github.com/crawlab-team/crawlab/core/models/models/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type SettingV2 struct {
	any                           `collection:"settings"`
	models.BaseModelV2[SettingV2] `bson:",inline"`
	Key                           string `json:"key" bson:"key"`
	Value                         bson.M `json:"value" bson:"value"`
}
