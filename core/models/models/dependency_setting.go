package models

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type DependencySetting struct {
	Id           primitive.ObjectID `json:"_id" bson:"_id"`
	Key          string             `json:"key" bson:"key"`
	Name         string             `json:"name" bson:"name"`
	Description  string             `json:"description" bson:"description"`
	Enabled      bool               `json:"enabled" bson:"enabled"`
	Cmd          string             `json:"cmd" bson:"cmd"`
	Proxy        string             `json:"proxy" bson:"proxy"`
	LastUpdateTs time.Time          `json:"last_update_ts" bson:"last_update_ts"`
}

func (j *DependencySetting) GetId() (id primitive.ObjectID) {
	return j.Id
}

func (j *DependencySetting) SetId(id primitive.ObjectID) {
	j.Id = id
}

type DependencySettingList []DependencySetting

func (l *DependencySettingList) GetModels() (res []interfaces.Model) {
	for i := range *l {
		d := (*l)[i]
		res = append(res, &d)
	}
	return res
}
