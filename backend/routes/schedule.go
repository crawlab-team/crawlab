package routes

import (
	"crawlab/constants"
	"crawlab/model"
	"crawlab/services"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"net/http"
)

func GetScheduleList(c *gin.Context) {
	results, err := model.GetScheduleList(nil)
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

func GetSchedule(c *gin.Context) {
	id := c.Param("id")

	result, err := model.GetSchedule(bson.ObjectIdHex(id))
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

func PostSchedule(c *gin.Context) {
	id := c.Param("id")

	// 绑定数据模型
	var newItem model.Schedule
	if err := c.ShouldBindJSON(&newItem); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}
	newItem.Id = bson.ObjectIdHex(id)

	// 如果node_id为空，则置为空ObjectId
	if newItem.NodeId == "" {
		newItem.NodeId = bson.ObjectIdHex(constants.ObjectIdNull)
	}

	// 更新数据库
	if err := model.UpdateSchedule(bson.ObjectIdHex(id), newItem); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 更新定时任务
	if err := services.Sched.Update(); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

func PutSchedule(c *gin.Context) {
	var item model.Schedule

	// 绑定数据模型
	if err := c.ShouldBindJSON(&item); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	// 如果node_id为空，则置为空ObjectId
	if item.NodeId == "" {
		item.NodeId = bson.ObjectIdHex(constants.ObjectIdNull)
	}

	// 更新数据库
	if err := model.AddSchedule(item); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 更新定时任务
	if err := services.Sched.Update(); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

func DeleteSchedule(c *gin.Context) {
	id := c.Param("id")

	// 删除定时任务
	if err := model.RemoveSchedule(bson.ObjectIdHex(id)); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 更新定时任务
	if err := services.Sched.Update(); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}
