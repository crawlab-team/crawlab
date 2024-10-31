package interfaces

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Model interface {
	GetId() (id primitive.ObjectID)
	SetId(id primitive.ObjectID)
	SetCreated(by primitive.ObjectID)
	SetUpdated(by primitive.ObjectID)
}
