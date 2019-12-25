package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

type SettingBody struct {
	AllowRegister string `json:"allow_register"`
}

func GetVersion(c *gin.Context) {
	version := viper.GetString("version")

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    version,
	})
}

func GetSetting(c *gin.Context) {
	allowRegister := viper.GetString("setting.allowRegister")

	body := SettingBody{AllowRegister: allowRegister}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    body,
	})
}
