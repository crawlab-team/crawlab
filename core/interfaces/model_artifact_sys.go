package interfaces

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ModelArtifactSys interface {
	GetCreateTs() time.Time
	SetCreateTs(ts time.Time)
	GetUpdateTs() time.Time
	SetUpdateTs(ts time.Time)
	GetDeleteTs() time.Time
	SetDeleteTs(ts time.Time)
	GetCreateUid() primitive.ObjectID
	SetCreateUid(id primitive.ObjectID)
	GetUpdateUid() primitive.ObjectID
	SetUpdateUid(id primitive.ObjectID)
	GetDeleteUid() primitive.ObjectID
	SetDeleteUid(id primitive.ObjectID)
}
