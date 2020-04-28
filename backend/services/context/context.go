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
func (c *Context) Success(data interface{}, metas ...interface{}) {
	var meta interface{}
	if len(metas) == 0 {
		meta = gin.H{}
	} else {
		meta = metas[0]
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
	causeError := errors2.Cause(err)
	switch causeError.(type) {
	case errors.OPError:
		opError := causeError.(errors.OPError)

		c.AbortWithStatusJSON(opError.HttpCode, gin.H{
			"status":  "ok",
			"message": "error",
			"error":   errStr,
		})

	case validator.ValidationErrors:
		validatorErrors := causeError.(validator.ValidationErrors)
		//firstError := validatorErrors[0].(validator.FieldError)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "ok",
			"message": "error",
			"error":   validatorErrors.Error(),
		})
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
