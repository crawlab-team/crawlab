package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type NotificationSettingV2 struct {
	any                                `collection:"notification_settings"`
	BaseModelV2[NotificationSettingV2] `bson:",inline"`
	Name                               string `json:"name" bson:"name"`
	Description                        string `json:"description" bson:"description"`
	Enabled                            bool   `json:"enabled" bson:"enabled"`

	Title                string `json:"title,omitempty" bson:"title,omitempty"`
	Template             string `json:"template" bson:"template"`
	TemplateMode         string `json:"template_mode" bson:"template_mode"`
	TemplateMarkdown     string `json:"template_markdown,omitempty" bson:"template_markdown,omitempty"`
	TemplateRichText     string `json:"template_rich_text,omitempty" bson:"template_rich_text,omitempty"`
	TemplateRichTextJson string `json:"template_rich_text_json,omitempty" bson:"template_rich_text_json,omitempty"`
	TemplateTheme        string `json:"template_theme,omitempty" bson:"template_theme,omitempty"`

	TaskTrigger string `json:"task_trigger" bson:"task_trigger"`
	Trigger     string `json:"trigger" bson:"trigger"`

	SenderEmail          string   `json:"sender_email,omitempty" bson:"sender_email,omitempty"`
	UseCustomSenderEmail bool     `json:"use_custom_sender_email,omitempty" bson:"use_custom_sender_email,omitempty"`
	SenderName           string   `json:"sender_name,omitempty" bson:"sender_name,omitempty"`
	MailTo               []string `json:"mail_to,omitempty" bson:"mail_to,omitempty"`
	MailCc               []string `json:"mail_cc,omitempty" bson:"mail_cc,omitempty"`
	MailBcc              []string `json:"mail_bcc,omitempty" bson:"mail_bcc,omitempty"`

	ChannelIds []primitive.ObjectID    `json:"channel_ids,omitempty" bson:"channel_ids,omitempty"`
	Channels   []NotificationChannelV2 `json:"channels,omitempty" bson:"-"`

	AlertId primitive.ObjectID `json:"alert_id,omitempty" bson:"alert_id,omitempty"`
}
