package models

import (
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Result bson.M

func (r *Result) GetId() (id primitive.ObjectID) {
	res, ok := r.Value()["_id"]
	if ok {
		id, ok = res.(primitive.ObjectID)
		if ok {
			return id
		}
	}
	return id
}

func (r *Result) SetId(id primitive.ObjectID) {
	(*r)["_id"] = id
}

func (r *Result) Value() map[string]interface{} {
	return *r
}

func (r *Result) SetValue(key string, value interface{}) {
	(*r)[key] = value
}

func (r *Result) GetValue(key string) (value interface{}) {
	return (*r)[key]
}

func (r *Result) GetTaskId() (id primitive.ObjectID) {
	res := r.GetValue(constants.TaskKey)
	if res == nil {
		return id
	}
	id, _ = res.(primitive.ObjectID)
	return id
}

func (r *Result) SetTaskId(id primitive.ObjectID) {
	r.SetValue(constants.TaskKey, id)
}

type ResultList []Result

func (l *ResultList) GetModels() (res []interfaces.Model) {
	for i := range *l {
		d := (*l)[i]
		res = append(res, &d)
	}
	return res
}
