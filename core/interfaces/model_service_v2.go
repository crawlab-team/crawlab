package interfaces

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ModelServiceV2[T any] interface {
	GetById(id primitive.ObjectID) (model *T, err error)
	Get(query bson.M, options *mongo.FindOptions) (model *T, err error)
	GetList(query bson.M, options *mongo.FindOptions) (models []T, err error)
	DeleteById(id primitive.ObjectID) (err error)
	Delete(query bson.M) (err error)
	DeleteList(query bson.M) (err error)
	UpdateById(id primitive.ObjectID, update bson.M) (err error)
	UpdateOne(query bson.M, update bson.M) (err error)
	UpdateMany(query bson.M, update bson.M) (err error)
	ReplaceById(id primitive.ObjectID, model T) (err error)
	Replace(query bson.M, model T) (err error)
	InsertOne(model T) (id primitive.ObjectID, err error)
	InsertMany(models []T) (ids []primitive.ObjectID, err error)
	Count(query bson.M) (total int, err error)
	GetCol() (col *mongo.Col)
}
