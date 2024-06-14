package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BaseModel struct {
	Id primitive.ObjectID `json:"_id" bson:"_id"`
}

func (d *BaseModel) GetId() (id primitive.ObjectID) {
	return d.Id
}
