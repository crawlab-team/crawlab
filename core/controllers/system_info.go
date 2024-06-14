package controllers

import (
	"github.com/crawlab-team/crawlab/core/entity"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

func getSystemInfo(c *gin.Context) {
	info := &entity.SystemInfo{
		Edition: viper.GetString("info.edition"),
		Version: viper.GetString("info.version"),
	}
	HandleSuccessWithData(c, info)
}

func getSystemInfoActions() []Action {
	return []Action{
		{
			Path:        "",
			Method:      http.MethodGet,
			HandlerFunc: getSystemInfo,
		},
	}
}

var SystemInfoController ActionController
