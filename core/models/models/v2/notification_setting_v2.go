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

	TaskTrigger   string `json:"task_trigger" bson:"task_trigger"`
	TriggerTarget string `json:"trigger_target" bson:"trigger_target"`
	Trigger       string `json:"trigger" bson:"trigger"`

	HasMail     bool   `json:"has_mail" bson:"has_mail"`
	SenderEmail string `json:"sender_email" bson:"sender_email"`
	SenderName  string `json:"sender_name" bson:"sender_name"`
	MailTo      string `json:"mail_to" bson:"mail_to"`
	MailCc      string `json:"mail_cc" bson:"mail_cc"`
	MailBcc     string `json:"mail_bcc" bson:"mail_bcc"`

	ChannelIds []primitive.ObjectID    `json:"channel_ids,omitempty" bson:"channel_ids,omitempty"`
	Channels   []NotificationChannelV2 `json:"channels,omitempty" bson:"-"`
}
