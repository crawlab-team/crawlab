package mongo

import (
	"github.com/crawlab-team/go-trace"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetMongoDb(dbName string, opts ...DbOption) (db *mongo.Database) {
	if dbName == "" {
		dbName = viper.GetString("mongo.db")
	}
	if dbName == "" {
		dbName = "test"
	}

	_opts := &DbOptions{}
	for _, op := range opts {
		op(_opts)
	}

	var c *mongo.Client
	if _opts.client == nil {
		var err error
		c, err = GetMongoClient()
		if err != nil {
			trace.PrintError(err)
			return nil
		}
	} else {
		c = _opts.client
	}

	return c.Database(dbName, nil)
}
