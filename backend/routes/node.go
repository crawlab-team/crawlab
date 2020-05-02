package routes

import (
	"crawlab/model"
	"crawlab/services"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"net/http"
)

// @Summary Get nodes
// @Description Get nodes
// @Tags node
// @Produce json
// @Param Authorization header string true "With the bearer started"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /nodes [get]
func GetNodeList(c *gin.Context) {
	nodes, err := model.GetNodeList(nil)
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	//for i, node := range nodes {
	//	nodes[i].IsMaster = services.IsMasterNode(node.Id.Hex())
	//}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    nodes,
	})
}

// @Summary Get node
// @Description Get node
// @Tags node
// @Produce json
// @Param Authorization header string true "With the bearer started"
// @Param id path string true "id"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /nodes/{id} [get]
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


// @Summary Post node
// @Description Post node
// @Tags node
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer started"
// @Param id path string true "post node"
// @Success 200 json string Response
// @Failure 500 json string Response
// @Router /nodes/{id} [post]
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

// @Summary Get tasks on node
// @Description Get tasks on node
// @Tags node
// @Produce json
// @Param Authorization header string true "With the bearer started"
// @Param id path string true "node id"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /nodes/{id}/tasks [get]
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

// @Summary Get system info
// @Description Get system info
// @Tags node
// @Produce json
// @Param Authorization header string true "With the bearer started"
// @Param id path string true "node id"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /nodes/{id}/system [get]
func GetSystemInfo(c *gin.Context) {
	id := c.Param("id")

	sysInfo, _ := services.GetSystemInfo(id)

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    sysInfo,
	})
}

// @Summary Delete node
// @Description Delete node
// @Tags node
// @Produce json
// @Param Authorization header string true "With the bearer started"
// @Param id path string true "node id"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /nodes/{id} [delete]
func DeleteNode(c *gin.Context) {
	id := c.Param("id")
	node, err := model.GetNode(bson.ObjectIdHex(id))
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
	err = node.Delete()
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}
