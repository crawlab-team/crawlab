package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type NotificationRequestV2 struct {
	any                                `collection:"notification_requests"`
	BaseModelV2[NotificationRequestV2] `bson:",inline"`
	Status                             string             `json:"status" bson:"status"`
	Error                              string             `json:"error,omitempty" bson:"error,omitempty"`
	SettingId                          primitive.ObjectID `json:"setting_id" bson:"setting_id"`
	ChannelId                          primitive.ObjectID `json:"channel_id" bson:"channel_id"`
}
