package routes

import (
	"crawlab/services"
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
)

func GetDocs(c *gin.Context) {
	type ResData struct {
		String string `json:"string"`
	}
	data, err := services.GetDocs()
	if err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
	}
	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    ResData{String:data},
	})
}
