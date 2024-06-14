package controllers

import (
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/delegate"
	"github.com/crawlab-team/crawlab/core/utils"
	"github.com/crawlab-team/crawlab/db/mongo"
	"github.com/crawlab-team/crawlab/trace"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"reflect"
	"time"
)

func NewListControllerDelegate(id ControllerId, svc interfaces.ModelBaseService) (d *ListControllerDelegate) {
	if svc == nil {
		panic(errors.ErrorControllerNoModelService)
	}

	return &ListControllerDelegate{
		id:  id,
		svc: svc,
		bc:  NewBasicControllerDelegate(id, svc),
	}
}

type ListControllerDelegate struct {
	id  ControllerId
	svc interfaces.ModelBaseService
	bc  BasicController
}

func (d *ListControllerDelegate) Get(c *gin.Context) {
	d.bc.Get(c)
}

func (d *ListControllerDelegate) Post(c *gin.Context) {
	d.bc.Post(c)
}

func (d *ListControllerDelegate) Put(c *gin.Context) {
	d.bc.Put(c)
}

func (d *ListControllerDelegate) Delete(c *gin.Context) {
	d.bc.Delete(c)
}

func (d *ListControllerDelegate) GetList(c *gin.Context) {
	// get all if query field "all" is set true
	all := MustGetFilterAll(c)
	if all {
		d.getAll(c)
		return
	}

	// get list and total
	l, total, err := d.getList(c)
	if err != nil {
		return
	}

	// response
	HandleSuccessWithListData(c, l, total)
}

func (d *ListControllerDelegate) PostList(c *gin.Context) {
	// bind
	docs, err := NewJsonBinder(d.id).BindList(c)
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// success ids
	var ids []primitive.ObjectID

	// reflect
	switch reflect.TypeOf(docs).Kind() {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(docs)
		for i := 0; i < s.Len(); i++ {
			item := s.Index(i)
			if !item.CanAddr() {
				HandleErrorInternalServerError(c, errors.ErrorModelInvalidType)
				return
			}
			ptr := item.Addr()
			doc, ok := ptr.Interface().(interfaces.Model)
			if !ok {
				HandleErrorInternalServerError(c, errors.ErrorModelInvalidType)
				return
			}
			if err := delegate.NewModelDelegate(doc, GetUserFromContext(c)).Add(); err != nil {
				_ = trace.TraceError(err)
				continue
			}
			ids = append(ids, doc.GetId())
		}
	}

	// check
	items, err := utils.GetArrayItems(docs)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	if len(ids) < len(items) {
		HandleErrorInternalServerError(c, errors.ErrorControllerAddError)
		return
	}

	// success
	HandleSuccessWithData(c, docs)
}

func (d *ListControllerDelegate) PutList(c *gin.Context) {
	payload, doc, err := NewJsonBinder(d.id).BindBatchRequestPayloadWithStringData(c)
	if err != nil {
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
	if err := d.svc.UpdateDoc(query, doc, payload.Fields); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccess(c)
}

func (d *ListControllerDelegate) DeleteList(c *gin.Context) {
	payload, err := NewJsonBinder(d.id).BindBatchRequestPayload(c)
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	if err := d.svc.DeleteList(bson.M{
		"_id": bson.M{
			"$in": payload.Ids,
		},
	}); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccess(c)
}

func (d *ListControllerDelegate) getAll(c *gin.Context) {
	// get list
	tic := time.Now()
	log.Debugf("getAll -> d.svc.GetMany:start")
	list, err := d.svc.GetList(nil, &mongo.FindOptions{
		Sort: bson.D{{"_id", -1}},
	})
	if err != nil {
		if err == mongo2.ErrNoDocuments {
			HandleErrorNotFound(c, err)
		} else {
			HandleErrorInternalServerError(c, err)
		}
		return
	}
	log.Debugf("getAll -> d.svc.GetMany:end. elapsed: %d ms", time.Now().Sub(tic).Milliseconds())
	tic = time.Now()

	// total count
	tic = time.Now()
	log.Debugf("getAll -> d.svc.Count:start")
	total, err := d.svc.Count(nil)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	log.Debugf("getAll -> d.svc.Count:end. elapsed: %d ms", time.Now().Sub(tic).Milliseconds())

	// response
	HandleSuccessWithListData(c, list, total)
}

func (d *ListControllerDelegate) getList(c *gin.Context) (l interfaces.List, total int, err error) {
	// params
	pagination := MustGetPagination(c)
	query := MustGetFilterQuery(c)
	sort := MustGetSortOption(c)

	// get list
	l, err = d.svc.GetList(query, &mongo.FindOptions{
		Sort:  sort,
		Skip:  pagination.Size * (pagination.Page - 1),
		Limit: pagination.Size,
	})
	if err != nil {
		if err.Error() == mongo2.ErrNoDocuments.Error() {
			HandleSuccessWithListData(c, nil, 0)
		} else {
			HandleErrorInternalServerError(c, err)
		}
		return
	}

	// total count
	total, err = d.svc.Count(query)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	return l, total, nil
}
