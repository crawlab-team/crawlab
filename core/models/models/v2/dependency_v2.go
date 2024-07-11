package models

import (
	"github.com/crawlab-team/crawlab/core/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DependencyV2 struct {
	any                       `collection:"dependencies"`
	BaseModelV2[DependencyV2] `bson:",inline"`
	Name                      string                  `json:"name" bson:"name"`
	Description               string                  `json:"description" bson:"description"`
	NodeId                    primitive.ObjectID      `json:"node_id" bson:"node_id"`
	Type                      string                  `json:"type" bson:"type"`
	LatestVersion             string                  `json:"latest_version" bson:"latest_version"`
	Version                   string                  `json:"version" bson:"version"`
	Result                    entity.DependencyResult `json:"result" bson:"-"`
}
