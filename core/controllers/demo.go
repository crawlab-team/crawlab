package controllers

import (
	"github.com/crawlab-team/crawlab/core/utils"
	"github.com/crawlab-team/crawlab/trace"
	"github.com/gin-gonic/gin"
	"net/http"
)

func getDemoActions() []Action {
	ctx := newDemoContext()
	return []Action{
		{
			Method:      http.MethodGet,
			Path:        "/import",
			HandlerFunc: ctx.import_,
		},
		{
			Method:      http.MethodGet,
			Path:        "/reimport",
			HandlerFunc: ctx.reimport,
		},
		{
			Method:      http.MethodGet,
			Path:        "/cleanup",
			HandlerFunc: ctx.cleanup,
		},
	}
}

type demoContext struct {
}

func (ctx *demoContext) import_(c *gin.Context) {
	if err := utils.ImportDemo(); err != nil {
		trace.PrintError(err)
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccess(c)
}

func (ctx *demoContext) reimport(c *gin.Context) {
	if err := utils.ReimportDemo(); err != nil {
		trace.PrintError(err)
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccess(c)
}

func (ctx *demoContext) cleanup(c *gin.Context) {
	if err := utils.ReimportDemo(); err != nil {
		trace.PrintError(err)
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccess(c)
}

var _demoCtx *demoContext

func newDemoContext() *demoContext {
	if _demoCtx != nil {
		return _demoCtx
	}

	_demoCtx = &demoContext{}

	return _demoCtx
}

var DemoController ActionController
