package utils

import (
	"github.com/crawlab-team/crawlab-db/generic"
	"github.com/upper/db/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetSqlQuery(query generic.ListQuery) (res db.Cond) {
	res = db.Cond{}
	for _, c := range query {
		switch c.Value.(type) {
		case primitive.ObjectID:
			c.Value = c.Value.(primitive.ObjectID).Hex()
		}
		switch c.Op {
		case generic.OpEqual:
			res[c.Key] = c.Value
		default:
			res[c.Key] = db.Cond{
				c.Op: c.Value,
			}
		}
	}
	// TODO: sort
	return res
}
