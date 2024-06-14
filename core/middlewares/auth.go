package middlewares

import (
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/core/user"
	"github.com/crawlab-team/crawlab/core/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func AuthorizationMiddleware() gin.HandlerFunc {
	userSvc, _ := user.GetUserService()
	return func(c *gin.Context) {
		// disable auth for test
		if viper.GetBool("auth.disabled") {
			modelSvc, err := service.GetService()
			if err != nil {
				utils.HandleErrorInternalServerError(c, err)
				return
			}
			u, err := modelSvc.GetUserByUsername(constants.DefaultAdminUsername, nil)
			if err != nil {
				utils.HandleErrorInternalServerError(c, err)
				return
			}
			c.Set(constants.UserContextKey, u)
			c.Next()
			return
		}

		// token string
		tokenStr := c.GetHeader("Authorization")

		// validate token
		u, err := userSvc.CheckToken(tokenStr)
		if err != nil {
			// validation failed, return error response
			utils.HandleErrorUnauthorized(c, errors.ErrorHttpUnauthorized)
			return
		}

		// set user in context
		c.Set(constants.UserContextKey, u)

		// validation success
		c.Next()
	}
}
