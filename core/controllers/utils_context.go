package controllers

import (
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/models/models/v2"
	"github.com/gin-gonic/gin"
)

func GetUserFromContext(c *gin.Context) (u *models.UserV2) {
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
