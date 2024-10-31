package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskV2 struct {
	any                 `collection:"tasks"`
	BaseModelV2[TaskV2] `bson:",inline"`
	SpiderId            primitive.ObjectID   `json:"spider_id" bson:"spider_id"`
	Status              string               `json:"status" bson:"status"`
	NodeId              primitive.ObjectID   `json:"node_id" bson:"node_id"`
	Cmd                 string               `json:"cmd" bson:"cmd"`
	Param               string               `json:"param" bson:"param"`
	Error               string               `json:"error" bson:"error"`
	Pid                 int                  `json:"pid" bson:"pid"`
	ScheduleId          primitive.ObjectID   `json:"schedule_id" bson:"schedule_id"`
	Type                string               `json:"type" bson:"type"`
	Mode                string               `json:"mode" bson:"mode"`
	NodeIds             []primitive.ObjectID `json:"node_ids" bson:"node_ids"`
	Priority            int                  `json:"priority" bson:"priority"`
	Stat                *TaskStatV2          `json:"stat,omitempty" bson:"-"`
	Spider              *SpiderV2            `json:"spider,omitempty" bson:"-"`
	Schedule            *ScheduleV2          `json:"schedule,omitempty" bson:"-"`
	Node                *NodeV2              `json:"node,omitempty" bson:"-"`
	UserId              primitive.ObjectID   `json:"-" bson:"-"`
}
