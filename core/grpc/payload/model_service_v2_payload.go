package payload

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ModelServiceV2Payload struct {
	Type        string             `json:"type,omitempty"`
	Id          primitive.ObjectID `json:"_id,omitempty"`
	Query       bson.M             `json:"query,omitempty"`
	FindOptions *mongo.FindOptions `json:"find_options,omitempty"`
	Model       any                `json:"model,omitempty"`
	Update      bson.M             `json:"update,omitempty"`
	Models      []any              `json:"models,omitempty"`
}
