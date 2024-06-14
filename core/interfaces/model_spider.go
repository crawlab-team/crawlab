package interfaces

import "go.mongodb.org/mongo-driver/bson/primitive"

type Spider interface {
	ModelWithNameDescription
	GetType() (ty string)
	GetMode() (mode string)
	SetMode(mode string)
	GetNodeIds() (ids []primitive.ObjectID)
	SetNodeIds(ids []primitive.ObjectID)
	GetCmd() (cmd string)
	SetCmd(cmd string)
	GetParam() (param string)
	SetParam(param string)
	GetPriority() (p int)
	SetPriority(p int)
	GetColId() (id primitive.ObjectID)
	SetColId(id primitive.ObjectID)
	GetIncrementalSync() (incrementalSync bool)
	SetIncrementalSync(incrementalSync bool)
	GetAutoInstall() (autoInstall bool)
	SetAutoInstall(autoInstall bool)
}
