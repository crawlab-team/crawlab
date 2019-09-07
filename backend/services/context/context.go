package context

import (
	"crawlab/constants"
	"crawlab/errors"
	"crawlab/model"
	"fmt"
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	errors2 "github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"runtime/debug"
)

type Context struct {
	*gin.Context
}

func (c *Context) User() *model.User {
	userIfe, exists := c.Get(constants.ContextUser)
	if !exists {
		return nil
	}
	user, ok := userIfe.(*model.User)
	if !ok {
		return nil
	}
	return user
}
func (c *Context) Success(data interface{}, metas ...gin.H) {
	var meta gin.H
	if len(metas) > 0 {
		meta = metas[0]
	}
	if meta == nil {
		meta = gin.H{}
	}
	if data == nil {
		data = gin.H{}
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "success",
		"data":    data,
		"meta":    meta,
		"error":   "",
	})
}
func (c *Context) Failed(err error, variables ...interface{}) {
	c.failed(err, http.StatusOK, variables...)
}
func (c *Context) failed(err error, httpCode int, variables ...interface{}) {
	errStr := err.Error()
	if len(variables) > 0 {
		errStr = fmt.Sprintf(errStr, variables...)
	}
	log.Errorf("handle error:" + errStr)
	debug.PrintStack()
	switch innerError := errors2.Cause(err).(type) {
	case *errors.OPError:

		c.AbortWithStatusJSON(innerError.HttpCode, gin.H{
			"status":  "ok",
			"message": "error",
			"error":   errStr,
			"code":    innerError.Code,
		})
		break
	case *validator.ValidationErrors:

		c.Failed(constants.ErrorBadRequest)
		break
	default:
		fmt.Println("deprecated....")
		c.AbortWithStatusJSON(httpCode, gin.H{
			"status":  "ok",
			"message": "error",
			"error":   errStr,
		})
	}
}
func (c *Context) FailedWithError(err error, httpCode ...int) {

	var code = 200
	if len(httpCode) > 0 {
		code = httpCode[0]
	}
	c.failed(err, code)

}

func WithGinContext(context *gin.Context) *Context {
	return &Context{Context: context}
}
