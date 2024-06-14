package interfaces

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ResultServiceMongo interface {
	GetId() (id primitive.ObjectID)
	SetId(id primitive.ObjectID)
	List(query bson.M, opts *mongo.FindOptions) (results []Result, err error)
	Count(query bson.M) (total int, err error)
	Insert(docs ...interface{}) (err error)
}
