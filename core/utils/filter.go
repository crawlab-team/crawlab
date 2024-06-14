package utils

import (
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"go.mongodb.org/mongo-driver/bson"
)

// FilterToQuery Translate entity.Filter to bson.M
func FilterToQuery(f interfaces.Filter) (q bson.M, err error) {
	if f == nil || f.IsNil() {
		return nil, nil
	}

	q = bson.M{}
	for _, cond := range f.GetConditions() {
		key := cond.GetKey()
		op := cond.GetOp()
		value := cond.GetValue()
		switch op {
		case constants.FilterOpNotSet:
			// do nothing
		case constants.FilterOpEqual:
			q[key] = cond.GetValue()
		case constants.FilterOpNotEqual:
			q[key] = bson.M{"$ne": value}
		case constants.FilterOpContains, constants.FilterOpRegex, constants.FilterOpSearch:
			q[key] = bson.M{"$regex": value, "$options": "i"}
		case constants.FilterOpNotContains:
			q[key] = bson.M{"$not": bson.M{"$regex": value}}
		case constants.FilterOpIn:
			q[key] = bson.M{"$in": value}
		case constants.FilterOpNotIn:
			q[key] = bson.M{"$nin": value}
		case constants.FilterOpGreaterThan:
			q[key] = bson.M{"$gt": value}
		case constants.FilterOpGreaterThanEqual:
			q[key] = bson.M{"$gte": value}
		case constants.FilterOpLessThan:
			q[key] = bson.M{"$lt": value}
		case constants.FilterOpLessThanEqual:
			q[key] = bson.M{"$lte": value}
		default:
			return nil, errors.ErrorFilterInvalidOperation
		}
	}
	return q, nil
}
