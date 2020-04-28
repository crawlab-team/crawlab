package routes

import (
	"crawlab/services"
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
)

func GetLatestRelease(c *gin.Context) {
	latestRelease, err := services.GetLatestRelease()
	if err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
	}
	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    latestRelease,
	})
}
