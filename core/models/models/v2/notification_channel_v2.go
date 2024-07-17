package models

type NotificationChannelV2 struct {
	any                                `collection:"notification_channels"`
	BaseModelV2[NotificationChannelV2] `bson:",inline"`
	Type                               string `json:"type" bson:"type"`
	Name                               string `json:"name" bson:"name"`
	Description                        string `json:"description" bson:"description"`
	Provider                           string `json:"provider" bson:"provider"`
	MailSettings                       struct {
		SMTPServer   string `json:"smtp_server" bson:"smtp_server"`
		SMTPPort     string `json:"smtp_port" bson:"smtp_port"`
		SMTPUser     string `json:"smtp_from_email_address" bson:"smtp_from_email_address"`
		SMTPPassword string `json:"smtp_from_email_password" bson:"smtp_from_email_password"`
	} `json:"mail_settings,omitempty" bson:"mail_settings,omitempty"`
	IMSettings struct {
		Webhook string `json:"webhook" bson:"webhook"`
	} `json:"im_settings,omitempty" bson:"im_settings,omitempty"`
}
