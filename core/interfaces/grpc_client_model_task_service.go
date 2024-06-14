package interfaces

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GrpcClientModelTaskService interface {
	ModelBaseService
	GetTaskById(id primitive.ObjectID) (s Task, err error)
	GetTask(query bson.M, opts *mongo.FindOptions) (s Task, err error)
	GetTaskList(query bson.M, opts *mongo.FindOptions) (res []Task, err error)
}
