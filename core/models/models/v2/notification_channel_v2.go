package models

type NotificationChannelV2 struct {
	any                                `collection:"notification_channels"`
	BaseModelV2[NotificationChannelV2] `bson:",inline"`
	Type                               string `json:"type" bson:"type"`
	Name                               string `json:"name" bson:"name"`
	Description                        string `json:"description" bson:"description"`
	Provider                           string `json:"provider" bson:"provider"`
	SMTPServer                         string `json:"smtp_server,omitempty" bson:"smtp_server,omitempty"`
	SMTPPort                           int    `json:"smtp_port,omitempty" bson:"smtp_port,omitempty"`
	SMTPUsername                       string `json:"smtp_username,omitempty" bson:"smtp_username,omitempty"`
	SMTPPassword                       string `json:"smtp_password,omitempty" bson:"smtp_password,omitempty"`
	WebhookUrl                         string `json:"webhook_url,omitempty" bson:"webhook_url,omitempty"`
	TelegramBotToken                   string `json:"telegram_bot_token,omitempty" bson:"telegram_bot_token,omitempty"`
	TelegramChatId                     string `json:"telegram_chat_id,omitempty" bson:"telegram_chat_id,omitempty"`
	GoogleOAuth2Json                   string `json:"google_oauth2_json,omitempty" bson:"google_oauth2_json,omitempty"`
}
