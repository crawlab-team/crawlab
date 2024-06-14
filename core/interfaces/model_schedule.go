package interfaces

import (
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Schedule interface {
	Model
	GetEnabled() (enabled bool)
	SetEnabled(enabled bool)
	GetEntryId() (id cron.EntryID)
	SetEntryId(id cron.EntryID)
	GetCron() (c string)
	SetCron(c string)
	GetSpiderId() (id primitive.ObjectID)
	SetSpiderId(id primitive.ObjectID)
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
}
