package models

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tag struct {
	Id          primitive.ObjectID `json:"_id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Color       string             `json:"color" bson:"color"`
	Description string             `json:"description" bson:"description"`
	Col         string             `json:"col" bson:"col"`
}

func (t *Tag) GetId() (id primitive.ObjectID) {
	return t.Id
}

func (t *Tag) SetId(id primitive.ObjectID) {
	t.Id = id
}

func (t *Tag) GetName() (res string) {
	return t.Name
}

func (t *Tag) GetColor() (res string) {
	return t.Color
}

func (t *Tag) SetCol(col string) {
	t.Col = col
}

type TagList []Tag

func (l *TagList) GetModels() (res []interfaces.Model) {
	for i := range *l {
		d := (*l)[i]
		res = append(res, &d)
	}
	return res
}
