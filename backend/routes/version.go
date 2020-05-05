package routes

import (
	"crawlab/services"
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
)

// @Summary Get  latest release
// @Description Get latest release
// @Tags version
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /releases/latest [get]
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
