package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type NotificationSettingV2 struct {
	Id          primitive.ObjectID        `json:"_id" bson:"_id"`
	Type        string                    `json:"type" bson:"type"`
	Name        string                    `json:"name" bson:"name"`
	Description string                    `json:"description" bson:"description"`
	Enabled     bool                      `json:"enabled" bson:"enabled"`
	Global      bool                      `json:"global" bson:"global"`
	Title       string                    `json:"title,omitempty" bson:"title,omitempty"`
	Template    string                    `json:"template,omitempty" bson:"template,omitempty"`
	TaskTrigger string                    `json:"task_trigger" bson:"task_trigger"`
	Mail        NotificationSettingMail   `json:"mail,omitempty" bson:"mail,omitempty"`
	Mobile      NotificationSettingMobile `json:"mobile,omitempty" bson:"mobile,omitempty"`
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
