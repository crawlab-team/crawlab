package controllers

import (
	"github.com/crawlab-team/crawlab/core/errors"
	models2 "github.com/crawlab-team/crawlab/core/models/models/v2"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/db/mongo"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
)

func GetProjectList(c *gin.Context) {
	// get all list
	all := MustGetFilterAll(c)
	if all {
		NewControllerV2[models2.ProjectV2]().getAll(c)
		return
	}

	// params
	pagination := MustGetPagination(c)
	query := MustGetFilterQuery(c)
	sort := MustGetSortOption(c)

	// get list
	projects, err := service.NewModelServiceV2[models2.ProjectV2]().GetMany(query, &mongo.FindOptions{
		Sort:  sort,
		Skip:  pagination.Size * (pagination.Page - 1),
		Limit: pagination.Size,
	})
	if err != nil {
		if err.Error() != mongo2.ErrNoDocuments.Error() {
			HandleErrorInternalServerError(c, err)
		}
		return
	}
	if len(projects) == 0 {
		HandleSuccessWithListData(c, []models2.ProjectV2{}, 0)
		return
	}

	// total count
	total, err := service.NewModelServiceV2[models2.ProjectV2]().Count(query)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// project ids
	var ids []primitive.ObjectID

	// count cache
	cache := map[primitive.ObjectID]int{}

	// iterate
	for _, p := range projects {
		ids = append(ids, p.Id)
		cache[p.Id] = 0
	}

	// spiders
	spiders, err := service.NewModelServiceV2[models2.SpiderV2]().GetMany(bson.M{
		"project_id": bson.M{
			"$in": ids,
		},
	}, nil)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	for _, s := range spiders {
		_, ok := cache[s.ProjectId]
		if !ok {
			HandleErrorInternalServerError(c, errors.ErrorControllerMissingInCache)
			return
		}
		cache[s.ProjectId]++
	}

	// assign
	var data []models2.ProjectV2
	for _, p := range projects {
		p.Spiders = cache[p.Id]
		data = append(data, p)
	}

	HandleSuccessWithListData(c, data, total)
}
