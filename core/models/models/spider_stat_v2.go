package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SpiderStatV2 struct {
	any                       `collection:"spider_stats"`
	BaseModelV2[SpiderStatV2] `bson:",inline"`
	LastTaskId                primitive.ObjectID `json:"last_task_id" bson:"last_task_id,omitempty"`
	LastTask                  *TaskV2            `json:"last_task,omitempty" bson:"-"`
	Tasks                     int                `json:"tasks" bson:"tasks"`
	Results                   int                `json:"results" bson:"results"`
	WaitDuration              int64              `json:"wait_duration" bson:"wait_duration,omitempty"`       // in second
	RuntimeDuration           int64              `json:"runtime_duration" bson:"runtime_duration,omitempty"` // in second
	TotalDuration             int64              `json:"total_duration" bson:"total_duration,omitempty"`     // in second
	AverageWaitDuration       int64              `json:"average_wait_duration" bson:"-"`                     // in second
	AverageRuntimeDuration    int64              `json:"average_runtime_duration" bson:"-"`                  // in second
	AverageTotalDuration      int64              `json:"average_total_duration" bson:"-"`                    // in second
}
