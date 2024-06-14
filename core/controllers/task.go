package controllers

import (
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/container"
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	delegate2 "github.com/crawlab-team/crawlab/core/models/delegate"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/core/result"
	"github.com/crawlab-team/crawlab/core/task/log"
	"github.com/crawlab-team/crawlab/core/utils"
	"github.com/crawlab-team/crawlab/db/generic"
	"github.com/crawlab-team/crawlab/db/mongo"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strings"
)

var TaskController *taskController

func getTaskActions() []Action {
	taskCtx := newTaskContext()
	return []Action{
		{
			Method:      http.MethodPost,
			Path:        "/run",
			HandlerFunc: taskCtx.run,
		},
		{
			Method:      http.MethodPost,
			Path:        "/:id/restart",
			HandlerFunc: taskCtx.restart,
		},
		{
			Method:      http.MethodPost,
			Path:        "/:id/cancel",
			HandlerFunc: taskCtx.cancel,
		},
		{
			Method:      http.MethodGet,
			Path:        "/:id/logs",
			HandlerFunc: taskCtx.getLogs,
		},
		{
			Method:      http.MethodGet,
			Path:        "/:id/data",
			HandlerFunc: taskCtx.getData,
		},
	}
}

type taskController struct {
	ListActionControllerDelegate
	d   ListActionControllerDelegate
	ctx *taskContext
}

func (ctr *taskController) Get(c *gin.Context) {
	ctr.ctx.getWithStatsSpider(c)
}

func (ctr *taskController) Delete(c *gin.Context) {
	if err := ctr.ctx._delete(c); err != nil {
		return
	}
	HandleSuccess(c)
}

func (ctr *taskController) GetList(c *gin.Context) {
	withStats := c.Query("stats")
	if withStats == "" {
		ctr.d.GetList(c)
		return
	}
	ctr.ctx.getListWithStats(c)
}

func (ctr *taskController) DeleteList(c *gin.Context) {
	if err := ctr.ctx._deleteList(c); err != nil {
		return
	}
	HandleSuccess(c)
}

type taskContext struct {
	modelSvc         service.ModelService
	modelTaskSvc     interfaces.ModelBaseService
	modelTaskStatSvc interfaces.ModelBaseService
	adminSvc         interfaces.SpiderAdminService
	schedulerSvc     interfaces.TaskSchedulerService
	l                log.Driver
}

func (ctx *taskContext) run(c *gin.Context) {
	// task
	var t models.Task
	if err := c.ShouldBindJSON(&t); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// validate spider id
	if t.GetSpiderId().IsZero() {
		HandleErrorBadRequest(c, errors.ErrorTaskEmptySpiderId)
		return
	}

	// spider
	s, err := ctx.modelSvc.GetSpiderById(t.GetSpiderId())
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// options
	opts := &interfaces.SpiderRunOptions{
		Mode:     t.Mode,
		NodeIds:  t.NodeIds,
		Cmd:      t.Cmd,
		Param:    t.Param,
		Priority: t.Priority,
	}

	// user
	if u := GetUserFromContext(c); u != nil {
		opts.UserId = u.GetId()
	}

	// run
	taskIds, err := ctx.adminSvc.Schedule(s.GetId(), opts)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccessWithData(c, taskIds)
}

func (ctx *taskContext) restart(c *gin.Context) {
	// id
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// task
	t, err := ctx.modelSvc.GetTaskById(id)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// options
	opts := &interfaces.SpiderRunOptions{
		Mode:     t.Mode,
		NodeIds:  t.NodeIds,
		Cmd:      t.Cmd,
		Param:    t.Param,
		Priority: t.Priority,
	}

	// user
	if u := GetUserFromContext(c); u != nil {
		opts.UserId = u.GetId()
	}

	// run
	taskIds, err := ctx.adminSvc.Schedule(t.SpiderId, opts)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccessWithData(c, taskIds)
}

func (ctx *taskContext) cancel(c *gin.Context) {
	// id
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// task
	t, err := ctx.modelSvc.GetTaskById(id)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// validate
	if !utils.IsCancellable(t.Status) {
		HandleErrorInternalServerError(c, errors.ErrorControllerNotCancellable)
		return
	}

	// cancel
	if err := ctx.schedulerSvc.Cancel(id, GetUserFromContext(c)); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccess(c)
}

func (ctx *taskContext) getLogs(c *gin.Context) {
	// id
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// pagination
	p, err := GetPagination(c)
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// logs
	logs, err := ctx.l.Find(id.Hex(), "", (p.Page-1)*p.Size, p.Size)
	if err != nil {
		if strings.HasSuffix(err.Error(), "Status:404 Not Found") {
			HandleSuccess(c)
			return
		}
		HandleErrorInternalServerError(c, err)
		return
	}
	total, err := ctx.l.Count(id.Hex(), "")
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccessWithListData(c, logs, total)
}

func (ctx *taskContext) getListWithStats(c *gin.Context) {
	// params
	pagination := MustGetPagination(c)
	query := MustGetFilterQuery(c)
	sort := MustGetSortOption(c)

	// get list
	list, err := ctx.modelTaskSvc.GetList(query, &mongo.FindOptions{
		Sort:  sort,
		Skip:  pagination.Size * (pagination.Page - 1),
		Limit: pagination.Size,
	})
	if err != nil {
		if err == mongo2.ErrNoDocuments {
			HandleErrorNotFound(c, err)
		} else {
			HandleErrorInternalServerError(c, err)
		}
		return
	}

	// check empty list
	if len(list.GetModels()) == 0 {
		HandleSuccessWithListData(c, nil, 0)
		return
	}

	// ids
	var ids []primitive.ObjectID
	for _, d := range list.GetModels() {
		t := d.(interfaces.Model)
		ids = append(ids, t.GetId())
	}

	// total count
	total, err := ctx.modelTaskSvc.Count(query)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// stat list
	query = bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	}
	stats, err := ctx.modelSvc.GetTaskStatList(query, nil)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// cache stat list to dict
	dict := map[primitive.ObjectID]models.TaskStat{}
	for _, s := range stats {
		dict[s.GetId()] = s
	}

	// iterate list again
	var data []interface{}
	for _, d := range list.GetModels() {
		t := d.(*models.Task)
		s, ok := dict[t.GetId()]
		if ok {
			t.Stat = &s
		}
		data = append(data, *t)
	}

	// response
	HandleSuccessWithListData(c, data, total)
}

func (ctx *taskContext) getWithStatsSpider(c *gin.Context) {
	// id
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// task
	t, err := ctx.modelSvc.GetTaskById(id)
	if err == mongo2.ErrNoDocuments {
		HandleErrorNotFound(c, err)
		return
	}
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// spider
	t.Spider, _ = ctx.modelSvc.GetSpiderById(t.SpiderId)

	// skip if task status is pending
	if t.Status == constants.TaskStatusPending {
		HandleSuccessWithData(c, t)
		return
	}

	// task stat
	t.Stat, _ = ctx.modelSvc.GetTaskStatById(id)

	HandleSuccessWithData(c, t)
}

func (ctx *taskContext) getData(c *gin.Context) {
	// id
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// pagination
	p, err := GetPagination(c)
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// task
	t, err := ctx.modelSvc.GetTaskById(id)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// result service
	resultSvc, err := result.GetResultService(t.SpiderId)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// query
	query := generic.ListQuery{
		generic.ListQueryCondition{
			Key:   constants.TaskKey,
			Op:    generic.OpEqual,
			Value: t.Id,
		},
	}

	// list
	data, err := resultSvc.List(query, &generic.ListOptions{
		Skip:  (p.Page - 1) * p.Size,
		Limit: p.Size,
		Sort:  []generic.ListSort{{"_id", generic.SortDirectionDesc}},
	})
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// total
	total, err := resultSvc.Count(query)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccessWithListData(c, data, total)
}

func (ctx *taskContext) _delete(c *gin.Context) (err error) {
	id := c.Param("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	if err := mongo.RunTransaction(func(context mongo2.SessionContext) (err error) {
		// delete task
		task, err := ctx.modelSvc.GetTaskById(oid)
		if err != nil {
			return err
		}
		if err := delegate2.NewModelDelegate(task, GetUserFromContext(c)).Delete(); err != nil {
			return err
		}

		// delete task stat
		taskStat, err := ctx.modelSvc.GetTaskStatById(oid)
		if err != nil {
			return err
		}
		if err := delegate2.NewModelDelegate(taskStat, GetUserFromContext(c)).Delete(); err != nil {
			return err
		}

		return nil
	}); err != nil {
		HandleErrorInternalServerError(c, err)
		return err
	}

	return nil
}

func (ctx *taskContext) _deleteList(c *gin.Context) (err error) {
	payload, err := NewJsonBinder(ControllerIdTask).BindBatchRequestPayload(c)
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	if err := mongo.RunTransaction(func(context mongo2.SessionContext) error {
		// delete tasks
		if err := ctx.modelTaskSvc.DeleteList(bson.M{
			"_id": bson.M{
				"$in": payload.Ids,
			},
		}); err != nil {
			return err
		}

		// delete task stats
		if err := ctx.modelTaskStatSvc.DeleteList(bson.M{
			"_id": bson.M{
				"$in": payload.Ids,
			},
		}); err != nil {
			return err
		}

		return nil
	}); err != nil {
		HandleErrorInternalServerError(c, err)
		return err
	}

	return nil
}

func newTaskContext() *taskContext {
	// context
	ctx := &taskContext{}

	// dependency injection
	if err := container.GetContainer().Invoke(func(
		modelSvc service.ModelService,
		adminSvc interfaces.SpiderAdminService,
		schedulerSvc interfaces.TaskSchedulerService,
	) {
		ctx.modelSvc = modelSvc
		ctx.adminSvc = adminSvc
		ctx.schedulerSvc = schedulerSvc
	}); err != nil {
		panic(err)
	}

	// model task service
	ctx.modelTaskSvc = ctx.modelSvc.GetBaseService(interfaces.ModelIdTask)

	// model task stat service
	ctx.modelTaskStatSvc = ctx.modelSvc.GetBaseService(interfaces.ModelIdTaskStat)

	// log driver
	l, err := log.GetLogDriver(log.DriverTypeFile)
	if err != nil {
		panic(err)
	}
	ctx.l = l

	return ctx
}

func newTaskController() *taskController {
	actions := getTaskActions()
	modelSvc, err := service.GetService()
	if err != nil {
		panic(err)
	}

	ctr := NewListPostActionControllerDelegate(ControllerIdTask, modelSvc.GetBaseService(interfaces.ModelIdTask), actions)
	d := NewListPostActionControllerDelegate(ControllerIdTask, modelSvc.GetBaseService(interfaces.ModelIdTask), actions)
	ctx := newTaskContext()

	return &taskController{
		ListActionControllerDelegate: *ctr,
		d:                            *d,
		ctx:                          ctx,
	}
}
