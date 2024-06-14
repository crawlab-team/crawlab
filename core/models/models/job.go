package models

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Job struct {
	Id     primitive.ObjectID `bson:"_id" json:"_id"`
	TaskId primitive.ObjectID `bson:"task_id" json:"task_id"`
}

func (j *Job) GetId() (id primitive.ObjectID) {
	return j.Id
}

func (j *Job) SetId(id primitive.ObjectID) {
	j.Id = id
}

type JobList []Job

func (l *JobList) GetModels() (res []interfaces.Model) {
	for i := range *l {
		d := (*l)[i]
		res = append(res, &d)
	}
	return res
}
