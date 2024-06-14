package controllers

import (
	"github.com/crawlab-team/crawlab/core/ds"
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	interfaces2 "github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/delegate"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/core/utils"
	"github.com/crawlab-team/crawlab/db/mongo"
	"github.com/crawlab-team/crawlab/trace"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

var DataSourceController *dataSourceController

func getDataSourceActions() []Action {
	ctx := newDataSourceContext()
	return []Action{
		{
			Path:        "/:id/change-password",
			Method:      http.MethodPost,
			HandlerFunc: ctx.changePassword,
		},
	}
}

type dataSourceController struct {
	ListActionControllerDelegate
	d   ListActionControllerDelegate
	ctx *dataSourceContext
}

func (ctr *dataSourceController) Post(c *gin.Context) {
	// data source
	var _ds models.DataSource
	if err := c.ShouldBindJSON(&_ds); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// add data source to db
	if err := mongo.RunTransaction(func(ctx mongo2.SessionContext) error {
		if err := delegate.NewModelDelegate(&_ds).Add(); err != nil {
			return trace.TraceError(err)
		}
		pwd, err := utils.EncryptAES(_ds.Password)
		if err != nil {
			return trace.TraceError(err)
		}
		p := models.Password{Id: _ds.Id, Password: pwd}
		if err := delegate.NewModelDelegate(&p).Add(); err != nil {
			return trace.TraceError(err)
		}
		return nil
	}); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// check data source status
	go func() { _ = ctr.ctx.dsSvc.CheckStatus(_ds.Id) }()

	HandleSuccess(c)
}

func (ctr *dataSourceController) Put(c *gin.Context) {
	// data source
	var _ds models.DataSource
	if err := c.ShouldBindJSON(&_ds); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	if err := delegate.NewModelDelegate(&_ds).Save(); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// check data source status
	go func() { _ = ctr.ctx.dsSvc.CheckStatus(_ds.Id) }()
}

type dataSourceContext struct {
	dsSvc interfaces.DataSourceService
}

var _dataSourceCtx *dataSourceContext

func newDataSourceContext() *dataSourceContext {
	if _dataSourceCtx != nil {
		return _dataSourceCtx
	}
	dsSvc, err := ds.GetDataSourceService()
	if err != nil {
		panic(err)
	}
	_dataSourceCtx = &dataSourceContext{
		dsSvc: dsSvc,
	}
	return _dataSourceCtx
}

func (ctx *dataSourceContext) changePassword(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	var payload map[string]string
	if err := c.ShouldBindJSON(&payload); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	password, ok := payload["password"]
	if !ok {
		HandleErrorBadRequest(c, errors.ErrorDataSourceMissingRequiredFields)
		return
	}
	if err := ctx.dsSvc.ChangePassword(id, password); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccess(c)
}

func newDataSourceController() *dataSourceController {
	actions := getDataSourceActions()
	modelSvc, err := service.GetService()
	if err != nil {
		panic(err)
	}

	ctr := NewListPostActionControllerDelegate(ControllerIdDataSource, modelSvc.GetBaseService(interfaces2.ModelIdDataSource), actions)
	d := NewListPostActionControllerDelegate(ControllerIdDataSource, modelSvc.GetBaseService(interfaces2.ModelIdDataSource), actions)
	ctx := newDataSourceContext()

	return &dataSourceController{
		ListActionControllerDelegate: *ctr,
		d:                            *d,
		ctx:                          ctx,
	}
}
