package models

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Task struct {
	Id         primitive.ObjectID   `json:"_id" bson:"_id"`
	SpiderId   primitive.ObjectID   `json:"spider_id" bson:"spider_id"`
	Status     string               `json:"status" bson:"status"`
	NodeId     primitive.ObjectID   `json:"node_id" bson:"node_id"`
	Cmd        string               `json:"cmd" bson:"cmd"`
	Param      string               `json:"param" bson:"param"`
	Error      string               `json:"error" bson:"error"`
	Pid        int                  `json:"pid" bson:"pid"`
	ScheduleId primitive.ObjectID   `json:"schedule_id" bson:"schedule_id"` // Schedule.Id
	Type       string               `json:"type" bson:"type"`
	Mode       string               `json:"mode" bson:"mode"`           // running mode of Task
	NodeIds    []primitive.ObjectID `json:"node_ids" bson:"node_ids"`   // list of Node.Id
	ParentId   primitive.ObjectID   `json:"parent_id" bson:"parent_id"` // parent Task.Id if it'Spider a sub-task
	Priority   int                  `json:"priority" bson:"priority"`
	Stat       *TaskStat            `json:"stat,omitempty" bson:"-"`
	HasSub     bool                 `json:"has_sub" json:"has_sub"` // whether to have sub-tasks
	SubTasks   []Task               `json:"sub_tasks,omitempty" bson:"-"`
	Spider     *Spider              `json:"spider,omitempty" bson:"-"`
	UserId     primitive.ObjectID   `json:"-" bson:"-"`
	CreateTs   time.Time            `json:"create_ts" bson:"create_ts"`
}

func (t *Task) GetId() (id primitive.ObjectID) {
	return t.Id
}

func (t *Task) SetId(id primitive.ObjectID) {
	t.Id = id
}

func (t *Task) GetNodeId() (id primitive.ObjectID) {
	return t.NodeId
}

func (t *Task) SetNodeId(id primitive.ObjectID) {
	t.NodeId = id
}

func (t *Task) GetNodeIds() (ids []primitive.ObjectID) {
	return t.NodeIds
}

func (t *Task) GetStatus() (status string) {
	return t.Status
}

func (t *Task) SetStatus(status string) {
	t.Status = status
}

func (t *Task) GetError() (error string) {
	return t.Error
}

func (t *Task) SetError(error string) {
	t.Error = error
}

func (t *Task) GetPid() (pid int) {
	return t.Pid
}

func (t *Task) SetPid(pid int) {
	t.Pid = pid
}

func (t *Task) GetSpiderId() (id primitive.ObjectID) {
	return t.SpiderId
}

func (t *Task) GetType() (ty string) {
	return t.Type
}

func (t *Task) GetCmd() (cmd string) {
	return t.Cmd
}

func (t *Task) GetParam() (param string) {
	return t.Param
}

func (t *Task) GetPriority() (p int) {
	return t.Priority
}

func (t *Task) GetUserId() (id primitive.ObjectID) {
	return t.UserId
}

func (t *Task) SetUserId(id primitive.ObjectID) {
	t.UserId = id
}

type TaskList []Task

func (l *TaskList) GetModels() (res []interfaces.Model) {
	for i := range *l {
		d := (*l)[i]
		res = append(res, &d)
	}
	return res
}

type TaskDailyItem struct {
	Date               string  `json:"date" bson:"_id"`
	TaskCount          int     `json:"task_count" bson:"task_count"`
	AvgRuntimeDuration float64 `json:"avg_runtime_duration" bson:"avg_runtime_duration"`
}
