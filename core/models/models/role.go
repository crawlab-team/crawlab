package models

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role struct {
	Id          primitive.ObjectID `json:"_id" bson:"_id"`
	Key         string             `json:"key" bson:"key"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
}

func (r *Role) GetId() (id primitive.ObjectID) {
	return r.Id
}

func (r *Role) SetId(id primitive.ObjectID) {
	r.Id = id
}

func (r *Role) GetKey() (key string) {
	return r.Key
}

func (r *Role) SetKey(key string) {
	r.Key = key
}

func (r *Role) GetName() (name string) {
	return r.Name
}

func (r *Role) SetName(name string) {
	r.Name = name
}

func (r *Role) GetDescription() (description string) {
	return r.Description
}

func (r *Role) SetDescription(description string) {
	r.Description = description
}

type RoleList []Role

func (l *RoleList) GetModels() (res []interfaces.Model) {
	for i := range *l {
		d := (*l)[i]
		res = append(res, &d)
	}
	return res
}
