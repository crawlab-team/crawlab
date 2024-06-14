package controllers

import (
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/user"
	"github.com/gin-gonic/gin"
)

func PostLogin(c *gin.Context) {
	var payload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	userSvc, err := user.GetUserServiceV2()
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	token, loggedInUser, err := userSvc.Login(payload.Username, payload.Password)
	if err != nil {
		HandleErrorUnauthorized(c, errors.ErrorUserUnauthorized)
		return
	}
	c.Set(constants.UserContextKey, loggedInUser)
	HandleSuccessWithData(c, token)
}

func PostLogout(c *gin.Context) {
	c.Set(constants.UserContextKey, nil)
	HandleSuccess(c)
}
