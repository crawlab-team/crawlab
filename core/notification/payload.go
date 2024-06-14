package notification

import "go.mongodb.org/mongo-driver/bson/primitive"

type SendPayload struct {
	TaskId primitive.ObjectID `json:"task_id"`
	Data   string             `json:"data"`
}
