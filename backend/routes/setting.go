package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

type SettingBody struct {
	AllowRegister  string `json:"allow_register"`
	EnableTutorial string `json:"enable_tutorial"`
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
	body := SettingBody{
		AllowRegister:  viper.GetString("setting.allowRegister"),
		EnableTutorial: viper.GetString("setting.enableTutorial"),
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    body,
	})
}
