package controllers

import (
	"github.com/crawlab-team/crawlab/core/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"path/filepath"
)

var SyncController ActionController

func getSyncActions() []Action {
	var ctx = newSyncContext()
	return []Action{
		{
			Method:      http.MethodGet,
			Path:        "/:id/scan",
			HandlerFunc: ctx.scan,
		},
		{
			Method:      http.MethodGet,
			Path:        "/:id/download",
			HandlerFunc: ctx.download,
		},
	}
}

type syncContext struct {
}

func (ctx *syncContext) scan(c *gin.Context) {
	id := c.Param("id")
	dir := ctx._getDir(id)
	files, err := utils.ScanDirectory(dir)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, files)
}

func (ctx *syncContext) download(c *gin.Context) {
	id := c.Param("id")
	filePath := c.Query("path")
	dir := ctx._getDir(id)
	c.File(filepath.Join(dir, filePath))
}

func (ctx *syncContext) _getDir(id string) string {
	workspacePath := viper.GetString("workspace")
	return filepath.Join(workspacePath, id)
}

func newSyncContext() syncContext {
	return syncContext{}
}
