package controllers

import (
	"github.com/crawlab-team/crawlab/core/entity"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func GetSystemInfo(c *gin.Context) {
	info := &entity.SystemInfo{
		Edition: viper.GetString("edition"),
		Version: viper.GetString("version"),
	}
	HandleSuccessWithData(c, info)
}
