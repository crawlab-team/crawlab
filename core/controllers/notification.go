package controllers

import (
	"github.com/crawlab-team/crawlab/core/notification"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

var NotificationController ActionController

func getNotificationActions() []Action {
	ctx := newNotificationContext()
	return []Action{
		{
			Method:      http.MethodGet,
			Path:        "/settings",
			HandlerFunc: ctx.GetSettingList,
		},
		{
			Method:      http.MethodGet,
			Path:        "/settings/:id",
			HandlerFunc: ctx.GetSetting,
		},
		{
			Method:      http.MethodPost,
			Path:        "/settings",
			HandlerFunc: ctx.PostSetting,
		},
		{
			Method:      http.MethodPut,
			Path:        "/settings/:id",
			HandlerFunc: ctx.PutSetting,
		},
		{
			Method:      http.MethodDelete,
			Path:        "/settings/:id",
			HandlerFunc: ctx.DeleteSetting,
		},
		{
			Method:      http.MethodPost,
			Path:        "/settings/:id/enable",
			HandlerFunc: ctx.EnableSetting,
		},
		{
			Method:      http.MethodPost,
			Path:        "/settings/:id/disable",
			HandlerFunc: ctx.DisableSetting,
		},
	}
}

type notificationContext struct {
	svc *notification.Service
}

func (ctx *notificationContext) GetSettingList(c *gin.Context) {
	query := MustGetFilterQuery(c)
	pagination := MustGetPagination(c)
	sort := MustGetSortOption(c)
	res, total, err := ctx.svc.GetSettingList(query, pagination, sort)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccessWithListData(c, res, total)
}

func (ctx *notificationContext) GetSetting(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	res, err := ctx.svc.GetSetting(id)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccessWithData(c, res)
}

func (ctx *notificationContext) PostSetting(c *gin.Context) {
	var s notification.NotificationSetting
	if err := c.ShouldBindJSON(&s); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	if err := ctx.svc.PosSetting(&s); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccess(c)
}

func (ctx *notificationContext) PutSetting(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	var s notification.NotificationSetting
	if err := c.ShouldBindJSON(&s); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	if err := ctx.svc.PutSetting(id, s); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccess(c)
}

func (ctx *notificationContext) DeleteSetting(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	if err := ctx.svc.DeleteSetting(id); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccess(c)
}

func (ctx *notificationContext) EnableSetting(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	if err := ctx.svc.EnableSetting(id); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccess(c)
}

func (ctx *notificationContext) DisableSetting(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	if err := ctx.svc.DisableSetting(id); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccess(c)
}

func newNotificationContext() *notificationContext {
	ctx := &notificationContext{
		svc: notification.GetService(),
	}
	return ctx
}
