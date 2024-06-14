package controllers

import (
	"bytes"
	"fmt"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/container"
	"github.com/crawlab-team/crawlab/core/entity"
	"github.com/crawlab-team/crawlab/core/errors"
	fs2 "github.com/crawlab-team/crawlab/core/fs"
	"github.com/crawlab-team/crawlab/core/interfaces"
	delegate2 "github.com/crawlab-team/crawlab/core/models/delegate"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/core/utils"
	"github.com/crawlab-team/crawlab/db/mongo"
	"github.com/crawlab-team/crawlab/trace"
	"github.com/crawlab-team/crawlab/vcs"
	"github.com/gin-gonic/gin"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var SpiderController *spiderController

func getSpiderActions() []Action {
	ctx := newSpiderContext()
	return []Action{
		{
			Method:      http.MethodGet,
			Path:        "/:id/files/list",
			HandlerFunc: ctx.listDir,
		},
		{
			Method:      http.MethodGet,
			Path:        "/:id/files/get",
			HandlerFunc: ctx.getFile,
		},
		{
			Method:      http.MethodGet,
			Path:        "/:id/files/info",
			HandlerFunc: ctx.getFileInfo,
		},
		{
			Method:      http.MethodPost,
			Path:        "/:id/files/save",
			HandlerFunc: ctx.saveFile,
		},
		{
			Method:      http.MethodPost,
			Path:        "/:id/files/save/dir",
			HandlerFunc: ctx.saveDir,
		},
		{
			Method:      http.MethodPost,
			Path:        "/:id/files/rename",
			HandlerFunc: ctx.renameFile,
		},
		{
			Method:      http.MethodPost,
			Path:        "/:id/files/delete",
			HandlerFunc: ctx.deleteFile,
		},
		{
			Method:      http.MethodPost,
			Path:        "/:id/files/copy",
			HandlerFunc: ctx.copyFile,
		},
		{
			Path:        "/:id/files/export",
			Method:      http.MethodPost,
			HandlerFunc: ctx.postExport,
		},
		{
			Method:      http.MethodPost,
			Path:        "/:id/run",
			HandlerFunc: ctx.run,
		},
		{
			Method:      http.MethodGet,
			Path:        "/:id/git",
			HandlerFunc: ctx.getGit,
		},
		{
			Method:      http.MethodGet,
			Path:        "/:id/git/remote-refs",
			HandlerFunc: ctx.getGitRemoteRefs,
		},
		{
			Method:      http.MethodPost,
			Path:        "/:id/git/checkout",
			HandlerFunc: ctx.gitCheckout,
		},
		{
			Method:      http.MethodPost,
			Path:        "/:id/git/pull",
			HandlerFunc: ctx.gitPull,
		},
		{
			Method:      http.MethodPost,
			Path:        "/:id/git/commit",
			HandlerFunc: ctx.gitCommit,
		},
		//{
		//	Method:      http.MethodPost,
		//	Path:        "/:id/clone",
		//	HandlerFunc: ctx.clone,
		//},
		{
			Path:        "/:id/data-source",
			Method:      http.MethodGet,
			HandlerFunc: ctx.getDataSource,
		},
		{
			Path:        "/:id/data-source/:ds_id",
			Method:      http.MethodPost,
			HandlerFunc: ctx.postDataSource,
		},
	}
}

type spiderController struct {
	ListActionControllerDelegate
	d   ListActionControllerDelegate
	ctx *spiderContext
}

func (ctr *spiderController) Get(c *gin.Context) {
	ctr.ctx._get(c)
}

func (ctr *spiderController) Post(c *gin.Context) {
	s, err := ctr.ctx._post(c)
	if err != nil {
		return
	}
	HandleSuccessWithData(c, s)
}

func (ctr *spiderController) Put(c *gin.Context) {
	s, err := ctr.ctx._put(c)
	if err != nil {
		return
	}
	HandleSuccessWithData(c, s)
}

func (ctr *spiderController) Delete(c *gin.Context) {
	if err := ctr.ctx._delete(c); err != nil {
		return
	}
	HandleSuccess(c)
}

func (ctr *spiderController) GetList(c *gin.Context) {
	withStats := c.Query("stats")
	if withStats == "" {
		ctr.d.GetList(c)
		return
	}
	ctr.ctx._getListWithStats(c)
}

func (ctr *spiderController) DeleteList(c *gin.Context) {
	if err := ctr.ctx._deleteList(c); err != nil {
		return
	}
	HandleSuccess(c)
}

type spiderContext struct {
	modelSvc           service.ModelService
	modelSpiderSvc     interfaces.ModelBaseService
	modelSpiderStatSvc interfaces.ModelBaseService
	modelTaskSvc       interfaces.ModelBaseService
	modelTaskStatSvc   interfaces.ModelBaseService
	adminSvc           interfaces.SpiderAdminService
}

func (ctx *spiderContext) listDir(c *gin.Context) {
	_, payload, fsSvc, err := ctx._processFileRequest(c, http.MethodGet)
	if err != nil {
		return
	}

	files, err := fsSvc.List(payload.Path)
	if err != nil {
		if err.Error() != "response status code: 404" {
			HandleErrorInternalServerError(c, err)
			return
		}
	}

	HandleSuccessWithData(c, files)
}

func (ctx *spiderContext) getFile(c *gin.Context) {
	_, payload, fsSvc, err := ctx._processFileRequest(c, http.MethodGet)
	if err != nil {
		return
	}

	data, err := fsSvc.GetFile(payload.Path)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	data = utils.TrimFileData(data)

	HandleSuccessWithData(c, string(data))
}

func (ctx *spiderContext) getFileInfo(c *gin.Context) {
	_, payload, fsSvc, err := ctx._processFileRequest(c, http.MethodGet)
	if err != nil {
		return
	}

	info, err := fsSvc.GetFileInfo(payload.Path)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccessWithData(c, info)
}

func (ctx *spiderContext) saveFile(c *gin.Context) {
	_, payload, fsSvc, err := ctx._processFileRequest(c, http.MethodPost)
	if err != nil {
		return
	}

	if err := fsSvc.Save(payload.Path, []byte(payload.Data)); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccess(c)
}

func (ctx *spiderContext) saveDir(c *gin.Context) {
	_, payload, fsSvc, err := ctx._processFileRequest(c, http.MethodPost)
	if err != nil {
		return
	}

	if err := fsSvc.CreateDir(payload.Path); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccess(c)
}

func (ctx *spiderContext) renameFile(c *gin.Context) {
	_, payload, fsSvc, err := ctx._processFileRequest(c, http.MethodPost)
	if err != nil {
		return
	}

	if err := fsSvc.Rename(payload.Path, payload.NewPath); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccess(c)
}

func (ctx *spiderContext) deleteFile(c *gin.Context) {
	_, payload, fsSvc, err := ctx._processFileRequest(c, http.MethodPost)
	if err != nil {
		return
	}

	if err := fsSvc.Delete(payload.Path); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccess(c)
}

func (ctx *spiderContext) copyFile(c *gin.Context) {
	_, payload, fsSvc, err := ctx._processFileRequest(c, http.MethodPost)
	if err != nil {
		return
	}

	if err := fsSvc.Copy(payload.Path, payload.NewPath); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccess(c)
}

func (ctx *spiderContext) run(c *gin.Context) {
	// spider id
	id, err := ctx._processActionRequest(c)
	if err != nil {
		return
	}

	// options
	var opts interfaces.SpiderRunOptions
	if err := c.ShouldBindJSON(&opts); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// user
	if u := GetUserFromContext(c); u != nil {
		opts.UserId = u.GetId()
	}

	// schedule
	taskIds, err := ctx.adminSvc.Schedule(id, &opts)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccessWithData(c, taskIds)
}

func (ctx *spiderContext) getGit(c *gin.Context) {
	// spider id
	id, err := ctx._processActionRequest(c)
	if err != nil {
		return
	}

	// git client
	gitClient, err := ctx._getGitClient(id)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// return null if git client is empty
	if gitClient == nil {
		HandleSuccess(c)
		return
	}

	// current branch
	currentBranch, err := gitClient.GetCurrentBranch()
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// branches
	branches, err := gitClient.GetBranches()
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	if branches == nil || len(branches) == 0 && currentBranch != "" {
		branches = []vcs.GitRef{{Name: currentBranch}}
	}

	// changes
	changes, err := gitClient.GetStatus()
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// logs
	logs, err := gitClient.GetLogsWithRefs()
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// ignore
	ignore, err := ctx._getGitIgnore(id)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// git
	_git, err := ctx.modelSvc.GetGitById(id)
	if err != nil {
		if err.Error() != mongo2.ErrNoDocuments.Error() {
			HandleErrorInternalServerError(c, err)
			return
		}
	}

	// response
	res := bson.M{
		"current_branch": currentBranch,
		"branches":       branches,
		"changes":        changes,
		"logs":           logs,
		"ignore":         ignore,
		"git":            _git,
	}

	HandleSuccessWithData(c, res)
}

func (ctx *spiderContext) getGitRemoteRefs(c *gin.Context) {
	// spider id
	id, err := ctx._processActionRequest(c)
	if err != nil {
		return
	}

	// remote name
	remoteName := c.Query("remote")
	if remoteName == "" {
		remoteName = vcs.GitRemoteNameOrigin
	}

	// git client
	gitClient, err := ctx._getGitClient(id)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// return null if git client is empty
	if gitClient == nil {
		HandleSuccess(c)
		return
	}

	// refs
	refs, err := gitClient.GetRemoteRefs(remoteName)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccessWithData(c, refs)
}

func (ctx *spiderContext) gitCheckout(c *gin.Context) {
	// payload
	var payload entity.GitPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// spider id
	id, err := ctx._processActionRequest(c)
	if err != nil {
		return
	}

	// git client
	gitClient, err := ctx._getGitClient(id)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// return null if git client is empty
	if gitClient == nil {
		HandleSuccess(c)
		return
	}

	// branch to pull
	var branch string
	if payload.Branch == "" {
		// by default current branch
		branch, err = gitClient.GetCurrentBranch()
		if err != nil {
			HandleErrorInternalServerError(c, err)
			return
		}
	} else {
		// payload branch
		branch = payload.Branch
	}

	// checkout
	if err := ctx._gitCheckout(gitClient, constants.GitRemoteNameOrigin, branch); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccess(c)
}

func (ctx *spiderContext) gitPull(c *gin.Context) {
	// payload
	var payload entity.GitPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// spider id
	id, err := ctx._processActionRequest(c)
	if err != nil {
		return
	}

	// git
	g, err := ctx.modelSvc.GetGitById(id)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// attempt to sync git
	_ = ctx.adminSvc.SyncGitOne(g)

	HandleSuccess(c)
}

func (ctx *spiderContext) gitCommit(c *gin.Context) {
	// payload
	var payload entity.GitPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// spider id
	id, err := ctx._processActionRequest(c)
	if err != nil {
		return
	}

	// git client
	gitClient, err := ctx._getGitClient(id)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// return null if git client is empty
	if gitClient == nil {
		HandleSuccess(c)
		return
	}

	// add
	for _, p := range payload.Paths {
		if err := gitClient.Add(p); err != nil {
			HandleErrorInternalServerError(c, err)
			return
		}
	}

	// commit
	if err := gitClient.Commit(payload.CommitMessage); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// push
	if err := gitClient.Push(
		vcs.WithRemoteNamePush(vcs.GitRemoteNameOrigin),
	); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccess(c)
}

func (ctx *spiderContext) getDataSource(c *gin.Context) {
	// spider id
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// spider
	s, err := ctx.modelSvc.GetSpiderById(id)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// data source
	ds, err := ctx.modelSvc.GetDataSourceById(s.DataSourceId)
	if err != nil {
		if err.Error() == mongo2.ErrNoDocuments.Error() {
			HandleSuccess(c)
			return
		}
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccessWithData(c, ds)
}

func (ctx *spiderContext) postDataSource(c *gin.Context) {
	// spider id
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// data source id
	dsId, err := primitive.ObjectIDFromHex(c.Param("ds_id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// spider
	s, err := ctx.modelSvc.GetSpiderById(id)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// data source
	if !dsId.IsZero() {
		_, err = ctx.modelSvc.GetDataSourceById(dsId)
		if err != nil {
			HandleErrorInternalServerError(c, err)
			return
		}
	}

	// save data source id
	s.DataSourceId = dsId
	if err := delegate2.NewModelDelegate(s).Save(); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccess(c)
}

func (ctx *spiderContext) postExport(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// zip file path
	zipFilePath, err := ctx.adminSvc.Export(id)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// download
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", zipFilePath))
	c.File(zipFilePath)
}

func (ctx *spiderContext) _get(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	s, err := ctx.modelSvc.GetSpiderById(id)
	if err == mongo2.ErrNoDocuments {
		HandleErrorNotFound(c, err)
		return
	}
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// stat
	s.Stat, err = ctx.modelSvc.GetSpiderStatById(s.GetId())
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// data collection
	if !s.ColId.IsZero() {
		col, err := ctx.modelSvc.GetDataCollectionById(s.ColId)
		if err != nil {
			if err != mongo2.ErrNoDocuments {
				HandleErrorInternalServerError(c, err)
				return
			}
		} else {
			s.ColName = col.Name
		}
	}

	HandleSuccessWithData(c, s)
}

func (ctx *spiderContext) _post(c *gin.Context) (s *models.Spider, err error) {
	// bind
	s = &models.Spider{}
	if err := c.ShouldBindJSON(&s); err != nil {
		HandleErrorBadRequest(c, err)
		return nil, err
	}

	// upsert data collection
	if err := ctx._upsertDataCollection(c, s); err != nil {
		HandleErrorInternalServerError(c, err)
		return nil, err
	}

	// add
	if err := delegate2.NewModelDelegate(s, GetUserFromContext(c)).Add(); err != nil {
		HandleErrorInternalServerError(c, err)
		return nil, err
	}

	// add stat
	st := &models.SpiderStat{
		Id: s.GetId(),
	}
	if err := delegate2.NewModelDelegate(st, GetUserFromContext(c)).Add(); err != nil {
		HandleErrorInternalServerError(c, err)
		return nil, err
	}

	return s, nil
}

func (ctx *spiderContext) _put(c *gin.Context) (s *models.Spider, err error) {
	// bind
	s = &models.Spider{}
	if err := c.ShouldBindJSON(&s); err != nil {
		HandleErrorBadRequest(c, err)
		return nil, err
	}

	// upsert data collection
	if err := ctx._upsertDataCollection(c, s); err != nil {
		HandleErrorInternalServerError(c, err)
		return nil, err
	}

	// save
	if err := delegate2.NewModelDelegate(s, GetUserFromContext(c)).Save(); err != nil {
		HandleErrorInternalServerError(c, err)
		return nil, err
	}

	return s, nil
}

func (ctx *spiderContext) _delete(c *gin.Context) (err error) {
	id := c.Param("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	if err := mongo.RunTransaction(func(context mongo2.SessionContext) (err error) {
		// delete spider
		s, err := ctx.modelSvc.GetSpiderById(oid)
		if err != nil {
			return err
		}
		if err := delegate2.NewModelDelegate(s, GetUserFromContext(c)).Delete(); err != nil {
			return err
		}

		// delete spider stat
		ss, err := ctx.modelSvc.GetSpiderStatById(oid)
		if err != nil {
			return err
		}
		if err := delegate2.NewModelDelegate(ss, GetUserFromContext(c)).Delete(); err != nil {
			return err
		}

		// related tasks
		tasks, err := ctx.modelSvc.GetTaskList(bson.M{"spider_id": oid}, nil)
		if err != nil {
			return err
		}

		// task ids
		var taskIds []primitive.ObjectID
		for _, t := range tasks {
			taskIds = append(taskIds, t.Id)
		}

		// delete related tasks
		if err := ctx.modelTaskSvc.DeleteList(bson.M{"_id": bson.M{"$in": taskIds}}); err != nil {
			return err
		}

		// delete related task stats
		if err := ctx.modelTaskStatSvc.DeleteList(bson.M{"_id": bson.M{"$in": taskIds}}); err != nil {
			return err
		}

		return nil
	}); err != nil {
		HandleErrorInternalServerError(c, err)
		return err
	}

	return nil
}

func (ctx *spiderContext) _getListWithStats(c *gin.Context) {
	// params
	pagination := MustGetPagination(c)
	query := MustGetFilterQuery(c)
	sort := MustGetSortOption(c)

	// get list
	l, err := ctx.modelSpiderSvc.GetList(query, &mongo.FindOptions{
		Sort:  sort,
		Skip:  pagination.Size * (pagination.Page - 1),
		Limit: pagination.Size,
	})
	if err != nil {
		if err.Error() == mongo2.ErrNoDocuments.Error() {
			HandleErrorNotFound(c, err)
		} else {
			HandleErrorInternalServerError(c, err)
		}
		return
	}

	// check empty list
	if len(l.GetModels()) == 0 {
		HandleSuccessWithListData(c, nil, 0)
		return
	}

	// ids
	var ids []primitive.ObjectID
	for _, d := range l.GetModels() {
		s := d.(*models.Spider)
		ids = append(ids, s.GetId())
	}

	// total count
	total, err := ctx.modelSpiderSvc.Count(query)
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
	stats, err := ctx.modelSvc.GetSpiderStatList(query, nil)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// cache stat list to dict
	dict := map[primitive.ObjectID]models.SpiderStat{}
	var tids []primitive.ObjectID
	for _, st := range stats {
		if st.Tasks > 0 {
			taskCount := int64(st.Tasks)
			st.AverageWaitDuration = int64(math.Round(float64(st.WaitDuration) / float64(taskCount)))
			st.AverageRuntimeDuration = int64(math.Round(float64(st.RuntimeDuration) / float64(taskCount)))
			st.AverageTotalDuration = int64(math.Round(float64(st.TotalDuration) / float64(taskCount)))
		}
		dict[st.GetId()] = st

		if !st.LastTaskId.IsZero() {
			tids = append(tids, st.LastTaskId)
		}
	}

	// task list and stats
	var tasks []models.Task
	dictTask := map[primitive.ObjectID]models.Task{}
	dictTaskStat := map[primitive.ObjectID]models.TaskStat{}
	if len(tids) > 0 {
		// task list
		queryTask := bson.M{
			"_id": bson.M{
				"$in": tids,
			},
		}
		tasks, err = ctx.modelSvc.GetTaskList(queryTask, nil)
		if err != nil {
			HandleErrorInternalServerError(c, err)
			return
		}

		// task stats list
		taskStats, err := ctx.modelSvc.GetTaskStatList(queryTask, nil)
		if err != nil {
			HandleErrorInternalServerError(c, err)
			return
		}

		// cache task stats to dict
		for _, st := range taskStats {
			dictTaskStat[st.GetId()] = st
		}

		// cache task list to dict
		for _, t := range tasks {
			st, ok := dictTaskStat[t.GetId()]
			if ok {
				t.Stat = &st
			}
			dictTask[t.GetSpiderId()] = t
		}
	}

	// iterate list again
	var data []interface{}
	for _, d := range l.GetModels() {
		s := d.(*models.Spider)

		// spider stat
		st, ok := dict[s.GetId()]
		if ok {
			s.Stat = &st

			// last task
			t, ok := dictTask[s.GetId()]
			if ok {
				s.Stat.LastTask = &t
			}
		}

		// add to list
		data = append(data, *s)
	}

	// response
	HandleSuccessWithListData(c, data, total)
}

func (ctx *spiderContext) _deleteList(c *gin.Context) (err error) {
	payload, err := NewJsonBinder(ControllerIdSpider).BindBatchRequestPayload(c)
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	if err := mongo.RunTransaction(func(context mongo2.SessionContext) (err error) {
		// delete spiders
		if err := ctx.modelSpiderSvc.DeleteList(bson.M{
			"_id": bson.M{
				"$in": payload.Ids,
			},
		}); err != nil {
			return err
		}

		// delete spider stats
		if err := ctx.modelSpiderStatSvc.DeleteList(bson.M{
			"_id": bson.M{
				"$in": payload.Ids,
			},
		}); err != nil {
			return err
		}

		// related tasks
		tasks, err := ctx.modelSvc.GetTaskList(bson.M{"spider_id": bson.M{"$in": payload.Ids}}, nil)
		if err != nil {
			return err
		}

		// task ids
		var taskIds []primitive.ObjectID
		for _, t := range tasks {
			taskIds = append(taskIds, t.Id)
		}

		// delete related tasks
		if err := ctx.modelTaskSvc.DeleteList(bson.M{"_id": bson.M{"$in": taskIds}}); err != nil {
			return err
		}

		// delete related task stats
		if err := ctx.modelTaskStatSvc.DeleteList(bson.M{"_id": bson.M{"$in": taskIds}}); err != nil {
			return err
		}

		return nil
	}); err != nil {
		HandleErrorInternalServerError(c, err)
		return err
	}

	return nil
}

func (ctx *spiderContext) _processFileRequest(c *gin.Context, method string) (id primitive.ObjectID, payload entity.FileRequestPayload, fsSvc interfaces.FsServiceV2, err error) {
	// id
	id, err = primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// payload
	contentType := c.GetHeader("Content-Type")
	if strings.HasPrefix(contentType, "multipart/form-data") {
		// multipart/form-data
		payload, err = ctx._getFileRequestMultipartPayload(c)
		if err != nil {
			HandleErrorBadRequest(c, err)
			return
		}
	} else {
		// query or application/json
		switch method {
		case http.MethodGet:
			err = c.ShouldBindQuery(&payload)
		default:
			err = c.ShouldBindJSON(&payload)
		}
		if err != nil {
			HandleErrorInternalServerError(c, err)
			return
		}
	}

	// fs service
	workspacePath := viper.GetString("workspace")
	fsSvc = fs2.NewFsServiceV2(filepath.Join(workspacePath, id.Hex()))

	return
}

func (ctx *spiderContext) _getFileRequestMultipartPayload(c *gin.Context) (payload entity.FileRequestPayload, err error) {
	fh, err := c.FormFile("file")
	if err != nil {
		return
	}
	f, err := fh.Open()
	if err != nil {
		return
	}
	buf := bytes.NewBuffer(nil)
	if _, err = io.Copy(buf, f); err != nil {
		return
	}
	payload.Path = c.PostForm("path")
	payload.Data = buf.String()
	return
}

func (ctx *spiderContext) _processActionRequest(c *gin.Context) (id primitive.ObjectID, err error) {
	// id
	id, err = primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	return
}

func (ctx *spiderContext) _upsertDataCollection(c *gin.Context, s *models.Spider) (err error) {
	if s.ColId.IsZero() {
		// validate
		if s.ColName == "" {
			return trace.TraceError(errors.ErrorControllerMissingRequestFields)
		}
		// no id
		dc, err := ctx.modelSvc.GetDataCollectionByName(s.ColName, nil)
		if err != nil {
			if err == mongo2.ErrNoDocuments {
				// not exists, add new
				dc = &models.DataCollection{Name: s.ColName}
				if err := delegate2.NewModelDelegate(dc, GetUserFromContext(c)).Add(); err != nil {
					return err
				}
			} else {
				// error
				return err
			}
		}
		s.ColId = dc.Id

		// create index
		_ = mongo.GetMongoCol(dc.Name).CreateIndex(mongo2.IndexModel{Keys: bson.M{constants.TaskKey: 1}})
		_ = mongo.GetMongoCol(dc.Name).CreateIndex(mongo2.IndexModel{Keys: bson.M{constants.HashKey: 1}})
	} else {
		// with id
		dc, err := ctx.modelSvc.GetDataCollectionById(s.ColId)
		if err != nil {
			return err
		}
		s.ColId = dc.Id
	}
	return nil
}

func (ctx *spiderContext) _getGitIgnore(id primitive.ObjectID) (ignore []string, err error) {
	workspacePath := viper.GetString("workspace")
	filePath := filepath.Join(workspacePath, id.Hex(), ".gitignore")
	if !utils.Exists(filePath) {
		return nil, nil
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, trace.TraceError(err)
	}
	ignore = strings.Split(string(data), "\n")
	return ignore, nil
}

func (ctx *spiderContext) _gitCheckout(gitClient *vcs.GitClient, remote, branch string) (err error) {
	if err := gitClient.CheckoutBranch(branch, vcs.WithBranch(branch)); err != nil {
		return trace.TraceError(err)
	}

	// pull
	return ctx._gitPull(gitClient, remote, branch)
}

func (ctx *spiderContext) _gitPull(gitClient *vcs.GitClient, remote, branch string) (err error) {
	// pull
	if err := gitClient.Pull(
		vcs.WithRemoteNamePull(remote),
		vcs.WithBranchNamePull(branch),
	); err != nil {
		return trace.TraceError(err)
	}

	// reset
	if err := gitClient.Reset(); err != nil {
		return trace.TraceError(err)
	}

	return nil
}

func (ctx *spiderContext) _getGitClient(id primitive.ObjectID) (gitClient *vcs.GitClient, err error) {
	// git
	g, err := ctx.modelSvc.GetGitById(id)
	if err != nil {
		if err != mongo2.ErrNoDocuments {
			return nil, trace.TraceError(err)
		}
		return nil, nil
	}

	// git client
	workspacePath := viper.GetString("workspace")
	gitClient, err = vcs.NewGitClient(vcs.WithPath(filepath.Join(workspacePath, id.Hex())))
	if err != nil {
		return nil, err
	}

	// set auth
	utils.InitGitClientAuth(g, gitClient)

	// remote name
	remoteName := vcs.GitRemoteNameOrigin

	// update remote
	r, err := gitClient.GetRemote(remoteName)
	if err == git.ErrRemoteNotFound {
		// remote not exists, create
		if _, err := gitClient.CreateRemote(&config.RemoteConfig{
			Name: remoteName,
			URLs: []string{g.Url},
		}); err != nil {
			return nil, trace.TraceError(err)
		}
	} else if err == nil {
		// remote exists, update if different
		if g.Url != r.Config().URLs[0] {
			if err := gitClient.DeleteRemote(remoteName); err != nil {
				return nil, trace.TraceError(err)
			}
			if _, err := gitClient.CreateRemote(&config.RemoteConfig{
				Name: remoteName,
				URLs: []string{g.Url},
			}); err != nil {
				return nil, trace.TraceError(err)
			}
		}
		gitClient.SetRemoteUrl(g.Url)
	} else {
		// error
		return nil, trace.TraceError(err)
	}

	// check if head reference exists
	_, err = gitClient.GetRepository().Head()
	if err == nil {
		return gitClient, nil
	}

	// align master/main branch
	ctx._alignBranch(gitClient)

	return gitClient, nil
}

func (ctx *spiderContext) _alignBranch(gitClient *vcs.GitClient) {
	// current branch
	currentBranch, err := gitClient.GetCurrentBranch()
	if err != nil {
		trace.PrintError(err)
		return
	}

	// skip if current branch is not master
	if currentBranch != vcs.GitBranchNameMaster {
		return
	}

	// remote refs
	refs, err := gitClient.GetRemoteRefs(vcs.GitRemoteNameOrigin)
	if err != nil {
		trace.PrintError(err)
		return
	}

	// main branch
	defaultRemoteBranch, err := ctx._getDefaultRemoteBranch(refs)
	if err != nil || defaultRemoteBranch == "" {
		return
	}

	// move branch
	if err := gitClient.MoveBranch(vcs.GitBranchNameMaster, defaultRemoteBranch); err != nil {
		trace.PrintError(err)
	}
}

func (ctx *spiderContext) _getDefaultRemoteBranch(refs []vcs.GitRef) (defaultRemoteBranchName string, err error) {
	// remote branch name
	for _, r := range refs {
		if r.Type != vcs.GitRefTypeBranch {
			continue
		}

		if r.Name == vcs.GitBranchNameMain {
			defaultRemoteBranchName = r.Name
			break
		}

		if r.Name == vcs.GitBranchNameMaster {
			defaultRemoteBranchName = r.Name
			continue
		}

		if defaultRemoteBranchName == "" {
			defaultRemoteBranchName = r.Name
			continue
		}
	}

	return defaultRemoteBranchName, nil
}

var _spiderCtx *spiderContext

func newSpiderContext() *spiderContext {
	if _spiderCtx != nil {
		return _spiderCtx
	}

	// context
	ctx := &spiderContext{}

	// dependency injection
	if err := container.GetContainer().Invoke(func(
		modelSvc service.ModelService,
		adminSvc interfaces.SpiderAdminService,
	) {
		ctx.modelSvc = modelSvc
		ctx.adminSvc = adminSvc
	}); err != nil {
		panic(err)
	}

	// model spider service
	ctx.modelSpiderSvc = ctx.modelSvc.GetBaseService(interfaces.ModelIdSpider)

	// model spider stat service
	ctx.modelSpiderStatSvc = ctx.modelSvc.GetBaseService(interfaces.ModelIdSpiderStat)

	// model task service
	ctx.modelTaskSvc = ctx.modelSvc.GetBaseService(interfaces.ModelIdTask)

	// model task stat service
	ctx.modelTaskStatSvc = ctx.modelSvc.GetBaseService(interfaces.ModelIdTaskStat)

	_spiderCtx = ctx

	return ctx
}

func newSpiderController() *spiderController {
	actions := getSpiderActions()
	modelSvc, err := service.GetService()
	if err != nil {
		panic(err)
	}

	ctr := NewListPostActionControllerDelegate(ControllerIdSpider, modelSvc.GetBaseService(interfaces.ModelIdSpider), actions)
	d := NewListPostActionControllerDelegate(ControllerIdSpider, modelSvc.GetBaseService(interfaces.ModelIdSpider), actions)
	ctx := newSpiderContext()

	return &spiderController{
		ListActionControllerDelegate: *ctr,
		d:                            *d,
		ctx:                          ctx,
	}
}
