package controllers

import (
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetSetting(c *gin.Context) {
	// key
	key := c.Param("id")

	// setting
	s, err := service.NewModelServiceV2[models.SettingV2]().GetOne(bson.M{"key": key}, nil)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccessWithData(c, s)
}

func PutSetting(c *gin.Context) {
	// key
	key := c.Param("id")

	// settings
	var s models.Setting
	if err := c.ShouldBindJSON(&s); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	modelSvc := service.NewModelServiceV2[models.SettingV2]()

	// setting
	_s, err := modelSvc.GetOne(bson.M{"key": key}, nil)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	u := GetUserFromContextV2(c)

	// save
	_s.Value = s.Value
	_s.SetUpdated(u.Id)
	err = modelSvc.ReplaceOne(bson.M{"key": key}, *_s)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccess(c)
}
