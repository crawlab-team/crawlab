package routes

import (
	"crawlab/model"
	"crawlab/services"
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"net/http"
	"runtime/debug"
)

// @Summary Get schedule list
// @Description Get schedule list
// @Tags schedule
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /schedules [get]
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

// @Summary Get schedule by id
// @Description Get schedule  by id
// @Tags schedule
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id path string true "schedule id"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /schedules/{id} [get]
func GetSchedule(c *gin.Context) {
	id := c.Param("id")
	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
		return
	}
	result, err := model.GetSchedule(bson.ObjectIdHex(id))
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	HandleSuccessData(c, result)
}

// @Summary Post schedule
// @Description Post schedule
// @Tags schedule
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id path string true "schedule id"
// @Param newItem body  model.Schedule true "schedule item"
// @Success 200 json string Response
// @Failure 500 json string Response
// @Router /schedules/{id} [post]
func PostSchedule(c *gin.Context) {
	id := c.Param("id")
	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
		return
	}
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

// @Summary Put schedule
// @Description Put schedule
// @Tags schedule
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param item body  model.Schedule true "schedule item"
// @Success 200 json string Response
// @Failure 500 json string Response
// @Router /schedules [put]
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

// @Summary Delete schedule
// @Description Delete schedule
// @Tags schedule
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id path string true "schedule id"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /schedules/{id} [delete]
func DeleteSchedule(c *gin.Context) {
	id := c.Param("id")
	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
		return
	}
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
// @Summary disable schedule
// @Description disable schedule
// @Tags schedule
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id path string true "schedule id"
// @Success 200 json string Response
// @Failure 500 json string Response
// @Router /schedules/{id}/disable [post]
func DisableSchedule(c *gin.Context) {
	id := c.Param("id")
	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
		return
	}
	if err := services.Sched.Disable(bson.ObjectIdHex(id)); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
	HandleSuccess(c)
}

// 运行定时任务
// @Summary enable schedule
// @Description enable schedule
// @Tags schedule
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id path string true "schedule id"
// @Success 200 json string Response
// @Failure 500 json string Response
// @Router /schedules/{id}/enable [post]
func EnableSchedule(c *gin.Context) {
	id := c.Param("id")
	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
		return
	}
	if err := services.Sched.Enable(bson.ObjectIdHex(id)); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
	HandleSuccess(c)
}

func PutBatchSchedules(c *gin.Context) {
	var schedules []model.Schedule
	if err := c.ShouldBindJSON(&schedules); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	for _, s := range schedules {
		// 验证cron表达式
		if err := services.ParserCron(s.Cron); err != nil {
			log.Errorf("parse cron error: " + err.Error())
			debug.PrintStack()
			HandleError(http.StatusInternalServerError, c, err)
			return
		}

		// 添加 UserID
		s.UserId = services.GetCurrentUserId(c)

		// 添加定时任务
		if err := model.AddSchedule(s); err != nil {
			log.Errorf("add schedule error: " + err.Error())
			debug.PrintStack()
			HandleError(http.StatusInternalServerError, c, err)
			return
		}
	}

	// 更新定时任务
	if err := services.Sched.Update(); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	HandleSuccess(c)
}
