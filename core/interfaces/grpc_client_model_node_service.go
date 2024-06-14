package interfaces

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GrpcClientModelNodeService interface {
	ModelBaseService
	GetNodeById(id primitive.ObjectID) (n Node, err error)
	GetNode(query bson.M, opts *mongo.FindOptions) (n Node, err error)
	GetNodeByKey(key string) (n Node, err error)
	GetNodeList(query bson.M, opts *mongo.FindOptions) (res []Node, err error)
}
