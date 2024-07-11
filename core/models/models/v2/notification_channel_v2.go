package models

type NotificationChannelV2 struct {
	any                                `collection:"notification_channels"`
	BaseModelV2[NotificationChannelV2] `bson:",inline"`
	Name                               string `json:"name" bson:"name"`
	Type                               string `json:"type" bson:"type"`
}
