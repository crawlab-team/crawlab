package controllers

import (
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/core/schedule"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func PostSchedule(c *gin.Context) {
	var s models.ScheduleV2
	if err := c.ShouldBindJSON(&s); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	u := GetUserFromContextV2(c)

	modelSvc := service.NewModelServiceV2[models.ScheduleV2]()

	s.SetCreated(u.Id)
	s.SetUpdated(u.Id)
	id, err := modelSvc.InsertOne(s)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	s.Id = id

	if s.Enabled {
		scheduleSvc, err := schedule.GetScheduleServiceV2()
		if err != nil {
			HandleErrorInternalServerError(c, err)
			return
		}
		if err := scheduleSvc.Enable(s, u.Id); err != nil {
			HandleErrorInternalServerError(c, err)
			return
		}
	}

	HandleSuccessWithData(c, s)
}

func PutScheduleById(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	var s models.ScheduleV2
	if err := c.ShouldBindJSON(&s); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	if s.Id != id {
		HandleErrorBadRequest(c, errors.ErrorHttpBadRequest)
		return
	}

	modelSvc := service.NewModelServiceV2[models.ScheduleV2]()
	err = modelSvc.ReplaceById(id, s)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	scheduleSvc, err := schedule.GetScheduleServiceV2()
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	u := GetUserFromContextV2(c)

	if s.Enabled {
		if err := scheduleSvc.Enable(s, u.Id); err != nil {
			HandleErrorInternalServerError(c, err)
			return
		}
	} else {
		if err := scheduleSvc.Disable(s, u.Id); err != nil {
			HandleErrorInternalServerError(c, err)
			return
		}
	}

	HandleSuccessWithData(c, s)
}

func PostScheduleEnable(c *gin.Context) {
	postScheduleEnableDisableFunc(true)(c)
}

func PostScheduleDisable(c *gin.Context) {
	postScheduleEnableDisableFunc(false)(c)
}

func postScheduleEnableDisableFunc(isEnable bool) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			HandleErrorBadRequest(c, err)
			return
		}
		svc, err := schedule.GetScheduleServiceV2()
		if err != nil {
			HandleErrorInternalServerError(c, err)
			return
		}
		s, err := service.NewModelServiceV2[models.ScheduleV2]().GetById(id)
		if err != nil {
			HandleErrorInternalServerError(c, err)
			return
		}
		u := GetUserFromContextV2(c)
		if isEnable {
			err = svc.Enable(*s, u.Id)
		} else {
			err = svc.Disable(*s, u.Id)
		}
		if err != nil {
			HandleErrorInternalServerError(c, err)
			return
		}
		HandleSuccess(c)
	}
}
