package controllers

import (
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/gin-gonic/gin"
)

func GetUserFromContext(c *gin.Context) (u interfaces.User) {
	value, ok := c.Get(constants.UserContextKey)
	if !ok {
		return nil
	}
	u, ok = value.(interfaces.User)
	if !ok {
		return nil
	}
	return u
}

func GetUserFromContextV2(c *gin.Context) (u *models.UserV2) {
	value, ok := c.Get(constants.UserContextKey)
	if !ok {
		return nil
	}
	u, ok = value.(*models.UserV2)
	if !ok {
		return nil
	}
	return u
}
