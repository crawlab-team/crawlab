package interfaces

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ModelBaseService interface {
	GetModelId() (id ModelId)
	SetModelId(id ModelId)
	GetById(id primitive.ObjectID) (res Model, err error)
	Get(query bson.M, opts *mongo.FindOptions) (res Model, err error)
	GetList(query bson.M, opts *mongo.FindOptions) (res List, err error)
	DeleteById(id primitive.ObjectID, args ...interface{}) (err error)
	Delete(query bson.M, args ...interface{}) (err error)
	DeleteList(query bson.M, args ...interface{}) (err error)
	ForceDeleteList(query bson.M, args ...interface{}) (err error)
	UpdateById(id primitive.ObjectID, update bson.M, args ...interface{}) (err error)
	Update(query bson.M, update bson.M, fields []string, args ...interface{}) (err error)
	UpdateDoc(query bson.M, doc Model, fields []string, args ...interface{}) (err error)
	Insert(u User, docs ...interface{}) (err error)
	Count(query bson.M) (total int, err error)
}

type ModelService interface {
	GetBaseService(id ModelId) (svc ModelBaseService)
}
