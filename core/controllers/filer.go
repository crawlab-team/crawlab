package controllers

import (
	"bufio"
	"fmt"
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
	"github.com/spf13/viper"
	"net/http"
	"strings"
)

var FilerController ActionController

func getFilerActions() []Action {
	filerCtx := newFilerContext()
	return []Action{
		{
			Method:      http.MethodGet,
			Path:        "*path",
			HandlerFunc: filerCtx.do,
		},
		{
			Method:      http.MethodPost,
			Path:        "*path",
			HandlerFunc: filerCtx.do,
		},
		{
			Method:      http.MethodPut,
			Path:        "*path",
			HandlerFunc: filerCtx.do,
		},
		{
			Method:      http.MethodDelete,
			Path:        "*path",
			HandlerFunc: filerCtx.do,
		},
	}
}

type filerContext struct {
	endpoint string
}

func (ctx *filerContext) do(c *gin.Context) {
	// request path
	requestPath := strings.Replace(c.Request.URL.Path, "/filer", "", 1)

	// request url
	requestUrl := fmt.Sprintf("%s%s", ctx.endpoint, requestPath)
	if c.Request.URL.RawQuery != "" {
		requestUrl += "?" + c.Request.URL.RawQuery
	}

	// request body
	bufR := bufio.NewScanner(c.Request.Body)
	requestBody := req.BodyJSON(bufR.Bytes())

	// request file uploads
	var requestFileUploads []req.FileUpload
	form, err := c.MultipartForm()
	if err == nil {
		for k, v := range form.File {
			for _, fh := range v {
				f, err := fh.Open()
				if err != nil {
					HandleErrorInternalServerError(c, err)
					return
				}
				requestFileUploads = append(requestFileUploads, req.FileUpload{
					FileName:  fh.Filename,
					FieldName: k,
					File:      f,
				})
			}
		}
	}

	// request header
	requestHeader := req.Header{}
	for k, v := range c.Request.Header {
		if len(v) > 0 {
			requestHeader[k] = v[0]
		}
	}

	// perform request
	res, err := req.Do(c.Request.Method, requestUrl, requestHeader, requestBody, requestFileUploads)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// status code check
	statusCode := res.Response().StatusCode
	if statusCode == http.StatusNotFound {
		HandleErrorNotFoundNoPrint(c, errors.ErrorControllerFilerNotFound)
		return
	}

	// response
	for k, v := range res.Response().Header {
		if len(v) > 0 {
			c.Header(k, v[0])
		}
	}
	_, _ = c.Writer.Write(res.Bytes())
	c.AbortWithStatus(statusCode)
}

var _filerCtx *filerContext

func newFilerContext() *filerContext {
	if _filerCtx != nil {
		return _filerCtx
	}

	ctx := &filerContext{
		endpoint: "http://localhost:8888",
	}

	if viper.GetString("fs.filer.proxy") != "" {
		ctx.endpoint = viper.GetString("fs.filer.proxy")
	}

	_filerCtx = ctx

	return ctx
}
