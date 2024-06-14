package interfaces

import "go.mongodb.org/mongo-driver/bson/primitive"

type Task interface {
	Model
	GetNodeId() (id primitive.ObjectID)
	SetNodeId(id primitive.ObjectID)
	GetNodeIds() (ids []primitive.ObjectID)
	GetStatus() (status string)
	SetStatus(status string)
	GetError() (error string)
	SetError(error string)
	GetPid() (pid int)
	SetPid(pid int)
	GetSpiderId() (id primitive.ObjectID)
	GetType() (ty string)
	GetCmd() (cmd string)
	GetParam() (param string)
	GetPriority() (p int)
	GetUserId() (id primitive.ObjectID)
	SetUserId(id primitive.ObjectID)
}
