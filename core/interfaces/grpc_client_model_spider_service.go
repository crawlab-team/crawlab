package interfaces

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GrpcClientModelSpiderService interface {
	ModelBaseService
	GetSpiderById(id primitive.ObjectID) (s Spider, err error)
	GetSpider(query bson.M, opts *mongo.FindOptions) (s Spider, err error)
	GetSpiderList(query bson.M, opts *mongo.FindOptions) (res []Spider, err error)
}
