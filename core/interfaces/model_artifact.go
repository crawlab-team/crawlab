package interfaces

import "go.mongodb.org/mongo-driver/bson/primitive"

type ModelArtifact interface {
	Model
	GetSys() (sys ModelArtifactSys)
	GetTagIds() (ids []primitive.ObjectID)
	SetTagIds(ids []primitive.ObjectID)
	SetObj(obj Model)
	SetDel(del bool)
}
