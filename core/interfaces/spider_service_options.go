package interfaces

import "go.mongodb.org/mongo-driver/bson/primitive"

type SpiderRunOptions struct {
	Mode       string               `json:"mode"`
	NodeIds    []primitive.ObjectID `json:"node_ids"`
	Cmd        string               `json:"cmd"`
	Param      string               `json:"param"`
	ScheduleId primitive.ObjectID   `json:"schedule_id"`
	Priority   int                  `json:"priority"`
	UserId     primitive.ObjectID   `json:"-"`
}

type SpiderCloneOptions struct {
	Name string
}
