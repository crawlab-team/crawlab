package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

func GetVersion(c *gin.Context) {
	version := viper.GetString("version")

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    version,
	})
}
