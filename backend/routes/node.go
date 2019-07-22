package routes

import (
	"crawlab/model"
	"crawlab/services"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"net/http"
)

func GetNodeList(c *gin.Context) {
	results, err := model.GetNodeList(nil)
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    results,
	})
}

func GetNode(c *gin.Context) {
	id := c.Param("id")

	result, err := model.GetNode(bson.ObjectIdHex(id))
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    result,
	})
}

func Ping(c *gin.Context) {
	data, err := services.GetNodeData()
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    data,
	})
}

func PostNode(c *gin.Context) {
	id := c.Param("id")

	item, err := model.GetNode(bson.ObjectIdHex(id))
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	var newItem model.Node
	if err := c.ShouldBindJSON(&newItem); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}
	newItem.Id = item.Id

	if err := model.UpdateNode(bson.ObjectIdHex(id), newItem); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

func GetNodeTaskList(c *gin.Context) {
	id := c.Param("id")

	tasks, err := model.GetNodeTaskList(bson.ObjectIdHex(id))
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    tasks,
	})
}

func GetSystemInfo(c *gin.Context) {
	sysInfo, _ := services.GetSystemInfo()

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    sysInfo,
	})
}
