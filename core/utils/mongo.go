package utils

import (
	"github.com/crawlab-team/crawlab/db/generic"
	"github.com/crawlab-team/crawlab/db/mongo"
	"go.mongodb.org/mongo-driver/bson"
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
