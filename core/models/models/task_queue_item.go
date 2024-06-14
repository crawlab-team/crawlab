package models

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskQueueItem struct {
	Id       primitive.ObjectID `json:"_id" bson:"_id"`
	Priority int                `json:"p" bson:"p"`
	NodeId   primitive.ObjectID `json:"nid,omitempty" bson:"nid,omitempty"`
}

func (t *TaskQueueItem) GetId() (id primitive.ObjectID) {
	return t.Id
}

func (t *TaskQueueItem) SetId(id primitive.ObjectID) {
	t.Id = id
}

type TaskQueueItemList []TaskQueueItem

func (l *TaskQueueItemList) GetModels() (res []interfaces.Model) {
	for i := range *l {
		d := (*l)[i]
		res = append(res, &d)
	}
	return res
}
