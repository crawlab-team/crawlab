package utils

import "github.com/google/uuid"

func NewUUIDString() (res string) {
	id, _ := uuid.NewUUID()
	return id.String()
}
