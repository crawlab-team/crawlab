package controllers

import (
	"github.com/crawlab-team/crawlab/core/container"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"time"
)

var StatsController ActionController

func getStatsActions() []Action {
	statsCtx := newStatsContext()
	return []Action{
		{
			Method:      http.MethodGet,
			Path:        "/overview",
			HandlerFunc: statsCtx.getOverview,
		},
		{
			Method:      http.MethodGet,
			Path:        "/daily",
			HandlerFunc: statsCtx.getDaily,
		},
		{
			Method:      http.MethodGet,
			Path:        "/tasks",
			HandlerFunc: statsCtx.getTasks,
		},
	}
}

type statsContext struct {
	statsSvc     interfaces.StatsService
	defaultQuery bson.M
}

func (svc *statsContext) getOverview(c *gin.Context) {
	data, err := svc.statsSvc.GetOverviewStats(svc.defaultQuery)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccessWithData(c, data)
}

func (svc *statsContext) getDaily(c *gin.Context) {
	data, err := svc.statsSvc.GetDailyStats(svc.defaultQuery)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccessWithData(c, data)
}

func (svc *statsContext) getTasks(c *gin.Context) {
	data, err := svc.statsSvc.GetTaskStats(svc.defaultQuery)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccessWithData(c, data)
}

func newStatsContext() *statsContext {
	// context
	ctx := &statsContext{
		defaultQuery: bson.M{
			"create_ts": bson.M{
				"$gte": time.Now().Add(-30 * 24 * time.Hour),
			},
		},
	}

	// dependency injection
	if err := container.GetContainer().Invoke(func(
		statsSvc interfaces.StatsService,
	) {
		ctx.statsSvc = statsSvc
	}); err != nil {
		panic(err)
	}

	return ctx
}
