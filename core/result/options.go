package result

import "go.mongodb.org/mongo-driver/bson/primitive"

type Option func(opts *Options)

type Options struct {
	registryKey string             // registry key
	SpiderId    primitive.ObjectID // data source id
}

func WithRegistryKey(key string) Option {
	return func(opts *Options) {
		opts.registryKey = key
	}
}
