package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type NotificationAlertV2 struct {
	any                              `collection:"notification_alerts"`
	BaseModelV2[NotificationAlertV2] `bson:",inline"`
	Name                             string             `json:"name" bson:"name"`
	Description                      string             `json:"description" bson:"description"`
	Enabled                          bool               `json:"enabled" bson:"enabled"`
	MetricTargetId                   primitive.ObjectID `json:"metric_target_id,omitempty" bson:"metric_target_id,omitempty"`
	MetricName                       string             `json:"metric_name" bson:"metric_name"`
	Operator                         string             `json:"operator" bson:"operator"`
	TargetValue                      float32            `json:"target_value" bson:"target_value"`
	Level                            string             `json:"level" bson:"level"`
}
