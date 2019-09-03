package routes

import (
	"github.com/gin-gonic/gin"
	"runtime/debug"
)

func HandleError(statusCode int, c *gin.Context, err error) {
	debug.PrintStack()
	c.JSON(statusCode, Response{
		Status:  "ok",
		Message: "error",
		Error:   err.Error(),
	})
}

func HandleErrorF(statusCode int, c *gin.Context, err string) {
	debug.PrintStack()
	c.JSON(statusCode, Response{
		Status:  "ok",
		Message: "error",
		Error:   err,
	})
}
