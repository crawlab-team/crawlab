package routes

import (
	"crawlab/model"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"net/http"
)

// 新增

// @Summary Put variable
// @Description Put variable
// @Tags variable
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param variable body model.Variable true "reqData body"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /variable [put]
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

// @Summary Post variable
// @Description Post variable
// @Tags variable
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param variable body model.Variable true "reqData body"
// @Param id path string true "variable id"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /variable/{id} [post]
func PostVariable(c *gin.Context) {
	var id = c.Param("id")

	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
		return
	}
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

// @Summary Delete variable
// @Description Delete variable
// @Tags variable
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id path string true "variable id"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /variable/{id} [delete]
func DeleteVariable(c *gin.Context) {
	var idStr = c.Param("id")
	if !bson.IsObjectIdHex(idStr) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
		return
	}
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

// @Summary Get variable  list
// @Description Get variable  list
// @Tags variable
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /variables [get]
func GetVariableList(c *gin.Context) {
	list := model.GetVariableList()
	HandleSuccessData(c, list)
}
