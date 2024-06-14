package controllers

import (
	"errors"
	"fmt"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/export"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/gin-gonic/gin"
)

func PostExport(c *gin.Context) {
	exportType := c.Param("type")
	exportTarget := c.Query("target")
	exportFilter, _ := GetFilter(c)

	var exportId string
	var err error
	switch exportType {
	case constants.ExportTypeCsv:
		exportId, err = export.GetCsvService().Export(exportType, exportTarget, exportFilter)
	case constants.ExportTypeJson:
		exportId, err = export.GetJsonService().Export(exportType, exportTarget, exportFilter)
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

func GetExport(c *gin.Context) {
	exportType := c.Param("type")
	exportId := c.Param("id")

	var exp interfaces.Export
	var err error
	switch exportType {
	case constants.ExportTypeCsv:
		exp, err = export.GetCsvService().GetExport(exportId)
	case constants.ExportTypeJson:
		exp, err = export.GetJsonService().GetExport(exportId)
	default:
		HandleErrorBadRequest(c, errors.New(fmt.Sprintf("invalid export type: %s", exportType)))
	}
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccessWithData(c, exp)
}

func GetExportDownload(c *gin.Context) {
	exportType := c.Param("type")
	exportId := c.Param("id")

	var exp interfaces.Export
	var err error
	switch exportType {
	case constants.ExportTypeCsv:
		exp, err = export.GetCsvService().GetExport(exportId)
	case constants.ExportTypeJson:
		exp, err = export.GetJsonService().GetExport(exportId)
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
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", exp.GetDownloadPath()))
	c.File(exp.GetDownloadPath())
}
