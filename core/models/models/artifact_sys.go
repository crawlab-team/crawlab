package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ArtifactSys struct {
	CreateTs  time.Time          `json:"create_ts" bson:"create_ts"`
	CreateUid primitive.ObjectID `json:"create_uid" bson:"create_uid"`
	UpdateTs  time.Time          `json:"update_ts" bson:"update_ts"`
	UpdateUid primitive.ObjectID `json:"update_uid" bson:"update_uid"`
	DeleteTs  time.Time          `json:"delete_ts" bson:"delete_ts"`
	DeleteUid primitive.ObjectID `json:"delete_uid" bson:"delete_uid"`
}

func (sys *ArtifactSys) GetCreateTs() time.Time {
	return sys.CreateTs
}

func (sys *ArtifactSys) SetCreateTs(ts time.Time) {
	sys.CreateTs = ts
}

func (sys *ArtifactSys) GetUpdateTs() time.Time {
	return sys.UpdateTs
}

func (sys *ArtifactSys) SetUpdateTs(ts time.Time) {
	sys.UpdateTs = ts
}

func (sys *ArtifactSys) GetDeleteTs() time.Time {
	return sys.DeleteTs
}

func (sys *ArtifactSys) SetDeleteTs(ts time.Time) {
	sys.DeleteTs = ts
}

func (sys *ArtifactSys) GetCreateUid() primitive.ObjectID {
	return sys.CreateUid
}

func (sys *ArtifactSys) SetCreateUid(id primitive.ObjectID) {
	sys.CreateUid = id
}

func (sys *ArtifactSys) GetUpdateUid() primitive.ObjectID {
	return sys.UpdateUid
}

func (sys *ArtifactSys) SetUpdateUid(id primitive.ObjectID) {
	sys.UpdateUid = id
}

func (sys *ArtifactSys) GetDeleteUid() primitive.ObjectID {
	return sys.DeleteUid
}

func (sys *ArtifactSys) SetDeleteUid(id primitive.ObjectID) {
	sys.DeleteUid = id
}
