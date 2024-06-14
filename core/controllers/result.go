package controllers

import (
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/core/result"
	"github.com/crawlab-team/crawlab/core/utils"
	"github.com/crawlab-team/crawlab/db/generic"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

var ResultController ActionController

func getResultActions() []Action {
	var resultCtx = newResultContext()
	return []Action{
		{
			Method:      http.MethodGet,
			Path:        "/:id",
			HandlerFunc: resultCtx.getList,
		},
	}
}

type resultContext struct {
	modelSvc service.ModelService
}

func (ctx *resultContext) getList(c *gin.Context) {
	// data collection id
	dcId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// data source id
	var dsId primitive.ObjectID
	dsIdStr := c.Query("data_source_id")
	if dsIdStr != "" {
		dsId, err = primitive.ObjectIDFromHex(dsIdStr)
		if err != nil {
			HandleErrorBadRequest(c, err)
			return
		}
	}

	// data collection
	dc, err := ctx.modelSvc.GetDataCollectionById(dcId)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// data source
	ds, err := ctx.modelSvc.GetDataSourceById(dsId)
	if err != nil {
		if err.Error() == mongo2.ErrNoDocuments.Error() {
			ds = &models.DataSource{}
		} else {
			HandleErrorInternalServerError(c, err)
			return
		}
	}

	// spider
	sq := bson.M{
		"col_id":         dc.Id,
		"data_source_id": ds.Id,
	}
	s, err := ctx.modelSvc.GetSpider(sq, nil)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// service
	svc, err := result.GetResultService(s.Id)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// params
	pagination := MustGetPagination(c)
	query := ctx.getListQuery(c)

	// get results
	data, err := svc.List(query, &generic.ListOptions{
		Sort:  []generic.ListSort{{"_id", generic.SortDirectionDesc}},
		Skip:  pagination.Size * (pagination.Page - 1),
		Limit: pagination.Size,
	})
	if err != nil {
		if err.Error() == mongo2.ErrNoDocuments.Error() {
			HandleSuccessWithListData(c, nil, 0)
			return
		}
		HandleErrorInternalServerError(c, err)
		return
	}

	// validate results
	if len(data) == 0 {
		HandleSuccessWithListData(c, nil, 0)
		return
	}

	// total count
	total, err := svc.Count(query)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// response
	HandleSuccessWithListData(c, data, total)
}

func (ctx *resultContext) getListQuery(c *gin.Context) (q generic.ListQuery) {
	f, err := GetFilter(c)
	if err != nil {
		return q
	}
	for _, cond := range f.Conditions {
		q = append(q, generic.ListQueryCondition{
			Key:   cond.Key,
			Op:    cond.Op,
			Value: utils.NormalizeObjectId(cond.Value),
		})
	}
	return q
}

func newResultContext() *resultContext {
	// context
	ctx := &resultContext{}

	var err error
	ctx.modelSvc, err = service.NewService()
	if err != nil {
		panic(err)
	}

	return ctx
}
