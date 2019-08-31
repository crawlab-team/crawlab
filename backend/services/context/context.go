package context

import (
	"crawlab/constants"
	"crawlab/errors"
	"crawlab/model"
	"fmt"
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	errors2 "github.com/pkg/errors"
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
func (c *Context) Success(data interface{}, meta interface{}) {
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
		"error":   "",
	})
}
func (c *Context) FailedWithError(err error, httpCode ...int) {

	var code = 200
	if len(httpCode) > 0 {
		code = httpCode[0]
	}
	log.Errorf("handle error:" + err.Error())
	debug.PrintStack()
	switch errors2.Cause(err).(type) {
	case errors.OPError:
		c.AbortWithStatusJSON(code, gin.H{
			"status":  "ok",
			"message": "error",
			"error":   err.Error(),
		})
		break
	default:
		fmt.Println("deprecated....")
		c.AbortWithStatusJSON(code, gin.H{
			"status":  "ok",
			"message": "error",
			"error":   err.Error(),
		})
	}

}

func WithGinContext(context *gin.Context) *Context {
	return &Context{Context: context}
}
