package controllers

import (
	"github.com/crawlab-team/crawlab/core/entity"
	"github.com/crawlab-team/crawlab/db/mongo"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

var FilterController ActionController

func getFilterActions() []Action {
	ctx := newFilterContext()
	return []Action{
		{
			Method:      http.MethodGet,
			Path:        "/:col",
			HandlerFunc: ctx.getColFieldOptions,
		},
		{
			Method:      http.MethodGet,
			Path:        "/:col/:value",
			HandlerFunc: ctx.getColFieldOptions,
		},
		{
			Method:      http.MethodGet,
			Path:        "/:col/:value/:label",
			HandlerFunc: ctx.getColFieldOptions,
		},
	}
}

type filterContext struct {
}

func (ctx *filterContext) getColFieldOptions(c *gin.Context) {
	colName := c.Param("col")
	value := c.Param("value")
	if value == "" {
		value = "_id"
	}
	label := c.Param("label")
	if label == "" {
		label = "name"
	}
	query := MustGetFilterQuery(c)
	pipelines := mongo2.Pipeline{}
	if query != nil {
		pipelines = append(pipelines, bson.D{{"$match", query}})
	}
	pipelines = append(
		pipelines,
		bson.D{
			{
				"$group",
				bson.M{
					"_id": bson.M{
						"value": "$" + value,
						"label": "$" + label,
					},
				},
			},
		},
	)
	pipelines = append(
		pipelines,
		bson.D{
			{
				"$project",
				bson.M{
					"value": "$_id.value",
					"label": bson.M{"$toString": "$_id.label"},
				},
			},
		},
	)
	pipelines = append(
		pipelines,
		bson.D{
			{
				"$sort",
				bson.D{
					{"value", 1},
				},
			},
		},
	)
	var options []entity.FilterSelectOption
	if err := mongo.GetMongoCol(colName).Aggregate(pipelines, nil).All(&options); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccessWithData(c, options)
}

func newFilterContext() *filterContext {
	return &filterContext{}
}
