package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
)

func HandleError(statusCode int, c *gin.Context, err error) {
	c.AbortWithStatusJSON(statusCode, Response{
		Status:  "error",
		Message: "failure",
		Error:   err.Error(),
	})
}

func HandleErrorF(statusCode int, c *gin.Context, err string) {
	debug.PrintStack()
	c.AbortWithStatusJSON(statusCode, Response{
		Status:  "ok",
		Message: "error",
		Error:   err,
	})
}

func HandleSuccess(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

func HandleSuccessData(c *gin.Context, data interface{}) {
	c.AbortWithStatusJSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    data,
	})
}
