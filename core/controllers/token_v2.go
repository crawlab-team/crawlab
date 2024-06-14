package controllers

import (
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/core/user"
	"github.com/gin-gonic/gin"
)

func PostToken(c *gin.Context) {
	var t models.TokenV2
	if err := c.ShouldBindJSON(&t); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	svc, err := user.GetUserServiceV2()
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	u := GetUserFromContextV2(c)
	t.SetCreated(u.Id)
	t.SetUpdated(u.Id)
	t.Token, err = svc.MakeToken(u)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	_, err = service.NewModelServiceV2[models.TokenV2]().InsertOne(t)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccess(c)
}
