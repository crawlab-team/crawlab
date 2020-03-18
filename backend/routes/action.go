package routes

import (
	"crawlab/model"
	"crawlab/services"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"net/http"
)

func GetAction(c *gin.Context) {
	id := c.Param("id")

	user, err := model.GetAction(bson.ObjectIdHex(id))
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    user,
	})
}

func GetActionList(c *gin.Context) {
	pageNum := c.GetInt("page_num")
	pageSize := c.GetInt("page_size")

	users, err := model.GetActionList(nil, (pageNum-1)*pageSize, pageSize, "-create_ts")
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	total, err := model.GetActionListTotal(nil)
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, ListResponse{
		Status:  "ok",
		Message: "success",
		Data:    users,
		Total:   total,
	})
}

func PutAction(c *gin.Context) {
	// 绑定请求数据
	var action model.Action
	if err := c.ShouldBindJSON(&action); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	action.UserId = services.GetCurrentUser(c).Id

	if err := action.Add(); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

func PostAction(c *gin.Context) {
	id := c.Param("id")

	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
	}

	var item model.Action
	if err := c.ShouldBindJSON(&item); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	if err := model.UpdateAction(bson.ObjectIdHex(id), item); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

func DeleteAction(c *gin.Context) {
	id := c.Param("id")

	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
		return
	}

	// 从数据库中删除该爬虫
	if err := model.RemoveAction(bson.ObjectIdHex(id)); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}
