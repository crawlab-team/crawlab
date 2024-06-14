package models

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Schedule struct {
	Id          primitive.ObjectID   `json:"_id" bson:"_id"`
	Name        string               `json:"name" bson:"name"`
	Description string               `json:"description" bson:"description"`
	SpiderId    primitive.ObjectID   `json:"spider_id" bson:"spider_id"`
	Cron        string               `json:"cron" bson:"cron"`
	EntryId     cron.EntryID         `json:"entry_id" bson:"entry_id"`
	Cmd         string               `json:"cmd" bson:"cmd"`
	Param       string               `json:"param" bson:"param"`
	Mode        string               `json:"mode" bson:"mode"`
	NodeIds     []primitive.ObjectID `json:"node_ids" bson:"node_ids"`
	Priority    int                  `json:"priority" bson:"priority"`
	Enabled     bool                 `json:"enabled" bson:"enabled"`
	UserId      primitive.ObjectID   `json:"user_id" bson:"user_id"`
}

func (s *Schedule) GetId() (id primitive.ObjectID) {
	return s.Id
}

func (s *Schedule) SetId(id primitive.ObjectID) {
	s.Id = id
}

func (s *Schedule) GetEnabled() (enabled bool) {
	return s.Enabled
}

func (s *Schedule) SetEnabled(enabled bool) {
	s.Enabled = enabled
}

func (s *Schedule) GetEntryId() (id cron.EntryID) {
	return s.EntryId
}

func (s *Schedule) SetEntryId(id cron.EntryID) {
	s.EntryId = id
}

func (s *Schedule) GetCron() (c string) {
	return s.Cron
}

func (s *Schedule) SetCron(c string) {
	s.Cron = c
}

func (s *Schedule) GetSpiderId() (id primitive.ObjectID) {
	return s.SpiderId
}

func (s *Schedule) SetSpiderId(id primitive.ObjectID) {
	s.SpiderId = id
}

func (s *Schedule) GetMode() (mode string) {
	return s.Mode
}

func (s *Schedule) SetMode(mode string) {
	s.Mode = mode
}

func (s *Schedule) GetNodeIds() (ids []primitive.ObjectID) {
	return s.NodeIds
}

func (s *Schedule) SetNodeIds(ids []primitive.ObjectID) {
	s.NodeIds = ids
}

func (s *Schedule) GetCmd() (cmd string) {
	return s.Cmd
}

func (s *Schedule) SetCmd(cmd string) {
	s.Cmd = cmd
}

func (s *Schedule) GetParam() (param string) {
	return s.Param
}

func (s *Schedule) SetParam(param string) {
	s.Param = param
}

func (s *Schedule) GetPriority() (p int) {
	return s.Priority
}

func (s *Schedule) SetPriority(p int) {
	s.Priority = p
}

type ScheduleList []Schedule

func (l *ScheduleList) GetModels() (res []interfaces.Model) {
	for i := range *l {
		d := (*l)[i]
		res = append(res, &d)
	}
	return res
}
