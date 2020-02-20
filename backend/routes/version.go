package routes

import (
	"crawlab/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetLatestRelease(c *gin.Context) {
	latestRelease, err := services.GetLatestRelease()
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    latestRelease,
	})
}
