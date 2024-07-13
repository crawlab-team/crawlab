package models

type NotificationSettingV2 struct {
	any                                `collection:"notification_settings"`
	BaseModelV2[NotificationSettingV2] `bson:",inline"`
	Type                               string                    `json:"type" bson:"type"`
	Name                               string                    `json:"name" bson:"name"`
	Description                        string                    `json:"description" bson:"description"`
	Enabled                            bool                      `json:"enabled" bson:"enabled"`
	Global                             bool                      `json:"global" bson:"global"`
	Title                              string                    `json:"title,omitempty" bson:"title,omitempty"`
	Template                           string                    `json:"template" bson:"template"`
	TemplateMode                       string                    `json:"template_mode" bson:"template_mode"`
	TemplateMarkdown                   string                    `json:"template_markdown,omitempty" bson:"template_markdown,omitempty"`
	TemplateRichText                   string                    `json:"template_rich_text,omitempty" bson:"template_rich_text,omitempty"`
	TemplateRichTextJson               string                    `json:"template_rich_text_json,omitempty" bson:"template_rich_text_json,omitempty"`
	TaskTrigger                        string                    `json:"task_trigger" bson:"task_trigger"`
	TriggerTarget                      string                    `json:"trigger_target" bson:"trigger_target"`
	Trigger                            string                    `json:"trigger" bson:"trigger"`
	Mail                               NotificationSettingMail   `json:"mail,omitempty" bson:"mail,omitempty"`
	Mobile                             NotificationSettingMobile `json:"mobile,omitempty" bson:"mobile,omitempty"`
}

type NotificationSettingMail struct {
	Server         string `json:"server" bson:"server"`
	Port           string `json:"port,omitempty" bson:"port,omitempty"`
	User           string `json:"user,omitempty" bson:"user,omitempty"`
	Password       string `json:"password,omitempty" bson:"password,omitempty"`
	SenderEmail    string `json:"sender_email,omitempty" bson:"sender_email,omitempty"`
	SenderIdentity string `json:"sender_identity,omitempty" bson:"sender_identity,omitempty"`
	To             string `json:"to,omitempty" bson:"to,omitempty"`
	Cc             string `json:"cc,omitempty" bson:"cc,omitempty"`
}

type NotificationSettingMobile struct {
	Webhook string `json:"webhook" bson:"webhook"`
}
