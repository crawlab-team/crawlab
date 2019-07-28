package routes

import (
	"crawlab/constants"
	"crawlab/model"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"net/http"
)

func GetHomeStats(c *gin.Context) {
	type DataOverview struct {
		TaskCount       int `json:"task_count"`
		SpiderCount     int `json:"spider_count"`
		ActiveNodeCount int `json:"active_node_count"`
		ScheduleCount   int `json:"schedule_count"`
	}

	type Data struct {
		Overview DataOverview          `json:"overview"`
		Daily    []model.TaskDailyItem `json:"daily"`
	}

	// 任务总数
	taskCount, err := model.GetTaskCount(nil)
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 在线节点总数
	activeNodeCount, err := model.GetNodeCount(bson.M{"status": constants.StatusOnline})
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 爬虫总数
	spiderCount, err := model.GetSpiderCount()
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 定时任务数
	scheduleCount, err := model.GetScheduleCount()
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 每日任务数
	items, err := model.GetDailyTaskStats(bson.M{})
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data: Data{
			Overview: DataOverview{
				ActiveNodeCount: activeNodeCount,
				TaskCount:       taskCount,
				SpiderCount:     spiderCount,
				ScheduleCount:   scheduleCount,
			},
			Daily: items,
		},
	})
}
