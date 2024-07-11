package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type DependencyLogV2 struct {
	any                          `collection:"dependency_logs"`
	BaseModelV2[DependencyLogV2] `bson:",inline"`
	TaskId                       primitive.ObjectID `json:"task_id" bson:"task_id"`
	Content                      string             `json:"content" bson:"content"`
}
