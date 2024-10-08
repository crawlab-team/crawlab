package database

import (
	"github.com/crawlab-team/crawlab/core/database/interfaces"
)

var serviceInstance interfaces.DatabaseRegistryService

func SetDatabaseRegistryService(svc interfaces.DatabaseRegistryService) {
	serviceInstance = svc
}

func GetDatabaseRegistryService() interfaces.DatabaseRegistryService {
	return serviceInstance
}
