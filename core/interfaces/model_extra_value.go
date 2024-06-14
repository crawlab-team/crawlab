package interfaces

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExtraValue interface {
	Model
	GetValue() (v interface{})
	SetValue(v interface{})
	GetObjectId() (oid primitive.ObjectID)
	SetObjectId(oid primitive.ObjectID)
	GetModel() (m string)
	SetModel(m string)
	GetType() (t string)
	SetType(t string)
}
