package controllers

import (
	"encoding/json"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/entity"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// GetSorts Get entity.Sort from gin.Context
func GetSorts(c *gin.Context) (sorts []entity.Sort, err error) {
	// bind
	sortStr := c.Query(constants.SortQueryField)
	if err := json.Unmarshal([]byte(sortStr), &sorts); err != nil {
		return nil, err
	}
	return sorts, nil
}

// GetSortsOption Get entity.Sort from gin.Context
func GetSortsOption(c *gin.Context) (sort bson.D, err error) {
	sorts, err := GetSorts(c)
	if err != nil {
		return nil, err
	}

	if sorts == nil || len(sorts) == 0 {
		return bson.D{{"_id", -1}}, nil
	}

	return SortsToOption(sorts)
}

func MustGetSortOption(c *gin.Context) (sort bson.D) {
	sort, err := GetSortsOption(c)
	if err != nil {
		return nil
	}
	return sort
}

// SortsToOption Translate entity.Sort to bson.D
func SortsToOption(sorts []entity.Sort) (sort bson.D, err error) {
	sort = bson.D{}
	for _, s := range sorts {
		switch s.Direction {
		case constants.ASCENDING:
			sort = append(sort, bson.E{Key: s.Key, Value: 1})
		case constants.DESCENDING:
			sort = append(sort, bson.E{Key: s.Key, Value: -1})
		}
	}
	if len(sort) == 0 {
		sort = bson.D{{"_id", -1}}
	}
	return sort, nil
}
