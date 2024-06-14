package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SpiderV2 struct {
	any                   `collection:"spiders"` // spider id
	BaseModelV2[SpiderV2] `bson:",inline"`
	Name                  string               `json:"name" bson:"name"`                     // spider name
	Type                  string               `json:"type" bson:"type"`                     // spider type
	ColId                 primitive.ObjectID   `json:"col_id" bson:"col_id"`                 // data collection id
	ColName               string               `json:"col_name,omitempty" bson:"-"`          // data collection name
	DataSourceId          primitive.ObjectID   `json:"data_source_id" bson:"data_source_id"` // data source id
	DataSource            *DataSourceV2        `json:"data_source,omitempty" bson:"-"`       // data source
	Description           string               `json:"description" bson:"description"`       // description
	ProjectId             primitive.ObjectID   `json:"project_id" bson:"project_id"`         // Project.Id
	Mode                  string               `json:"mode" bson:"mode"`                     // default Task.Mode
	NodeIds               []primitive.ObjectID `json:"node_ids" bson:"node_ids"`             // default Task.NodeIds
	Stat                  *SpiderStatV2        `json:"stat,omitempty" bson:"-"`

	// execution
	Cmd         string `json:"cmd" bson:"cmd"`     // execute command
	Param       string `json:"param" bson:"param"` // default task param
	Priority    int    `json:"priority" bson:"priority"`
	AutoInstall bool   `json:"auto_install" bson:"auto_install"`

	// settings
	IncrementalSync bool `json:"incremental_sync" bson:"incremental_sync"` // whether to incrementally sync files
}
