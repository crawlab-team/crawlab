package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type DependencyTaskV2 struct {
	any                           `collection:"dependency_tasks"`
	BaseModelV2[DependencyTaskV2] `bson:",inline"`
	Status                        string             `json:"status" bson:"status"`
	Error                         string             `json:"error" bson:"error"`
	SettingId                     primitive.ObjectID `json:"setting_id" bson:"setting_id"`
	Type                          string             `json:"type" bson:"type"`
	NodeId                        primitive.ObjectID `json:"node_id" bson:"node_id"`
	Action                        string             `json:"action" bson:"action"`
	DepNames                      []string           `json:"dep_names" bson:"dep_names"`
}
