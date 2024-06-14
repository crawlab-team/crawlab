package models

import (
	"time"
)

type TaskStatV2 struct {
	any                     `collection:"task_stats"`
	BaseModelV2[TaskStatV2] `bson:",inline"`
	CreateTs                time.Time `json:"create_ts" bson:"create_ts,omitempty"`
	StartTs                 time.Time `json:"start_ts" bson:"start_ts,omitempty"`
	EndTs                   time.Time `json:"end_ts" bson:"end_ts,omitempty"`
	WaitDuration            int64     `json:"wait_duration" bson:"wait_duration,omitempty"`       // in millisecond
	RuntimeDuration         int64     `json:"runtime_duration" bson:"runtime_duration,omitempty"` // in millisecond
	TotalDuration           int64     `json:"total_duration" bson:"total_duration,omitempty"`     // in millisecond
	ResultCount             int64     `json:"result_count" bson:"result_count"`
	ErrorLogCount           int64     `json:"error_log_count" bson:"error_log_count"`
}
