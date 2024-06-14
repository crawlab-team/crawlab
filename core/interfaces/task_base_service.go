package interfaces

import "go.mongodb.org/mongo-driver/bson/primitive"

type TaskBaseService interface {
	WithConfigPath
	Module
	SaveTask(t Task, status string) (err error)
	IsStopped() (res bool)
	GetQueue(nodeId primitive.ObjectID) (queue string)
}
