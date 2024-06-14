package models

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Git struct {
	Id            primitive.ObjectID `json:"_id" bson:"_id"`
	Url           string             `json:"url" bson:"url"`
	AuthType      string             `json:"auth_type" bson:"auth_type"`
	Username      string             `json:"username" bson:"username"`
	Password      string             `json:"password" bson:"password"`
	CurrentBranch string             `json:"current_branch" bson:"current_branch"`
	AutoPull      bool               `json:"auto_pull" bson:"auto_pull"`
}

func (g *Git) GetId() (id primitive.ObjectID) {
	return g.Id
}

func (g *Git) SetId(id primitive.ObjectID) {
	g.Id = id
}

func (g *Git) GetUrl() (url string) {
	return g.Url
}

func (g *Git) SetUrl(url string) {
	g.Url = url
}

func (g *Git) GetAuthType() (authType string) {
	return g.AuthType
}

func (g *Git) SetAuthType(authType string) {
	g.AuthType = authType
}

func (g *Git) GetUsername() (username string) {
	return g.Username
}

func (g *Git) SetUsername(username string) {
	g.Username = username
}

func (g *Git) GetPassword() (password string) {
	return g.Password
}

func (g *Git) SetPassword(password string) {
	g.Password = password
}

func (g *Git) GetCurrentBranch() (currentBranch string) {
	return g.CurrentBranch
}

func (g *Git) SetCurrentBranch(currentBranch string) {
	g.CurrentBranch = currentBranch
}

func (g *Git) GetAutoPull() (autoPull bool) {
	return g.AutoPull
}

func (g *Git) SetAutoPull(autoPull bool) {
	g.AutoPull = autoPull
}

type GitList []Git

func (l *GitList) GetModels() (res []interfaces.Model) {
	for i := range *l {
		d := (*l)[i]
		res = append(res, &d)
	}
	return res
}
