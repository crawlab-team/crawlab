package models

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Permission struct {
	Id          primitive.ObjectID `json:"_id" bson:"_id"`
	Key         string             `json:"key" bson:"key"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Type        string             `json:"type" bson:"type"`
	Target      []string           `json:"target" bson:"target"`
	Allow       []string           `json:"allow" bson:"allow"`
	Deny        []string           `json:"deny" bson:"deny"`
}

func (p *Permission) GetId() (id primitive.ObjectID) {
	return p.Id
}

func (p *Permission) SetId(id primitive.ObjectID) {
	p.Id = id
}

func (p *Permission) GetKey() (key string) {
	return p.Key
}

func (p *Permission) SetKey(key string) {
	p.Key = key
}

func (p *Permission) GetName() (name string) {
	return p.Name
}

func (p *Permission) SetName(name string) {
	p.Name = name
}

func (p *Permission) GetDescription() (description string) {
	return p.Description
}

func (p *Permission) SetDescription(description string) {
	p.Description = description
}

func (p *Permission) GetType() (t string) {
	return p.Type
}

func (p *Permission) SetType(t string) {
	p.Type = t
}

func (p *Permission) GetTarget() (target []string) {
	return p.Target
}

func (p *Permission) SetTarget(target []string) {
	p.Target = target
}

func (p *Permission) GetAllow() (include []string) {
	return p.Allow
}

func (p *Permission) SetAllow(include []string) {
	p.Allow = include
}

func (p *Permission) GetDeny() (exclude []string) {
	return p.Deny
}

func (p *Permission) SetDeny(exclude []string) {
	p.Deny = exclude
}

type PermissionList []Permission

func (l *PermissionList) GetModels() (res []interfaces.Model) {
	for i := range *l {
		d := (*l)[i]
		res = append(res, &d)
	}
	return res
}
