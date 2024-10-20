package middlewares

import (
	"errors"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/models/models/v2"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/core/user"
	"github.com/crawlab-team/crawlab/core/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
)

func AuthorizationMiddlewareV2() gin.HandlerFunc {
	userSvc, _ := user.GetUserServiceV2()
	return func(c *gin.Context) {
		// disable auth for test
		if viper.GetBool("auth.disabled") {
			u, err := service.NewModelServiceV2[models.UserV2]().GetOne(bson.M{"username": constants.DefaultAdminUsername}, nil)
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
			utils.HandleErrorUnauthorized(c, errors.New("invalid token"))
			return
		}

		// set user in context
		c.Set(constants.UserContextKey, u)

		// validation success
		c.Next()
	}
}
