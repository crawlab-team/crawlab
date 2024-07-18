package models

type NotificationChannelV2 struct {
	any                                `collection:"notification_channels"`
	BaseModelV2[NotificationChannelV2] `bson:",inline"`
	Type                               string `json:"type" bson:"type"`
	Name                               string `json:"name" bson:"name"`
	Description                        string `json:"description" bson:"description"`
	Provider                           string `json:"provider" bson:"provider"`
	SMTPServer                         string `json:"smtp_server" bson:"smtp_server"`
	SMTPPort                           string `json:"smtp_port" bson:"smtp_port"`
	SMTPUsername                       string `json:"smtp_username" bson:"smtp_username"`
	SMTPPassword                       string `json:"smtp_password" bson:"smtp_password"`
	WebhookUrl                         string `json:"webhook_url" bson:"webhook_url"`
}
