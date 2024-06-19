package controllers

import (
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/entity"
	"github.com/crawlab-team/crawlab/trace"
	"github.com/gin-gonic/gin"
	"net/http"
)

func handleError(statusCode int, c *gin.Context, err error, print bool) {
	if print {
		trace.PrintError(err)
	}
	c.AbortWithStatusJSON(statusCode, entity.Response{
		Status:  constants.HttpResponseStatusOk,
		Message: constants.HttpResponseMessageError,
		Error:   err.Error(),
	})
}

func HandleError(statusCode int, c *gin.Context, err error) {
	handleError(statusCode, c, err, true)
}

func HandleErrorNoPrint(statusCode int, c *gin.Context, err error) {
	handleError(statusCode, c, err, false)
}

func HandleErrorBadRequest(c *gin.Context, err error) {
	HandleError(http.StatusBadRequest, c, err)
}

func HandleErrorForbidden(c *gin.Context, err error) {
	HandleError(http.StatusForbidden, c, err)
}

func HandleErrorUnauthorized(c *gin.Context, err error) {
	HandleError(http.StatusUnauthorized, c, err)
}

func HandleErrorNotFound(c *gin.Context, err error) {
	HandleError(http.StatusNotFound, c, err)
}

func HandleErrorNotFoundNoPrint(c *gin.Context, err error) {
	HandleErrorNoPrint(http.StatusNotFound, c, err)
}

func HandleErrorInternalServerError(c *gin.Context, err error) {
	HandleError(http.StatusInternalServerError, c, err)
}

func HandleSuccess(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusOK, entity.Response{
		Status:  constants.HttpResponseStatusOk,
		Message: constants.HttpResponseMessageSuccess,
	})
}

func HandleSuccessWithData(c *gin.Context, data interface{}) {
	c.AbortWithStatusJSON(http.StatusOK, entity.Response{
		Status:  constants.HttpResponseStatusOk,
		Message: constants.HttpResponseMessageSuccess,
		Data:    data,
	})
}

func HandleSuccessWithListData(c *gin.Context, data interface{}, total int) {
	c.AbortWithStatusJSON(http.StatusOK, entity.ListResponse{
		Status:  constants.HttpResponseStatusOk,
		Message: constants.HttpResponseMessageSuccess,
		Data:    data,
		Total:   total,
	})
}
