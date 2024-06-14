package interfaces

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GrpcClientModelEnvironmentService interface {
	ModelBaseService
	GetEnvironmentById(id primitive.ObjectID) (s Environment, err error)
	GetEnvironment(query bson.M, opts *mongo.FindOptions) (s Environment, err error)
	GetEnvironmentList(query bson.M, opts *mongo.FindOptions) (res []Environment, err error)
}
