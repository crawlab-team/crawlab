package interfaces

import "go.mongodb.org/mongo-driver/bson/primitive"

type ResultServiceRegistry interface {
	Register(key string, fn ResultServiceRegistryFn)
	Unregister(key string)
	Get(key string) (fn ResultServiceRegistryFn)
}

type ResultServiceRegistryFn func(colId primitive.ObjectID, dsId primitive.ObjectID) (ResultService, error)
