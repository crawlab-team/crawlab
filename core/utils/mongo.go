package utils

import (
	"context"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/db/generic"
	"github.com/crawlab-team/crawlab/db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"time"
)

func GetMongoQuery(query generic.ListQuery) (res bson.M) {
	res = bson.M{}
	for _, c := range query {
		switch c.Op {
		case generic.OpEqual:
			res[c.Key] = c.Value
		default:
			res[c.Key] = bson.M{
				c.Op: c.Value,
			}
		}
	}
	return res
}

func GetMongoOpts(opts *generic.ListOptions) (res *mongo.FindOptions) {
	var sort bson.D
	for _, s := range opts.Sort {
		direction := 1
		if s.Direction == generic.SortDirectionAsc {
			direction = 1
		} else if s.Direction == generic.SortDirectionDesc {
			direction = -1
		}
		sort = append(sort, bson.E{Key: s.Key, Value: direction})
	}
	return &mongo.FindOptions{
		Skip:  opts.Skip,
		Limit: opts.Limit,
		Sort:  sort,
	}
}

func GetMongoClient(ds *models.DataSource) (c *mongo2.Client, err error) {
	return getMongoClient(context.Background(), ds)
}

func GetMongoClientWithTimeout(ds *models.DataSource, timeout time.Duration) (c *mongo2.Client, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return getMongoClient(ctx, ds)
}

func getMongoClient(ctx context.Context, ds *models.DataSource) (c *mongo2.Client, err error) {
	// normalize settings
	if ds.Host == "" {
		ds.Host = constants.DefaultHost
	}
	if ds.Port == "" {
		ds.Port = constants.DefaultMongoPort
	}

	// options
	var opts []mongo.ClientOption
	opts = append(opts, mongo.WithContext(ctx))
	opts = append(opts, mongo.WithUri(ds.Url))
	opts = append(opts, mongo.WithHost(ds.Host))
	opts = append(opts, mongo.WithPort(ds.Port))
	opts = append(opts, mongo.WithDb(ds.Database))
	opts = append(opts, mongo.WithUsername(ds.Username))
	opts = append(opts, mongo.WithPassword(ds.Password))
	opts = append(opts, mongo.WithHosts(ds.Hosts))

	// extra
	if ds.Extra != nil {
		// auth source
		authSource, ok := ds.Extra["auth_source"]
		if ok {
			opts = append(opts, mongo.WithAuthSource(authSource))
		}

		// auth mechanism
		authMechanism, ok := ds.Extra["auth_mechanism"]
		if ok {
			opts = append(opts, mongo.WithAuthMechanism(authMechanism))
		}
	}

	// client
	return mongo.GetMongoClient(opts...)
}
