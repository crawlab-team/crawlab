package routes

import (
	"crawlab/model"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"net/http"
)

// 新增
func PutVariable(c *gin.Context) {
	var variable model.Variable
	if err := c.ShouldBindJSON(&variable); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}
	if err := variable.Add(); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
	HandleSuccess(c)
}

// 修改
func PostVariable(c *gin.Context) {
	var id = c.Param("id")
	var variable model.Variable
	if err := c.ShouldBindJSON(&variable); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}
	variable.Id = bson.ObjectIdHex(id)
	if err := variable.Save(); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
	HandleSuccess(c)
}

// 删除
func DeleteVariable(c *gin.Context) {
	var idStr = c.Param("id")
	var id = bson.ObjectIdHex(idStr)
	variable, err := model.GetVariable(id)
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
	variable.Id = id
	if err := variable.Delete(); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
	HandleSuccess(c)

}

// 列表
func GetVariableList(c *gin.Context) {
	list := model.GetVariableList()
	HandleSuccessData(c, list)
}
