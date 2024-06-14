package controllers

import (
	"encoding/json"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/entity"
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"strings"
)

// GetFilter Get entity.Filter from gin.Context
func GetFilter(c *gin.Context) (f *entity.Filter, err error) {
	// bind
	condStr := c.Query(constants.FilterQueryFieldConditions)
	var conditions []*entity.Condition
	if err := json.Unmarshal([]byte(condStr), &conditions); err != nil {
		return nil, err
	}

	// attempt to convert object id
	for i, cond := range conditions {
		v := reflect.ValueOf(cond.Value)
		switch v.Kind() {
		case reflect.String:
			item := cond.Value.(string)
			id, err := primitive.ObjectIDFromHex(item)
			if err == nil {
				conditions[i].Value = id
			} else {
				conditions[i].Value = item
			}
		case reflect.Slice, reflect.Array:
			var items []interface{}
			for i := 0; i < v.Len(); i++ {
				vItem := v.Index(i)
				item := vItem.Interface()

				// string
				stringItem, ok := item.(string)
				if ok {
					id, err := primitive.ObjectIDFromHex(stringItem)
					if err == nil {
						items = append(items, id)
					} else {
						items = append(items, stringItem)
					}
					continue
				}

				// default
				items = append(items, item)
			}
			conditions[i].Value = items
		}
	}

	return &entity.Filter{
		IsOr:       false,
		Conditions: conditions,
	}, nil
}

// GetFilterQuery Get bson.M from gin.Context
func GetFilterQuery(c *gin.Context) (q bson.M, err error) {
	f, err := GetFilter(c)
	if err != nil {
		return nil, err
	}

	if f == nil {
		return nil, nil
	}

	// TODO: implement logic OR

	return utils.FilterToQuery(f)
}

func MustGetFilterQuery(c *gin.Context) (q bson.M) {
	q, err := GetFilterQuery(c)
	if err != nil {
		return nil
	}
	return q
}

// GetFilterAll Get all from gin.Context
func GetFilterAll(c *gin.Context) (res bool, err error) {
	resStr := c.Query(constants.FilterQueryFieldAll)
	switch strings.ToUpper(resStr) {
	case "1":
		return true, nil
	case "0":
		return false, nil
	case "Y":
		return true, nil
	case "N":
		return false, nil
	case "T":
		return true, nil
	case "F":
		return false, nil
	case "TRUE":
		return true, nil
	case "FALSE":
		return false, nil
	default:
		return false, errors.ErrorFilterInvalidOperation
	}
}

func MustGetFilterAll(c *gin.Context) (res bool) {
	res, err := GetFilterAll(c)
	if err != nil {
		return false
	}
	return res
}
