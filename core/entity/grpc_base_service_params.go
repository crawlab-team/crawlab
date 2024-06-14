package entity

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GrpcBaseServiceParams struct {
	Query       bson.M             `json:"q"`
	Id          primitive.ObjectID `json:"id"`
	Update      bson.M             `json:"u"`
	Doc         interfaces.Model   `json:"d"`
	Fields      []string           `json:"f"`
	FindOptions *mongo.FindOptions `json:"o"`
	Docs        []interface{}      `json:"dl"`
	User        interfaces.User    `json:"U"`
}

func (params *GrpcBaseServiceParams) Value() interface{} {
	return params
}
