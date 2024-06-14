package models

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Env struct {
	Name  string `json:"name" bson:"name"`
	Value string `json:"value" bson:"value"`
}

type Spider struct {
	Id           primitive.ObjectID   `json:"_id" bson:"_id"`                       // spider id
	Name         string               `json:"name" bson:"name"`                     // spider name
	Type         string               `json:"type" bson:"type"`                     // spider type
	ColId        primitive.ObjectID   `json:"col_id" bson:"col_id"`                 // data collection id
	ColName      string               `json:"col_name,omitempty" bson:"-"`          // data collection name
	DataSourceId primitive.ObjectID   `json:"data_source_id" bson:"data_source_id"` // data source id
	DataSource   *DataSource          `json:"data_source,omitempty" bson:"-"`       // data source
	Description  string               `json:"description" bson:"description"`       // description
	ProjectId    primitive.ObjectID   `json:"project_id" bson:"project_id"`         // Project.Id
	Mode         string               `json:"mode" bson:"mode"`                     // default Task.Mode
	NodeIds      []primitive.ObjectID `json:"node_ids" bson:"node_ids"`             // default Task.NodeIds
	Stat         *SpiderStat          `json:"stat,omitempty" bson:"-"`

	// execution
	Cmd         string `json:"cmd" bson:"cmd"`     // execute command
	Param       string `json:"param" bson:"param"` // default task param
	Priority    int    `json:"priority" bson:"priority"`
	AutoInstall bool   `json:"auto_install" bson:"auto_install"`

	// settings
	IncrementalSync bool `json:"incremental_sync" bson:"incremental_sync"` // whether to incrementally sync files
}

func (s *Spider) GetId() (id primitive.ObjectID) {
	return s.Id
}

func (s *Spider) SetId(id primitive.ObjectID) {
	s.Id = id
}

func (s *Spider) GetName() (name string) {
	return s.Name
}

func (s *Spider) SetName(name string) {
	s.Name = name
}

func (s *Spider) GetDescription() (description string) {
	return s.Description
}

func (s *Spider) SetDescription(description string) {
	s.Description = description
}

func (s *Spider) GetType() (ty string) {
	return s.Type
}

func (s *Spider) GetMode() (mode string) {
	return s.Mode
}

func (s *Spider) SetMode(mode string) {
	s.Mode = mode
}

func (s *Spider) GetNodeIds() (ids []primitive.ObjectID) {
	return s.NodeIds
}

func (s *Spider) SetNodeIds(ids []primitive.ObjectID) {
	s.NodeIds = ids
}

func (s *Spider) GetCmd() (cmd string) {
	return s.Cmd
}

func (s *Spider) SetCmd(cmd string) {
	s.Cmd = cmd
}

func (s *Spider) GetParam() (param string) {
	return s.Param
}

func (s *Spider) SetParam(param string) {
	s.Param = param
}

func (s *Spider) GetPriority() (p int) {
	return s.Priority
}

func (s *Spider) SetPriority(p int) {
	s.Priority = p
}

func (s *Spider) GetColId() (id primitive.ObjectID) {
	return s.ColId
}

func (s *Spider) SetColId(id primitive.ObjectID) {
	s.ColId = id
}

func (s *Spider) GetIncrementalSync() (incrementalSync bool) {
	return s.IncrementalSync
}

func (s *Spider) SetIncrementalSync(incrementalSync bool) {
	s.IncrementalSync = incrementalSync
}

func (s *Spider) GetAutoInstall() (autoInstall bool) {
	return s.AutoInstall
}

func (s *Spider) SetAutoInstall(autoInstall bool) {
	s.AutoInstall = autoInstall
}

type SpiderList []Spider

func (l *SpiderList) GetModels() (res []interfaces.Model) {
	for i := range *l {
		d := (*l)[i]
		res = append(res, &d)
	}
	return res
}
