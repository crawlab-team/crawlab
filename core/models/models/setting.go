package models

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Setting struct {
	Id    primitive.ObjectID `json:"_id" bson:"_id"`
	Key   string             `json:"key" bson:"key"`
	Value bson.M             `json:"value" bson:"value"`
}

func (s *Setting) GetId() (id primitive.ObjectID) {
	return s.Id
}

func (s *Setting) SetId(id primitive.ObjectID) {
	s.Id = id
}

type SettingList []Setting

func (l *SettingList) GetModels() (res []interfaces.Model) {
	for i := range *l {
		d := (*l)[i]
		res = append(res, &d)
	}
	return res
}
