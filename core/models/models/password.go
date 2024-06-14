package models

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Password struct {
	Id       primitive.ObjectID `json:"_id" bson:"_id"`
	Password string             `json:"password" bson:"p"`
}

func (p *Password) GetId() (id primitive.ObjectID) {
	return p.Id
}

func (p *Password) SetId(id primitive.ObjectID) {
	p.Id = id
}

type PasswordList []Password

func (l *PasswordList) GetModels() (res []interfaces.Model) {
	for i := range *l {
		d := (*l)[i]
		res = append(res, &d)
	}
	return res
}
