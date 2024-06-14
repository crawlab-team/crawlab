package controllers

import (
	"github.com/crawlab-team/crawlab/core/stats"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

var statsDefaultQuery = bson.M{
	"create_ts": bson.M{
		"$gte": time.Now().Add(-30 * 24 * time.Hour),
	},
}

func GetStatsOverview(c *gin.Context) {
	data, err := stats.GetStatsService().GetOverviewStats(statsDefaultQuery)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccessWithData(c, data)
}

func GetStatsDaily(c *gin.Context) {
	data, err := stats.GetStatsService().GetDailyStats(statsDefaultQuery)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccessWithData(c, data)
}

func GetStatsTasks(c *gin.Context) {
	data, err := stats.GetStatsService().GetTaskStats(statsDefaultQuery)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccessWithData(c, data)
}
