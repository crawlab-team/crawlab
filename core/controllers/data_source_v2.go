package controllers

import (
	"github.com/crawlab-team/crawlab/core/ds"
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/models/models/v2"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func PostDataSource(c *gin.Context) {
	// data source
	var payload struct {
		Name        string            `json:"name"`
		Type        string            `json:"type"`
		Description string            `json:"description"`
		Host        string            `json:"host"`
		Port        string            `json:"port"`
		Url         string            `json:"url"`
		Hosts       []string          `json:"hosts"`
		Database    string            `json:"database"`
		Username    string            `json:"username"`
		Password    string            `json:"-,omitempty"`
		ConnectType string            `json:"connect_type"`
		Status      string            `json:"status"`
		Error       string            `json:"error"`
		Extra       map[string]string `json:"extra,omitempty"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	u := GetUserFromContextV2(c)

	// add data source to db
	dataSource := models.DataSourceV2{
		Name:        payload.Name,
		Type:        payload.Type,
		Description: payload.Description,
		Host:        payload.Host,
		Port:        payload.Port,
		Url:         payload.Url,
		Hosts:       payload.Hosts,
		Database:    payload.Database,
		Username:    payload.Username,
		Password:    payload.Password,
		ConnectType: payload.ConnectType,
		Status:      payload.Status,
		Error:       payload.Error,
		Extra:       payload.Extra,
	}
	dataSource.SetCreated(u.Id)
	dataSource.SetUpdated(u.Id)
	id, err := service.NewModelServiceV2[models.DataSourceV2]().InsertOne(dataSource)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	dataSource.Id = id

	// check data source status
	go func() {
		_ = ds.GetDataSourceServiceV2().CheckStatus(id)
	}()

	HandleSuccessWithData(c, dataSource)
}

func PutDataSourceById(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// data source
	var dataSource models.DataSourceV2
	if err := c.ShouldBindJSON(&dataSource); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	err = service.NewModelServiceV2[models.DataSourceV2]().ReplaceById(id, dataSource)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// check data source status
	go func() {
		_ = ds.GetDataSourceServiceV2().CheckStatus(id)
	}()
}

func PostDataSourceChangePassword(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	var payload struct {
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	if payload.Password == "" {
		HandleErrorBadRequest(c, errors.ErrorDataSourceMissingRequiredFields)
		return
	}
	u := GetUserFromContextV2(c)
	if err := ds.GetDataSourceServiceV2().ChangePassword(id, payload.Password, u.Id); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccess(c)
}
