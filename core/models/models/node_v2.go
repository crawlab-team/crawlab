package models

import (
	"time"
)

type NodeV2 struct {
	any                 `collection:"nodes"`
	BaseModelV2[NodeV2] `bson:",inline"`
	Key                 string    `json:"key" bson:"key"`
	Name                string    `json:"name" bson:"name"`
	Ip                  string    `json:"ip" bson:"ip"`
	Port                string    `json:"port" bson:"port"`
	Mac                 string    `json:"mac" bson:"mac"`
	Hostname            string    `json:"hostname" bson:"hostname"`
	Description         string    `json:"description" bson:"description"`
	IsMaster            bool      `json:"is_master" bson:"is_master"`
	Status              string    `json:"status" bson:"status"`
	Enabled             bool      `json:"enabled" bson:"enabled"`
	Active              bool      `json:"active" bson:"active"`
	ActiveAt            time.Time `json:"active_at" bson:"active_ts"`
	AvailableRunners    int       `json:"available_runners" bson:"available_runners"`
	MaxRunners          int       `json:"max_runners" bson:"max_runners"`
}
