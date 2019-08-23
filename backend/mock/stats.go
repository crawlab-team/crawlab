package mock

import (
	"crawlab/model"
	"github.com/gin-gonic/gin"
	"net/http"
)



var taskDailyItems = []model.TaskDailyItem{
	{
		Date:               "2019/08/19",
		TaskCount:          2,
		AvgRuntimeDuration: 1000,
	},
	{
		Date:               "2019/08/20",
		TaskCount:          3,
		AvgRuntimeDuration: 10130,
	},
}

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
	taskCount := 10

	// 在线节点总数
	activeNodeCount := 4

	// 爬虫总数
	spiderCount := 5
	// 定时任务数
	scheduleCount := 2

	// 每日任务数
	items := taskDailyItems

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
