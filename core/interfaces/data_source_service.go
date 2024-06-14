package interfaces

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type DataSourceService interface {
	ChangePassword(id primitive.ObjectID, password string) (err error)
	Monitor()
	CheckStatus(id primitive.ObjectID) (err error)
	SetTimeout(duration time.Duration)
	SetMonitorInterval(duration time.Duration)
}
