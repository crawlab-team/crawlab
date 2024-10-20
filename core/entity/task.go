package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StreamMessageTaskData struct {
	TaskId  primitive.ObjectID `json:"task_id"`
	Records []Result           `json:"data"`
	Logs    []string           `json:"logs"`
}
