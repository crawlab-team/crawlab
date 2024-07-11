package models

import (
	"time"
)

type DependencySettingV2 struct {
	any                            `collection:"dependency_settings"`
	BaseModelV2[DependencySetting] `bson:",inline"`
	Key                            string    `json:"key" bson:"key"`
	Name                           string    `json:"name" bson:"name"`
	Description                    string    `json:"description" bson:"description"`
	Enabled                        bool      `json:"enabled" bson:"enabled"`
	Cmd                            string    `json:"cmd" bson:"cmd"`
	Proxy                          string    `json:"proxy" bson:"proxy"`
	LastUpdateTs                   time.Time `json:"last_update_ts" bson:"last_update_ts"`
}
