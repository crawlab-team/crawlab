package controllers

import (
	"github.com/crawlab-team/crawlab/core/container"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/delegate"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/gin-gonic/gin"
)

var TokenController *tokenController

var TokenActions []Action

type tokenController struct {
	ListActionControllerDelegate
	d   ListActionControllerDelegate
	ctx *tokenContext
}

func (ctr *tokenController) Post(c *gin.Context) {
	var err error
	var t models.Token
	if err := c.ShouldBindJSON(&t); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	u, err := ctr.ctx.userSvc.GetCurrentUser(c)
	if err != nil {
		HandleErrorUnauthorized(c, err)
		return
	}
	t.Token, err = ctr.ctx.userSvc.MakeToken(u)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	if err := delegate.NewModelDelegate(&t, GetUserFromContext(c)).Add(); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccess(c)
}

type tokenContext struct {
	modelSvc service.ModelService
	userSvc  interfaces.UserService
}

func newTokenContext() *tokenContext {
	// context
	ctx := &tokenContext{}

	// dependency injection
	if err := container.GetContainer().Invoke(func(
		modelSvc service.ModelService,
		userSvc interfaces.UserService,
	) {
		ctx.modelSvc = modelSvc
		ctx.userSvc = userSvc
	}); err != nil {
		panic(err)
	}

	return ctx
}

func newTokenController() *tokenController {
	modelSvc, err := service.GetService()
	if err != nil {
		panic(err)
	}

	ctr := NewListPostActionControllerDelegate(ControllerIdToken, modelSvc.GetBaseService(interfaces.ModelIdToken), TokenActions)
	d := NewListPostActionControllerDelegate(ControllerIdToken, modelSvc.GetBaseService(interfaces.ModelIdToken), TokenActions)
	ctx := newTokenContext()

	return &tokenController{
		ListActionControllerDelegate: *ctr,
		d:                            *d,
		ctx:                          ctx,
	}
}
