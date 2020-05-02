package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

type SettingBody struct {
	AllowRegister     string `json:"allow_register"`
	EnableTutorial    string `json:"enable_tutorial"`
	RunOnMaster       string `json:"run_on_master"`
	EnableDemoSpiders string `json:"enable_demo_spiders"`
}

// @Summary Get version
// @Description Get version
// @Tags setting
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /version [get]
func GetVersion(c *gin.Context) {
	version := viper.GetString("version")

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    version,
	})
}

// @Summary Get setting
// @Description Get setting
// @Tags setting
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /setting [get]
func GetSetting(c *gin.Context) {
	body := SettingBody{
		AllowRegister:     viper.GetString("setting.allowRegister"),
		EnableTutorial:    viper.GetString("setting.enableTutorial"),
		RunOnMaster:       viper.GetString("setting.runOnMaster"),
		EnableDemoSpiders: viper.GetString("setting.enableDemoSpiders"),
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    body,
	})
}
