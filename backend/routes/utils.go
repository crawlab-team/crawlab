package routes

import (
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"runtime/debug"
)

func HandleError(statusCode int, c *gin.Context, err error) {
	log.Errorf("handle error:" + err.Error())
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
