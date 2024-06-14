package controllers

import (
	"errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/db/mongo"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
)

type BaseControllerV2[T any] struct {
	modelSvc *service.ModelServiceV2[T]
	actions  []Action
}

func (ctr *BaseControllerV2[T]) GetById(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	model, err := ctr.modelSvc.GetById(id)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccessWithData(c, model)
}

func (ctr *BaseControllerV2[T]) GetList(c *gin.Context) {
	// get all if query field "all" is set true
	all := MustGetFilterAll(c)
	if all {
		ctr.getAll(c)
		return
	}

	// get list
	ctr.getList(c)
}

func (ctr *BaseControllerV2[T]) Post(c *gin.Context) {
	var model T
	if err := c.ShouldBindJSON(&model); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	u := GetUserFromContextV2(c)
	m := any(&model).(interfaces.ModelV2)
	m.SetId(primitive.NewObjectID())
	m.SetCreated(u.Id)
	m.SetUpdated(u.Id)
	col := ctr.modelSvc.GetCol()
	res, err := col.GetCollection().InsertOne(col.GetContext(), m)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	result, err := ctr.modelSvc.GetById(res.InsertedID.(primitive.ObjectID))
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccessWithData(c, result)
}

func (ctr *BaseControllerV2[T]) PutById(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	var model T
	if err := c.ShouldBindJSON(&model); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	u := GetUserFromContextV2(c)
	m := any(&model).(interfaces.ModelV2)
	m.SetId(primitive.NewObjectID())
	m.SetUpdated(u.Id)

	if err := ctr.modelSvc.ReplaceById(id, model); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	result, err := ctr.modelSvc.GetById(id)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccessWithData(c, result)
}

func (ctr *BaseControllerV2[T]) PatchList(c *gin.Context) {
	type Payload struct {
		Ids    []primitive.ObjectID `json:"ids"`
		Update bson.M               `json:"update"`
	}

	var payload Payload
	if err := c.ShouldBindJSON(&payload); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// query
	query := bson.M{
		"_id": bson.M{
			"$in": payload.Ids,
		},
	}

	// update
	if err := ctr.modelSvc.UpdateMany(query, bson.M{"$set": payload.Update}); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccess(c)
}

func (ctr *BaseControllerV2[T]) DeleteById(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	if err := ctr.modelSvc.DeleteById(id); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccess(c)
}

func (ctr *BaseControllerV2[T]) DeleteList(c *gin.Context) {
	type Payload struct {
		Ids []primitive.ObjectID `json:"ids"`
	}

	var payload Payload
	if err := c.ShouldBindJSON(&payload); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	if err := ctr.modelSvc.DeleteMany(bson.M{
		"_id": bson.M{
			"$in": payload.Ids,
		},
	}); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccess(c)
}

func (ctr *BaseControllerV2[T]) getAll(c *gin.Context) {
	models, err := ctr.modelSvc.GetMany(nil, &mongo.FindOptions{
		Sort: bson.D{{"_id", -1}},
	})
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	total, err := ctr.modelSvc.Count(nil)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccessWithListData(c, models, total)
}

func (ctr *BaseControllerV2[T]) getList(c *gin.Context) {
	// params
	pagination := MustGetPagination(c)
	query := MustGetFilterQuery(c)
	sort := MustGetSortOption(c)

	// get list
	models, err := ctr.modelSvc.GetMany(query, &mongo.FindOptions{
		Sort:  sort,
		Skip:  pagination.Size * (pagination.Page - 1),
		Limit: pagination.Size,
	})
	if err != nil {
		if errors.Is(err, mongo2.ErrNoDocuments) {
			HandleSuccessWithListData(c, nil, 0)
		} else {
			HandleErrorInternalServerError(c, err)
		}
		return
	}

	// total count
	total, err := ctr.modelSvc.Count(query)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// response
	HandleSuccessWithListData(c, models, total)
}

func NewControllerV2[T any](actions ...Action) *BaseControllerV2[T] {
	ctr := &BaseControllerV2[T]{
		modelSvc: service.NewModelServiceV2[T](),
		actions:  actions,
	}
	return ctr
}
