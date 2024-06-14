package models

import (
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ScheduleV2 struct {
	any                     `collection:"schedules"`
	BaseModelV2[ScheduleV2] `bson:",inline"`
	Name                    string               `json:"name" bson:"name"`
	Description             string               `json:"description" bson:"description"`
	SpiderId                primitive.ObjectID   `json:"spider_id" bson:"spider_id"`
	Cron                    string               `json:"cron" bson:"cron"`
	EntryId                 cron.EntryID         `json:"entry_id" bson:"entry_id"`
	Cmd                     string               `json:"cmd" bson:"cmd"`
	Param                   string               `json:"param" bson:"param"`
	Mode                    string               `json:"mode" bson:"mode"`
	NodeIds                 []primitive.ObjectID `json:"node_ids" bson:"node_ids"`
	Priority                int                  `json:"priority" bson:"priority"`
	Enabled                 bool                 `json:"enabled" bson:"enabled"`
	UserId                  primitive.ObjectID   `json:"user_id" bson:"user_id"`
}
