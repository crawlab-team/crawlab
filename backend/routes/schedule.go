package routes

import (
	"crawlab/model"
	"crawlab/services"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"net/http"
)

func GetScheduleList(c *gin.Context) {
	query := bson.M{}

	// 获取校验
	query = services.GetAuthQuery(query, c)

	results, err := model.GetScheduleList(query)
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
	HandleSuccessData(c, results)
}

func GetSchedule(c *gin.Context) {
	id := c.Param("id")

	result, err := model.GetSchedule(bson.ObjectIdHex(id))
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	HandleSuccessData(c, result)
}

func PostSchedule(c *gin.Context) {
	id := c.Param("id")

	// 绑定数据模型
	var newItem model.Schedule
	if err := c.ShouldBindJSON(&newItem); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	// 验证cron表达式
	if err := services.ParserCron(newItem.Cron); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	newItem.Id = bson.ObjectIdHex(id)
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

	HandleSuccess(c)
}

func PutSchedule(c *gin.Context) {
	var item model.Schedule

	// 绑定数据模型
	if err := c.ShouldBindJSON(&item); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	// 验证cron表达式
	if err := services.ParserCron(item.Cron); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 加入用户ID
	item.UserId = services.GetCurrentUserId(c)

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

	HandleSuccess(c)
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

	HandleSuccess(c)
}

// 停止定时任务
func DisableSchedule(c *gin.Context) {
	id := c.Param("id")
	if err := services.Sched.Disable(bson.ObjectIdHex(id)); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
	HandleSuccess(c)
}

// 运行定时任务
func EnableSchedule(c *gin.Context) {
	id := c.Param("id")
	if err := services.Sched.Enable(bson.ObjectIdHex(id)); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
	HandleSuccess(c)
}
