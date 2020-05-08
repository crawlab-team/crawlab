package routes

import (
	"crawlab/services"
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
)

// @Summary Get docs
// @Description Get docs
// @Tags docs
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /docs [get]
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
