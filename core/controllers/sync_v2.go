package controllers

import (
	"github.com/crawlab-team/crawlab/core/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"path/filepath"
)

func GetSyncScan(c *gin.Context) {
	id := c.Param("id")
	path := c.Query("path")

	workspacePath := viper.GetString("workspace")
	dirPath := filepath.Join(workspacePath, id, path)
	files, err := utils.ScanDirectory(dirPath)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, files)
}

func GetSyncDownload(c *gin.Context) {
	id := c.Param("id")
	path := c.Query("path")
	workspacePath := viper.GetString("workspace")
	filePath := filepath.Join(workspacePath, id, path)
	c.File(filePath)
}
