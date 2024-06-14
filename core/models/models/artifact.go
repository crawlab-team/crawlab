package models

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Artifact struct {
	Id     primitive.ObjectID   `bson:"_id" json:"_id"`
	Col    string               `bson:"_col" json:"_col"`
	Del    bool                 `bson:"_del" json:"_del"`
	TagIds []primitive.ObjectID `bson:"_tid" json:"_tid"`
	Sys    *ArtifactSys         `bson:"_sys" json:"_sys"`
	Obj    interface{}          `bson:"_obj" json:"_obj"`
}

func (a *Artifact) GetId() (id primitive.ObjectID) {
	return a.Id
}

func (a *Artifact) SetId(id primitive.ObjectID) {
	a.Id = id
}

func (a *Artifact) GetSys() (sys interfaces.ModelArtifactSys) {
	if a.Sys == nil {
		a.Sys = &ArtifactSys{}
	}
	return a.Sys
}

func (a *Artifact) GetTagIds() (ids []primitive.ObjectID) {
	return a.TagIds
}

func (a *Artifact) SetTagIds(ids []primitive.ObjectID) {
	a.TagIds = ids
}

func (a *Artifact) SetObj(obj interfaces.Model) {
	a.Obj = obj
}

func (a *Artifact) SetDel(del bool) {
	a.Del = del
}

type ArtifactList []Artifact

func (l *ArtifactList) GetModels() (res []interfaces.Model) {
	for i := range *l {
		d := (*l)[i]
		res = append(res, &d)
	}
	return res
}
