package models

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExtraValue struct {
	Id       primitive.ObjectID `json:"_id" bson:"_id"`
	ObjectId primitive.ObjectID `json:"oid" bson:"oid"`
	Model    string             `json:"model" bson:"m"`
	Type     string             `json:"type" bson:"t"`
	Value    interface{}        `json:"value" bson:"v"`
}

func (ev *ExtraValue) GetId() (id primitive.ObjectID) {
	return ev.Id
}

func (ev *ExtraValue) SetId(id primitive.ObjectID) {
	ev.Id = id
}

func (ev *ExtraValue) GetValue() (v interface{}) {
	return ev.Value
}

func (ev *ExtraValue) SetValue(v interface{}) {
	ev.Value = v
}

func (ev *ExtraValue) GetObjectId() (oid primitive.ObjectID) {
	return ev.ObjectId
}

func (ev *ExtraValue) SetObjectId(oid primitive.ObjectID) {
	ev.ObjectId = oid
}

func (ev *ExtraValue) GetModel() (m string) {
	return ev.Model
}

func (ev *ExtraValue) SetModel(m string) {
	ev.Model = m
}

func (ev *ExtraValue) GetType() (t string) {
	return ev.Type
}

func (ev *ExtraValue) SetType(t string) {
	ev.Type = t
}

type ExtraValueList []ExtraValue

func (l *ExtraValueList) GetModels() (res []interfaces.Model) {
	for i := range *l {
		d := (*l)[i]
		res = append(res, &d)
	}
	return res
}
