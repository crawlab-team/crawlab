package models

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Token struct {
	Id    primitive.ObjectID `json:"_id" bson:"_id"`
	Name  string             `json:"name" bson:"name"`
	Token string             `json:"token" bson:"token"`
}

func (t *Token) GetId() (id primitive.ObjectID) {
	return t.Id
}

func (t *Token) SetId(id primitive.ObjectID) {
	t.Id = id
}

type TokenList []Token

func (l *TokenList) GetModels() (res []interfaces.Model) {
	for i := range *l {
		d := (*l)[i]
		res = append(res, &d)
	}
	return res
}
