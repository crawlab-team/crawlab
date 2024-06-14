package interfaces

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GrpcClientModelTaskStatService interface {
	ModelBaseService
	GetTaskStatById(id primitive.ObjectID) (s TaskStat, err error)
	GetTaskStat(query bson.M, opts *mongo.FindOptions) (s TaskStat, err error)
	GetTaskStatList(query bson.M, opts *mongo.FindOptions) (res []TaskStat, err error)
}
