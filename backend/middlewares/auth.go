package middlewares

import (
	"crawlab/constants"
	"crawlab/services"
	"crawlab/services/context"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type AuthorizationComponent func(ctx *context.Context)

func CheckAccountEnabled(ctx *context.Context) {
	user := ctx.User()
	if !user.Enable && user.Username != "admin" {
		ctx.Failed(constants.ErrorAccountDisabled)
		return
	}
}
func CheckNeedResetPassword(ctx *context.Context) {
	user := ctx.User()

	if user.RePasswordTs.Before(time.Now()) {
		ctx.Failed(constants.ErrorNeedResetPassword)
		return
	}
}
func CheckUserPermission(ctx *context.Context) {
	user := ctx.User()
	// 如果为普通权限，校验请求地址是否符合要求
	if user.Role == constants.RoleNormal {
		if strings.HasPrefix(strings.ToLower(ctx.Request.URL.Path), "/users") {
			ctx.Failed(constants.ErrorAccountNoPermission)
			return
		}
	}
}
func TryAttachCurrentUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取token string
		tokenStr := ctx.GetHeader("Authorization")
		if len(tokenStr) < 20 {
			return
		}
		// 校验token
		user, err := services.CheckToken(tokenStr)
		if err == nil {
			ctx.Set(constants.ContextUser, &user)
		}
		ctx.Next()
	}
}

func AuthorizationMiddleware(components ...AuthorizationComponent) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithGinContext(c)

		if ctx.User() == nil {

			ctx.Failed(constants.ErrorTokenExpired)
			return
		}

		for _, component := range components {
			if c.IsAborted() {
				return
			}
			component(ctx)
		}
		if c.IsAborted() {
			return
		}
		// 校验成功
		ctx.Next()
	}
}
