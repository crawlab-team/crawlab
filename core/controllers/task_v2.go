package controllers

import (
	"errors"
	log2 "github.com/apex/log"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/models/v2"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/core/result"
	"github.com/crawlab-team/crawlab/core/spider/admin"
	"github.com/crawlab-team/crawlab/core/task/log"
	"github.com/crawlab-team/crawlab/core/task/scheduler"
	"github.com/crawlab-team/crawlab/core/utils"
	"github.com/crawlab-team/crawlab/db/generic"
	"github.com/crawlab-team/crawlab/db/mongo"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func GetTaskById(c *gin.Context) {
	// id
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// task
	t, err := service.NewModelServiceV2[models.TaskV2]().GetById(id)
	if errors.Is(err, mongo2.ErrNoDocuments) {
		HandleErrorNotFound(c, err)
		return
	}
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// spider
	t.Spider, _ = service.NewModelServiceV2[models.SpiderV2]().GetById(t.SpiderId)

	// skip if task status is pending
	if t.Status == constants.TaskStatusPending {
		HandleSuccessWithData(c, t)
		return
	}

	// task stat
	t.Stat, _ = service.NewModelServiceV2[models.TaskStatV2]().GetById(id)

	HandleSuccessWithData(c, t)
}

func GetTaskList(c *gin.Context) {
	withStats := c.Query("stats")
	if withStats == "" {
		NewControllerV2[models.TaskV2]().GetList(c)
		return
	}

	// params
	pagination := MustGetPagination(c)
	query := MustGetFilterQuery(c)
	sort := MustGetSortOption(c)

	// get tasks
	tasks, err := service.NewModelServiceV2[models.TaskV2]().GetMany(query, &mongo.FindOptions{
		Sort:  sort,
		Skip:  pagination.Size * (pagination.Page - 1),
		Limit: pagination.Size,
	})
	if err != nil {
		if errors.Is(err, mongo2.ErrNoDocuments) {
			HandleErrorNotFound(c, err)
		} else {
			HandleErrorInternalServerError(c, err)
		}
		return
	}

	// check empty list
	if len(tasks) == 0 {
		HandleSuccessWithListData(c, nil, 0)
		return
	}

	// ids
	var taskIds []primitive.ObjectID
	var spiderIds []primitive.ObjectID
	for _, t := range tasks {
		taskIds = append(taskIds, t.Id)
		spiderIds = append(spiderIds, t.SpiderId)
	}

	// total count
	total, err := service.NewModelServiceV2[models.TaskV2]().Count(query)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// stat list
	stats, err := service.NewModelServiceV2[models.TaskStatV2]().GetMany(bson.M{
		"_id": bson.M{
			"$in": taskIds,
		},
	}, nil)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// cache stat list to dict
	statsDict := map[primitive.ObjectID]models.TaskStatV2{}
	for _, s := range stats {
		statsDict[s.Id] = s
	}

	// spider list
	spiders, err := service.NewModelServiceV2[models.SpiderV2]().GetMany(bson.M{
		"_id": bson.M{
			"$in": spiderIds,
		},
	}, nil)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// cache spider list to dict
	spiderDict := map[primitive.ObjectID]models.SpiderV2{}
	for _, s := range spiders {
		spiderDict[s.Id] = s
	}

	// iterate list again
	for i, t := range tasks {
		// task stat
		ts, ok := statsDict[t.Id]
		if ok {
			tasks[i].Stat = &ts
		}

		// spider
		s, ok := spiderDict[t.SpiderId]
		if ok {
			tasks[i].Spider = &s
		}
	}

	// response
	HandleSuccessWithListData(c, tasks, total)
}

func DeleteTaskById(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// delete in db
	if err := mongo.RunTransaction(func(context mongo2.SessionContext) (err error) {
		// delete task
		_, err = service.NewModelServiceV2[models.TaskV2]().GetById(id)
		if err != nil {
			return err
		}
		err = service.NewModelServiceV2[models.TaskV2]().DeleteById(id)
		if err != nil {
			return err
		}

		// delete task stat
		_, err = service.NewModelServiceV2[models.TaskStatV2]().GetById(id)
		if err != nil {
			log2.Warnf("delete task stat error: %s", err.Error())
			return nil
		}
		err = service.NewModelServiceV2[models.TaskStatV2]().DeleteById(id)
		if err != nil {
			log2.Warnf("delete task stat error: %s", err.Error())
			return nil
		}

		return nil
	}); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// delete task logs
	logPath := filepath.Join(viper.GetString("log.path"), id.Hex())
	if err := os.RemoveAll(logPath); err != nil {
		log2.Warnf("failed to remove task log directory: %s", logPath)
	}

	HandleSuccess(c)
}

func DeleteList(c *gin.Context) {
	var payload struct {
		Ids []primitive.ObjectID `json:"ids"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	if err := mongo.RunTransaction(func(context mongo2.SessionContext) error {
		// delete tasks
		if err := service.NewModelServiceV2[models.TaskV2]().DeleteMany(bson.M{
			"_id": bson.M{
				"$in": payload.Ids,
			},
		}); err != nil {
			return err
		}

		// delete task stats
		if err := service.NewModelServiceV2[models.TaskV2]().DeleteMany(bson.M{
			"_id": bson.M{
				"$in": payload.Ids,
			},
		}); err != nil {
			log2.Warnf("delete task stat error: %s", err.Error())
			return nil
		}

		return nil
	}); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// delete tasks logs
	wg := sync.WaitGroup{}
	wg.Add(len(payload.Ids))
	for _, id := range payload.Ids {
		go func(id string) {
			// delete task logs
			logPath := filepath.Join(viper.GetString("log.path"), id)
			if err := os.RemoveAll(logPath); err != nil {
				log2.Warnf("failed to remove task log directory: %s", logPath)
			}
			wg.Done()
		}(id.Hex())
	}
	wg.Wait()

	HandleSuccess(c)
}

func PostTaskRun(c *gin.Context) {
	// task
	var t models.TaskV2
	if err := c.ShouldBindJSON(&t); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// validate spider id
	if t.SpiderId.IsZero() {
		HandleErrorBadRequest(c, errors.New("spider id is required"))
		return
	}

	// spider
	s, err := service.NewModelServiceV2[models.SpiderV2]().GetById(t.SpiderId)
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
	if u := GetUserFromContextV2(c); u != nil {
		opts.UserId = u.Id
	}

	// run
	adminSvc, err := admin.GetSpiderAdminServiceV2()
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	taskIds, err := adminSvc.Schedule(s.Id, opts)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccessWithData(c, taskIds)

}

func PostTaskRestart(c *gin.Context) {
	// id
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// task
	t, err := service.NewModelServiceV2[models.TaskV2]().GetById(id)
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
	if u := GetUserFromContextV2(c); u != nil {
		opts.UserId = u.Id
	}

	// run
	adminSvc, err := admin.GetSpiderAdminServiceV2()
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	taskIds, err := adminSvc.Schedule(t.SpiderId, opts)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccessWithData(c, taskIds)
}

func PostTaskCancel(c *gin.Context) {
	// id
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// task
	t, err := service.NewModelServiceV2[models.TaskV2]().GetById(id)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// validate
	if !utils.IsCancellable(t.Status) {
		HandleErrorInternalServerError(c, errors.New("task is not cancellable"))
		return
	}

	u := GetUserFromContextV2(c)

	// cancel
	schedulerSvc, err := scheduler.GetTaskSchedulerServiceV2()
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	if err := schedulerSvc.Cancel(id, u.Id); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccess(c)
}

func GetTaskLogs(c *gin.Context) {
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
	logDriver, err := log.GetFileLogDriver()
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	logs, err := logDriver.Find(id.Hex(), "", (p.Page-1)*p.Size, p.Size)
	if err != nil {
		if strings.HasSuffix(err.Error(), "Status:404 Not Found") {
			HandleSuccess(c)
			return
		}
		HandleErrorInternalServerError(c, err)
		return
	}
	total, err := logDriver.Count(id.Hex(), "")
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccessWithListData(c, logs, total)
}

func GetTaskData(c *gin.Context) {
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
	t, err := service.NewModelServiceV2[models.TaskV2]().GetById(id)
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
