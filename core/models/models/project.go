package models

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Project struct {
	Id          primitive.ObjectID `json:"_id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Spiders     int                `json:"spiders" bson:"-"`
}

func (p *Project) GetId() (id primitive.ObjectID) {
	return p.Id
}

func (p *Project) SetId(id primitive.ObjectID) {
	p.Id = id
}

func (p *Project) GetName() (name string) {
	return p.Name
}

func (p *Project) SetName(name string) {
	p.Name = name
}

func (p *Project) GetDescription() (description string) {
	return p.Description
}

func (p *Project) SetDescription(description string) {
	p.Description = description
}

type ProjectList []Project

func (l *ProjectList) GetModels() (res []interfaces.Model) {
	for i := range *l {
		d := (*l)[i]
		res = append(res, &d)
	}
	return res
}
