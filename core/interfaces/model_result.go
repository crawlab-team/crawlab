package interfaces

import "go.mongodb.org/mongo-driver/bson/primitive"

type Result interface {
	Value() map[string]interface{}
	SetValue(key string, value interface{})
	GetValue(key string) (value interface{})
	GetTaskId() (id primitive.ObjectID)
	SetTaskId(id primitive.ObjectID)
}
