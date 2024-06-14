package controllers

import (
	"errors"
	"fmt"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/export"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/gin-gonic/gin"
	"net/http"
)

var ExportController ActionController

func getExportActions() []Action {
	ctx := newExportContext()
	return []Action{
		{
			Method:      http.MethodPost,
			Path:        "/:type",
			HandlerFunc: ctx.postExport,
		},
		{
			Method:      http.MethodGet,
			Path:        "/:type/:id",
			HandlerFunc: ctx.getExport,
		},
		{
			Method:      http.MethodGet,
			Path:        "/:type/:id/download",
			HandlerFunc: ctx.getExportDownload,
		},
	}
}

type exportContext struct {
	csvSvc  interfaces.ExportService
	jsonSvc interfaces.ExportService
}

func (ctx *exportContext) postExport(c *gin.Context) {
	exportType := c.Param("type")
	exportTarget := c.Query("target")
	exportFilter, _ := GetFilter(c)

	var exportId string
	var err error
	switch exportType {
	case constants.ExportTypeCsv:
		exportId, err = ctx.csvSvc.Export(exportType, exportTarget, exportFilter)
	case constants.ExportTypeJson:
		exportId, err = ctx.jsonSvc.Export(exportType, exportTarget, exportFilter)
	default:
		HandleErrorBadRequest(c, errors.New(fmt.Sprintf("invalid export type: %s", exportType)))
		return
	}
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccessWithData(c, exportId)
}

func (ctx *exportContext) getExport(c *gin.Context) {
	exportType := c.Param("type")
	exportId := c.Param("id")

	var exp interfaces.Export
	var err error
	switch exportType {
	case constants.ExportTypeCsv:
		exp, err = ctx.csvSvc.GetExport(exportId)
	case constants.ExportTypeJson:
		exp, err = ctx.jsonSvc.GetExport(exportId)
	default:
		HandleErrorBadRequest(c, errors.New(fmt.Sprintf("invalid export type: %s", exportType)))
	}
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccessWithData(c, exp)
}

func (ctx *exportContext) getExportDownload(c *gin.Context) {
	exportType := c.Param("type")
	exportId := c.Param("id")

	var exp interfaces.Export
	var err error
	switch exportType {
	case constants.ExportTypeCsv:
		exp, err = ctx.csvSvc.GetExport(exportId)
	case constants.ExportTypeJson:
		exp, err = ctx.jsonSvc.GetExport(exportId)
	default:
		HandleErrorBadRequest(c, errors.New(fmt.Sprintf("invalid export type: %s", exportType)))
	}
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	switch exportType {
	case constants.ExportTypeCsv:
		c.Header("Content-Type", "text/csv")
	case constants.ExportTypeJson:
		c.Header("Content-Type", "text/plain")
	default:
		HandleErrorBadRequest(c, errors.New(fmt.Sprintf("invalid export type: %s", exportType)))
	}
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", exp.GetDownloadPath()))
	c.File(exp.GetDownloadPath())
}

func newExportContext() *exportContext {
	return &exportContext{
		csvSvc:  export.GetCsvService(),
		jsonSvc: export.GetJsonService(),
	}
}
