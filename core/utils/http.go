package utils

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

func HandleErrorUnauthorized(c *gin.Context, err error) {
	HandleError(http.StatusUnauthorized, c, err)
}

func HandleErrorInternalServerError(c *gin.Context, err error) {
	HandleError(http.StatusInternalServerError, c, err)
}
