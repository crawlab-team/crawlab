package interfaces

import "go.mongodb.org/mongo-driver/bson/primitive"

type TaskStatsService interface {
	TaskBaseService
	InsertData(id primitive.ObjectID, records ...interface{}) (err error)
	InsertLogs(id primitive.ObjectID, logs ...string) (err error)
}
