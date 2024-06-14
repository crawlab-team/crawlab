package models

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type TaskStat struct {
	Id              primitive.ObjectID `json:"_id" bson:"_id"`
	CreateTs        time.Time          `json:"create_ts" bson:"create_ts,omitempty"`
	StartTs         time.Time          `json:"start_ts" bson:"start_ts,omitempty"`
	EndTs           time.Time          `json:"end_ts" bson:"end_ts,omitempty"`
	WaitDuration    int64              `json:"wait_duration" bson:"wait_duration,omitempty"`       // in millisecond
	RuntimeDuration int64              `json:"runtime_duration" bson:"runtime_duration,omitempty"` // in millisecond
	TotalDuration   int64              `json:"total_duration" bson:"total_duration,omitempty"`     // in millisecond
	ResultCount     int64              `json:"result_count" bson:"result_count"`
	ErrorLogCount   int64              `json:"error_log_count" bson:"error_log_count"`
}

func (s *TaskStat) GetId() (id primitive.ObjectID) {
	return s.Id
}

func (s *TaskStat) SetId(id primitive.ObjectID) {
	s.Id = id
}

func (s *TaskStat) GetCreateTs() (ts time.Time) {
	return s.CreateTs
}

func (s *TaskStat) SetCreateTs(ts time.Time) {
	s.CreateTs = ts
}

func (s *TaskStat) GetStartTs() (ts time.Time) {
	return s.StartTs
}

func (s *TaskStat) SetStartTs(ts time.Time) {
	s.StartTs = ts
}

func (s *TaskStat) GetEndTs() (ts time.Time) {
	return s.EndTs
}

func (s *TaskStat) SetEndTs(ts time.Time) {
	s.EndTs = ts
}

func (s *TaskStat) GetWaitDuration() (d int64) {
	return s.WaitDuration
}

func (s *TaskStat) SetWaitDuration(d int64) {
	s.WaitDuration = d
}

func (s *TaskStat) GetRuntimeDuration() (d int64) {
	return s.RuntimeDuration
}

func (s *TaskStat) SetRuntimeDuration(d int64) {
	s.RuntimeDuration = d
}

func (s *TaskStat) GetTotalDuration() (d int64) {
	return s.WaitDuration + s.RuntimeDuration
}

func (s *TaskStat) SetTotalDuration(d int64) {
	s.TotalDuration = d
}

func (s *TaskStat) GetResultCount() (c int64) {
	return s.ResultCount
}

func (s *TaskStat) SetResultCount(c int64) {
	s.ResultCount = c
}

func (s *TaskStat) GetErrorLogCount() (c int64) {
	return s.ErrorLogCount
}

func (s *TaskStat) SetErrorLogCount(c int64) {
	s.ErrorLogCount = c
}

type TaskStatList []TaskStat

func (l *TaskStatList) GetModels() (res []interfaces.Model) {
	for i := range *l {
		d := (*l)[i]
		res = append(res, &d)
	}
	return res
}
