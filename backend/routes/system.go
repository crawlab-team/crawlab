package routes

import (
	"crawlab/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetLangList(c *gin.Context) {
	nodeId := c.Param("id")
	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    services.GetLangList(nodeId),
	})
}

func GetDepList(c *gin.Context) {
	nodeId := c.Param("id")
	lang := c.Query("lang")
	depName := c.Query("dep_name")
	depList, err := services.GetDepList(nodeId, lang, depName)
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    depList,
	})
}
