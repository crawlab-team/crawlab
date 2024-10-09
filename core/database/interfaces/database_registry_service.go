package interfaces

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DatabaseRegistryService interface {
	Start()
	CheckStatus()
	GetDatabaseService(id primitive.ObjectID) (res DatabaseService, err error)
}
