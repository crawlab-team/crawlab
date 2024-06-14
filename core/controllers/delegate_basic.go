package controllers

import (
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	delegate2 "github.com/crawlab-team/crawlab/core/models/delegate"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
)

func NewBasicControllerDelegate(id ControllerId, svc interfaces.ModelBaseService) (d *BasicControllerDelegate) {
	return &BasicControllerDelegate{
		id:  id,
		svc: svc,
	}
}

type BasicControllerDelegate struct {
	id  ControllerId
	svc interfaces.ModelBaseService
}

func (d *BasicControllerDelegate) Get(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	doc, err := d.svc.GetById(id)
	if err == mongo2.ErrNoDocuments {
		HandleErrorNotFound(c, err)
		return
	}
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccessWithData(c, doc)
}

func (d *BasicControllerDelegate) Post(c *gin.Context) {
	doc, err := NewJsonBinder(d.id).Bind(c)
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	if err := delegate2.NewModelDelegate(doc, GetUserFromContext(c)).Add(); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccessWithData(c, doc)
}

func (d *BasicControllerDelegate) Put(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	doc, err := NewJsonBinder(d.id).Bind(c)
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	if doc.GetId() != id {
		HandleErrorBadRequest(c, errors.ErrorHttpBadRequest)
		return
	}
	_, err = d.svc.GetById(id)
	if err != nil {
		HandleErrorNotFound(c, err)
		return
	}
	if err := delegate2.NewModelDelegate(doc, GetUserFromContext(c)).Save(); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccessWithData(c, doc)
}

func (d *BasicControllerDelegate) Delete(c *gin.Context) {
	id := c.Param("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	doc, err := d.svc.GetById(oid)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	if err := delegate2.NewModelDelegate(doc, GetUserFromContext(c)).Delete(); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccess(c)
}
