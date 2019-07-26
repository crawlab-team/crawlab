package middlewares

import (
	"crawlab/routes"
	"crawlab/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthorizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if c.Request.URL.Path == "/login" || (c.Request.URL.Path == "/users" && c.Request.Method == "PUT") {
			c.Next()
		} else {
			_, err := services.CheckToken(tokenStr)
			if err == nil {
				c.Next()
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, routes.Response{
					Status:  "ok",
					Message: "unauthorized",
					Error:   "unauthorized",
				})
			}
		}
	}
}
