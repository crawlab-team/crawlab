package mongo

import "go.mongodb.org/mongo-driver/mongo"

type DbOption func(options *DbOptions)

type DbOptions struct {
	client *mongo.Client
}

func WithDbClient(c *mongo.Client) DbOption {
	return func(options *DbOptions) {
		options.client = c
	}
}
