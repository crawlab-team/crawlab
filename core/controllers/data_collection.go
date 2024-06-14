package controllers

import (
	"github.com/crawlab-team/crawlab/core/container"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/db/mongo"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

var DataCollectionController *dataCollectionController

func getDataCollectionActions() []Action {
	ctx := newDataCollectionContext()
	return []Action{
		{
			Method:      http.MethodPost,
			Path:        "/:id/indexes",
			HandlerFunc: ctx.postIndexes,
		},
	}
}

type dataCollectionController struct {
	ListActionControllerDelegate
	d   ListActionControllerDelegate
	ctx *dataCollectionContext
}

type dataCollectionContext struct {
	modelSvc  service.ModelService
	resultSvc interfaces.ResultService
}

func (ctx *dataCollectionContext) postIndexes(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	dc, err := ctx.modelSvc.GetDataCollectionById(id)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	for _, f := range dc.Fields {
		if err := mongo.GetMongoCol(dc.Name).CreateIndex(mongo2.IndexModel{
			Keys: f.Key,
		}); err != nil {
			HandleErrorInternalServerError(c, err)
			return
		}
	}

	HandleSuccess(c)
}

var _dataCollectionCtx *dataCollectionContext

func newDataCollectionContext() *dataCollectionContext {
	if _dataCollectionCtx != nil {
		return _dataCollectionCtx
	}

	// context
	ctx := &dataCollectionContext{}

	// dependency injection
	if err := container.GetContainer().Invoke(func(
		modelSvc service.ModelService,
	) {
		ctx.modelSvc = modelSvc
	}); err != nil {
		panic(err)
	}

	_dataCollectionCtx = ctx

	return ctx
}

func newDataCollectionController() *dataCollectionController {
	actions := getDataCollectionActions()
	modelSvc, err := service.GetService()
	if err != nil {
		panic(err)
	}

	ctr := NewListPostActionControllerDelegate(ControllerIdDataCollection, modelSvc.GetBaseService(interfaces.ModelIdDataCollection), actions)
	d := NewListPostActionControllerDelegate(ControllerIdDataCollection, modelSvc.GetBaseService(interfaces.ModelIdDataCollection), actions)
	ctx := newDataCollectionContext()

	return &dataCollectionController{
		ListActionControllerDelegate: *ctr,
		d:                            *d,
		ctx:                          ctx,
	}
}
