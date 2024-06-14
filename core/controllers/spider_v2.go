package controllers

import (
	"errors"
	"fmt"
	log2 "github.com/apex/log"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/entity"
	"github.com/crawlab-team/crawlab/core/fs"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/core/spider/admin"
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
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func GetSpiderById(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	s, err := service.NewModelServiceV2[models.SpiderV2]().GetById(id)
	if errors.Is(err, mongo2.ErrNoDocuments) {
		HandleErrorNotFound(c, err)
		return
	}
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// stat
	s.Stat, err = service.NewModelServiceV2[models.SpiderStatV2]().GetById(s.Id)
	if err != nil {
		if !errors.Is(err, mongo2.ErrNoDocuments) {
			HandleErrorInternalServerError(c, err)
			return
		}
	}

	// data collection
	if !s.ColId.IsZero() {
		col, err := service.NewModelServiceV2[models.DataCollectionV2]().GetById(s.ColId)
		if err != nil {
			if !errors.Is(err, mongo2.ErrNoDocuments) {
				HandleErrorInternalServerError(c, err)
				return
			}
		} else {
			s.ColName = col.Name
		}
	}

	HandleSuccessWithData(c, s)
}

func GetSpiderList(c *gin.Context) {
	withStats := c.Query("stats")
	if withStats == "" {
		NewControllerV2[models.SpiderV2]().GetList(c)
		return
	}

	// params
	pagination := MustGetPagination(c)
	query := MustGetFilterQuery(c)
	sort := MustGetSortOption(c)

	// get list
	spiders, err := service.NewModelServiceV2[models.SpiderV2]().GetMany(query, &mongo.FindOptions{
		Sort:  sort,
		Skip:  pagination.Size * (pagination.Page - 1),
		Limit: pagination.Size,
	})
	if err != nil {
		if err.Error() != mongo2.ErrNoDocuments.Error() {
			HandleErrorInternalServerError(c, err)
		}
		return
	}
	if len(spiders) == 0 {
		HandleSuccessWithListData(c, []models.SpiderV2{}, 0)
		return
	}

	// ids
	var ids []primitive.ObjectID
	for _, s := range spiders {
		ids = append(ids, s.Id)
	}

	// total count
	total, err := service.NewModelServiceV2[models.SpiderV2]().Count(query)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// stat list
	spiderStats, err := service.NewModelServiceV2[models.SpiderStatV2]().GetMany(bson.M{"_id": bson.M{"$in": ids}}, nil)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// cache stat list to dict
	dict := map[primitive.ObjectID]models.SpiderStatV2{}
	var taskIds []primitive.ObjectID
	for _, st := range spiderStats {
		if st.Tasks > 0 {
			taskCount := int64(st.Tasks)
			st.AverageWaitDuration = int64(math.Round(float64(st.WaitDuration) / float64(taskCount)))
			st.AverageRuntimeDuration = int64(math.Round(float64(st.RuntimeDuration) / float64(taskCount)))
			st.AverageTotalDuration = int64(math.Round(float64(st.TotalDuration) / float64(taskCount)))
		}
		dict[st.Id] = st

		if !st.LastTaskId.IsZero() {
			taskIds = append(taskIds, st.LastTaskId)
		}
	}

	// task list and stats
	var tasks []models.TaskV2
	dictTask := map[primitive.ObjectID]models.TaskV2{}
	dictTaskStat := map[primitive.ObjectID]models.TaskStatV2{}
	if len(taskIds) > 0 {
		// task list
		queryTask := bson.M{
			"_id": bson.M{
				"$in": taskIds,
			},
		}
		tasks, err = service.NewModelServiceV2[models.TaskV2]().GetMany(queryTask, nil)
		if err != nil {
			HandleErrorInternalServerError(c, err)
			return
		}

		// task stats list
		taskStats, err := service.NewModelServiceV2[models.TaskStatV2]().GetMany(queryTask, nil)
		if err != nil {
			HandleErrorInternalServerError(c, err)
			return
		}

		// cache task stats to dict
		for _, st := range taskStats {
			dictTaskStat[st.Id] = st
		}

		// cache task list to dict
		for _, t := range tasks {
			st, ok := dictTaskStat[t.Id]
			if ok {
				t.Stat = &st
			}
			dictTask[t.SpiderId] = t
		}
	}

	// iterate list again
	var data []models.SpiderV2
	for _, s := range spiders {
		// spider stat
		st, ok := dict[s.Id]
		if ok {
			s.Stat = &st

			// last task
			t, ok := dictTask[s.Id]
			if ok {
				s.Stat.LastTask = &t
			}
		}

		// add to list
		data = append(data, s)
	}

	// response
	HandleSuccessWithListData(c, data, total)
}

func PostSpider(c *gin.Context) {
	// bind
	var s models.SpiderV2
	if err := c.ShouldBindJSON(&s); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// upsert data collection
	if err := upsertSpiderDataCollection(&s); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	u := GetUserFromContextV2(c)

	// add
	s.SetCreated(u.Id)
	s.SetUpdated(u.Id)
	id, err := service.NewModelServiceV2[models.SpiderV2]().InsertOne(s)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	s.SetId(id)

	// add stat
	st := models.SpiderStatV2{}
	st.SetId(id)
	st.SetCreated(u.Id)
	st.SetUpdated(u.Id)
	_, err = service.NewModelServiceV2[models.SpiderStatV2]().InsertOne(st)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// create folder
	err = getSpiderFsSvcById(id).CreateDir(".")
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccessWithData(c, s)
}

func PutSpiderById(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// bind
	var s models.SpiderV2
	if err := c.ShouldBindJSON(&s); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// upsert data collection
	if err := upsertSpiderDataCollection(&s); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	u := GetUserFromContextV2(c)

	modelSvc := service.NewModelServiceV2[models.SpiderV2]()

	// save
	s.SetUpdated(u.Id)
	err = modelSvc.ReplaceById(id, s)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	_s, err := modelSvc.GetById(id)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	s = *_s

	HandleSuccessWithData(c, s)
}

func DeleteSpiderById(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	if err := mongo.RunTransaction(func(context mongo2.SessionContext) (err error) {
		// delete spider
		err = service.NewModelServiceV2[models.SpiderV2]().DeleteById(id)
		if err != nil {
			return err
		}

		// delete spider stat
		err = service.NewModelServiceV2[models.SpiderStatV2]().DeleteById(id)
		if err != nil {
			return err
		}

		// related tasks
		tasks, err := service.NewModelServiceV2[models.TaskV2]().GetMany(bson.M{"spider_id": id}, nil)
		if err != nil {
			return err
		}

		if len(tasks) == 0 {
			return nil
		}

		// task ids
		var taskIds []primitive.ObjectID
		for _, t := range tasks {
			taskIds = append(taskIds, t.Id)
		}

		// delete related tasks
		err = service.NewModelServiceV2[models.TaskV2]().DeleteMany(bson.M{"_id": bson.M{"$in": taskIds}})
		if err != nil {
			return err
		}

		// delete related task stats
		err = service.NewModelServiceV2[models.TaskStatV2]().DeleteMany(bson.M{"_id": bson.M{"$in": taskIds}})
		if err != nil {
			return err
		}

		// delete tasks logs
		wg := sync.WaitGroup{}
		wg.Add(len(taskIds))
		for _, id := range taskIds {
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

		return nil
	}); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccess(c)
}

func DeleteSpiderList(c *gin.Context) {
	var payload struct {
		Ids []primitive.ObjectID `json:"ids"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	if err := mongo.RunTransaction(func(context mongo2.SessionContext) (err error) {
		// delete spiders
		if err := service.NewModelServiceV2[models.SpiderV2]().DeleteMany(bson.M{
			"_id": bson.M{
				"$in": payload.Ids,
			},
		}); err != nil {
			return err
		}

		// delete spider stats
		if err := service.NewModelServiceV2[models.SpiderStatV2]().DeleteMany(bson.M{
			"_id": bson.M{
				"$in": payload.Ids,
			},
		}); err != nil {
			return err
		}

		// related tasks
		tasks, err := service.NewModelServiceV2[models.TaskV2]().GetMany(bson.M{"spider_id": bson.M{"$in": payload.Ids}}, nil)
		if err != nil {
			return err
		}

		if len(tasks) == 0 {
			return nil
		}

		// task ids
		var taskIds []primitive.ObjectID
		for _, t := range tasks {
			taskIds = append(taskIds, t.Id)
		}

		// delete related tasks
		if err := service.NewModelServiceV2[models.TaskV2]().DeleteMany(bson.M{"_id": bson.M{"$in": taskIds}}); err != nil {
			return err
		}

		// delete related task stats
		if err := service.NewModelServiceV2[models.TaskStatV2]().DeleteMany(bson.M{"_id": bson.M{"$in": taskIds}}); err != nil {
			return err
		}

		// delete tasks logs
		wg := sync.WaitGroup{}
		wg.Add(len(taskIds))
		for _, id := range taskIds {
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

		return nil
	}); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccess(c)
}

func GetSpiderListDir(c *gin.Context) {
	path := c.Query("path")

	fsSvc, err := getSpiderFsSvc(c)
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	files, err := fsSvc.List(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			HandleErrorInternalServerError(c, err)
			return
		}
	}

	HandleSuccessWithData(c, files)
}

func GetSpiderFile(c *gin.Context) {
	path := c.Query("path")

	fsSvc, err := getSpiderFsSvc(c)
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	data, err := fsSvc.GetFile(path)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccessWithData(c, string(data))
}

func GetSpiderFileInfo(c *gin.Context) {
	path := c.Query("path")

	fsSvc, err := getSpiderFsSvc(c)
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	info, err := fsSvc.GetFileInfo(path)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccessWithData(c, info)
}

func PostSpiderSaveFile(c *gin.Context) {
	fsSvc, err := getSpiderFsSvc(c)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	if c.GetHeader("Content-Type") == "application/json" {
		var payload struct {
			Path string `json:"path"`
			Data string `json:"data"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			HandleErrorBadRequest(c, err)
			return
		}
		if err := fsSvc.Save(payload.Path, []byte(payload.Data)); err != nil {
			HandleErrorInternalServerError(c, err)
			return
		}
	} else {
		path, ok := c.GetPostForm("path")
		if !ok {
			HandleErrorBadRequest(c, errors.New("missing required field 'path'"))
			return
		}
		file, err := c.FormFile("file")
		if err != nil {
			HandleErrorBadRequest(c, err)
			return
		}
		f, err := file.Open()
		if err != nil {
			HandleErrorBadRequest(c, err)
			return
		}
		fileData, err := io.ReadAll(f)
		if err != nil {
			HandleErrorBadRequest(c, err)
			return
		}
		if err := fsSvc.Save(path, fileData); err != nil {
			HandleErrorInternalServerError(c, err)
			return
		}
	}

	HandleSuccess(c)
}

func PostSpiderSaveFiles(c *gin.Context) {
	fsSvc, err := getSpiderFsSvc(c)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(len(form.File))
	for path := range form.File {
		go func(path string) {
			file, err := c.FormFile(path)
			if err != nil {
				log2.Warnf("invalid file header: %s", path)
				log2.Error(err.Error())
				wg.Done()
				return
			}
			f, err := file.Open()
			if err != nil {
				log2.Warnf("unable to open file: %s", path)
				log2.Error(err.Error())
				wg.Done()
				return
			}
			fileData, err := io.ReadAll(f)
			if err != nil {
				log2.Warnf("unable to read file: %s", path)
				log2.Error(err.Error())
				wg.Done()
				return
			}
			if err := fsSvc.Save(path, fileData); err != nil {
				log2.Warnf("unable to save file: %s", path)
				log2.Error(err.Error())
				wg.Done()
				return
			}
			wg.Done()
		}(path)
	}
	wg.Wait()

	HandleSuccess(c)
}

func PostSpiderSaveDir(c *gin.Context) {
	var payload struct {
		Path    string `json:"path"`
		NewPath string `json:"new_path"`
		Data    string `json:"data"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	fsSvc, err := getSpiderFsSvc(c)
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	if err := fsSvc.CreateDir(payload.Path); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccess(c)
}

func PostSpiderRenameFile(c *gin.Context) {
	var payload struct {
		Path    string `json:"path"`
		NewPath string `json:"new_path"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	fsSvc, err := getSpiderFsSvc(c)
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	if err := fsSvc.Rename(payload.Path, payload.NewPath); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
}

func DeleteSpiderFile(c *gin.Context) {
	var payload struct {
		Path string `json:"path"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	if payload.Path == "~" {
		payload.Path = "."
	}

	fsSvc, err := getSpiderFsSvc(c)
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	if err := fsSvc.Delete(payload.Path); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	_, err = fsSvc.GetFileInfo(".")
	if err != nil {
		_ = fsSvc.CreateDir("/")
	}

	HandleSuccess(c)
}

func PostSpiderCopyFile(c *gin.Context) {
	var payload struct {
		Path    string `json:"path"`
		NewPath string `json:"new_path"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	fsSvc, err := getSpiderFsSvc(c)
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	if err := fsSvc.Copy(payload.Path, payload.NewPath); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccess(c)
}

func PostSpiderExport(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	adminSvc, err := admin.GetSpiderAdminServiceV2()
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// zip file path
	zipFilePath, err := adminSvc.Export(id)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// download
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", zipFilePath))
	c.File(zipFilePath)
}

func PostSpiderRun(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
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

	adminSvc, err := admin.GetSpiderAdminServiceV2()
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// schedule
	taskIds, err := adminSvc.Schedule(id, &opts)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccessWithData(c, taskIds)
}

func GetSpiderGit(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// git client
	gitClient, err := getSpiderGitClient(id)
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
	ignore, err := getSpiderGitIgnore(id)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// git
	_git, err := service.NewModelServiceV2[models.GitV2]().GetById(id)
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

func GetSpiderGitRemoteRefs(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// remote name
	remoteName := c.Query("remote")
	if remoteName == "" {
		remoteName = vcs.GitRemoteNameOrigin
	}

	// git client
	gitClient, err := getSpiderGitClient(id)
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

func PostSpiderGitCheckout(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// payload
	var payload struct {
		Paths         []string `json:"paths"`
		CommitMessage string   `json:"commit_message"`
		Branch        string   `json:"branch"`
		Tag           string   `json:"tag"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// git client
	gitClient, err := getSpiderGitClient(id)
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
	if err := gitSpiderCheckout(gitClient, constants.GitRemoteNameOrigin, branch); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccess(c)
}

func PostSpiderGitPull(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// payload
	var payload struct {
		Paths         []string `json:"paths"`
		CommitMessage string   `json:"commit_message"`
		Branch        string   `json:"branch"`
		Tag           string   `json:"tag"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// git
	g, err := service.NewModelServiceV2[models.GitV2]().GetById(id)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// attempt to sync git
	adminSvc, err := admin.GetSpiderAdminServiceV2()
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	_ = adminSvc.SyncGitOne(g)

	HandleSuccess(c)
}

func PostSpiderGitCommit(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// payload
	var payload entity.GitPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// git client
	gitClient, err := getSpiderGitClient(id)
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

func GetSpiderDataSource(c *gin.Context) {
	// spider id
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// spider
	s, err := service.NewModelServiceV2[models.SpiderV2]().GetById(id)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// data source
	ds, err := service.NewModelServiceV2[models.DataSourceV2]().GetById(s.DataSourceId)
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

func PostSpiderDataSource(c *gin.Context) {
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
	s, err := service.NewModelServiceV2[models.SpiderV2]().GetById(id)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// data source
	if !dsId.IsZero() {
		_, err = service.NewModelServiceV2[models.DataSourceV2]().GetById(dsId)
		if err != nil {
			HandleErrorInternalServerError(c, err)
			return
		}
	}

	// save data source id
	u := GetUserFromContextV2(c)
	s.DataSourceId = dsId
	s.SetUpdatedBy(u.Id)
	_, err = service.NewModelServiceV2[models.SpiderV2]().InsertOne(*s)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccess(c)
}

func getSpiderFsSvc(c *gin.Context) (svc interfaces.FsServiceV2, err error) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return nil, err
	}

	workspacePath := viper.GetString("workspace")
	fsSvc := fs.NewFsServiceV2(filepath.Join(workspacePath, id.Hex()))

	return fsSvc, nil
}

func getSpiderFsSvcById(id primitive.ObjectID) interfaces.FsServiceV2 {
	workspacePath := viper.GetString("workspace")
	fsSvc := fs.NewFsServiceV2(filepath.Join(workspacePath, id.Hex()))
	return fsSvc
}

func getSpiderGitClient(id primitive.ObjectID) (client *vcs.GitClient, err error) {
	// git
	g, err := service.NewModelServiceV2[models.GitV2]().GetById(id)
	if err != nil {
		if !errors.Is(err, mongo2.ErrNoDocuments) {
			return nil, trace.TraceError(err)
		}
		return nil, nil
	}

	// git client
	workspacePath := viper.GetString("workspace")
	client, err = vcs.NewGitClient(vcs.WithPath(filepath.Join(workspacePath, id.Hex())))
	if err != nil {
		return nil, err
	}

	// set auth
	utils.InitGitClientAuthV2(g, client)

	// remote name
	remoteName := vcs.GitRemoteNameOrigin

	// update remote
	r, err := client.GetRemote(remoteName)
	if errors.Is(err, git.ErrRemoteNotFound) {
		// remote not exists, create
		if _, err := client.CreateRemote(&config.RemoteConfig{
			Name: remoteName,
			URLs: []string{g.Url},
		}); err != nil {
			return nil, trace.TraceError(err)
		}
	} else if err == nil {
		// remote exists, update if different
		if g.Url != r.Config().URLs[0] {
			if err := client.DeleteRemote(remoteName); err != nil {
				return nil, trace.TraceError(err)
			}
			if _, err := client.CreateRemote(&config.RemoteConfig{
				Name: remoteName,
				URLs: []string{g.Url},
			}); err != nil {
				return nil, trace.TraceError(err)
			}
		}
		client.SetRemoteUrl(g.Url)
	} else {
		// error
		return nil, trace.TraceError(err)
	}

	// check if head reference exists
	_, err = client.GetRepository().Head()
	if err == nil {
		return client, nil
	}

	// align master/main branch
	alignSpiderGitBranch(client)

	return client, nil
}

func alignSpiderGitBranch(gitClient *vcs.GitClient) {
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
	defaultRemoteBranch, err := getSpiderDefaultRemoteBranch(refs)
	if err != nil || defaultRemoteBranch == "" {
		return
	}

	// move branch
	if err := gitClient.MoveBranch(vcs.GitBranchNameMaster, defaultRemoteBranch); err != nil {
		trace.PrintError(err)
	}
}

func getSpiderDefaultRemoteBranch(refs []vcs.GitRef) (defaultRemoteBranchName string, err error) {
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

func getSpiderGitIgnore(id primitive.ObjectID) (ignore []string, err error) {
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

func gitSpiderCheckout(gitClient *vcs.GitClient, remote, branch string) (err error) {
	if err := gitClient.CheckoutBranch(branch, vcs.WithBranch(branch)); err != nil {
		return trace.TraceError(err)
	}

	// pull
	return spiderGitPull(gitClient, remote, branch)
}

func spiderGitPull(gitClient *vcs.GitClient, remote, branch string) (err error) {
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

func upsertSpiderDataCollection(s *models.SpiderV2) (err error) {
	modelSvc := service.NewModelServiceV2[models.DataCollectionV2]()
	if s.ColId.IsZero() {
		// validate
		if s.ColName == "" {
			return errors.New("data collection name is required")
		}
		// no id
		dc, err := modelSvc.GetOne(bson.M{"name": s.ColName}, nil)
		if err != nil {
			if errors.Is(err, mongo2.ErrNoDocuments) {
				// not exists, add new
				dc = &models.DataCollectionV2{Name: s.ColName}
				dcId, err := modelSvc.InsertOne(*dc)
				if err != nil {
					return err
				}
				dc.SetId(dcId)
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
		dc, err := modelSvc.GetById(s.ColId)
		if err != nil {
			return err
		}
		s.ColId = dc.Id
	}
	return nil
}
