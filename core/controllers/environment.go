package controllers

import (
	"github.com/crawlab-team/crawlab/core/container"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/service"
)

var EnvironmentController *environmentController

var EnvironmentActions []Action

type environmentController struct {
	ListActionControllerDelegate
	d   ListActionControllerDelegate
	ctx *environmentContext
}

type environmentContext struct {
	modelSvc service.ModelService
	userSvc  interfaces.UserService
}

func newEnvironmentContext() *environmentContext {
	// context
	ctx := &environmentContext{}

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

func newEnvironmentController() *environmentController {
	modelSvc, err := service.GetService()
	if err != nil {
		panic(err)
	}

	ctr := NewListPostActionControllerDelegate(ControllerIdEnvironment, modelSvc.GetBaseService(interfaces.ModelIdEnvironment), EnvironmentActions)
	d := NewListPostActionControllerDelegate(ControllerIdEnvironment, modelSvc.GetBaseService(interfaces.ModelIdEnvironment), EnvironmentActions)
	ctx := newEnvironmentContext()

	return &environmentController{
		ListActionControllerDelegate: *ctr,
		d:                            *d,
		ctx:                          ctx,
	}
}
