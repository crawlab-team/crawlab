package controllers

import (
	"errors"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/fs"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/core/spider/admin"
	"github.com/crawlab-team/crawlab/core/utils"
	"github.com/crawlab-team/crawlab/db/mongo"
	"github.com/crawlab-team/crawlab/trace"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"math"
	"os"
	"path/filepath"
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
	// get all list
	all := MustGetFilterAll(c)
	if all {
		NewControllerV2[models.SpiderV2]().getAll(c)
		return
	}

	// get list
	withStats := c.Query("stats")
	if withStats == "" {
		NewControllerV2[models.SpiderV2]().GetList(c)
		return
	}

	// get list with stats
	getSpiderListWithStats(c)
}

func getSpiderListWithStats(c *gin.Context) {
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

	// user
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
	fsSvc, err := getSpiderFsSvcById(id)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	err = fsSvc.CreateDir(".")
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
					log.Warnf("failed to remove task log directory: %s", logPath)
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

	go func() {
		// spider
		s, err := service.NewModelServiceV2[models.SpiderV2]().GetById(id)
		if err != nil {
			log.Errorf("failed to get spider: %s", err.Error())
			trace.PrintError(err)
			return
		}

		// skip spider with git
		if !s.GitId.IsZero() {
			return
		}

		// delete spider directory
		fsSvc, err := getSpiderFsSvcById(id)
		if err != nil {
			log.Errorf("failed to get spider fs service: %s", err.Error())
			trace.PrintError(err)
			return
		}
		err = fsSvc.Delete(".")
		if err != nil {
			log.Errorf("failed to delete spider directory: %s", err.Error())
			trace.PrintError(err)
			return
		}
	}()

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
					log.Warnf("failed to remove task log directory: %s", logPath)
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

	// delete spider directories
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(len(payload.Ids))
		for _, id := range payload.Ids {
			go func(id primitive.ObjectID) {
				defer wg.Done()

				// spider
				s, err := service.NewModelServiceV2[models.SpiderV2]().GetById(id)
				if err != nil {
					log.Errorf("failed to get spider: %s", err.Error())
					trace.PrintError(err)
					return
				}

				// skip spider with git
				if !s.GitId.IsZero() {
					return
				}

				// delete spider directory
				fsSvc, err := getSpiderFsSvcById(id)
				if err != nil {
					log.Errorf("failed to get spider fs service: %s", err.Error())
					trace.PrintError(err)
					return
				}
				err = fsSvc.Delete(".")
				if err != nil {
					log.Errorf("failed to delete spider directory: %s", err.Error())
					trace.PrintError(err)
					return
				}
			}(id)
		}
		wg.Wait()
	}()

	HandleSuccess(c)
}

func GetSpiderListDir(c *gin.Context) {
	rootPath, err := getSpiderRootPath(c)
	if err != nil {
		HandleErrorForbidden(c, err)
		return
	}
	GetBaseFileListDir(rootPath, c)
}

func GetSpiderFile(c *gin.Context) {
	rootPath, err := getSpiderRootPath(c)
	if err != nil {
		HandleErrorForbidden(c, err)
		return
	}
	GetBaseFileFile(rootPath, c)
}

func GetSpiderFileInfo(c *gin.Context) {
	rootPath, err := getSpiderRootPath(c)
	if err != nil {
		HandleErrorForbidden(c, err)
		return
	}
	GetBaseFileFileInfo(rootPath, c)
}

func PostSpiderSaveFile(c *gin.Context) {
	rootPath, err := getSpiderRootPath(c)
	if err != nil {
		HandleErrorForbidden(c, err)
		return
	}
	PostBaseFileSaveFile(rootPath, c)
}

func PostSpiderSaveFiles(c *gin.Context) {
	rootPath, err := getSpiderRootPath(c)
	if err != nil {
		HandleErrorForbidden(c, err)
		return
	}
	PostBaseFileSaveFiles(rootPath, c)
}

func PostSpiderSaveDir(c *gin.Context) {
	rootPath, err := getSpiderRootPath(c)
	if err != nil {
		HandleErrorForbidden(c, err)
		return
	}
	PostBaseFileSaveDir(rootPath, c)
}

func PostSpiderRenameFile(c *gin.Context) {
	rootPath, err := getSpiderRootPath(c)
	if err != nil {
		HandleErrorForbidden(c, err)
		return
	}
	PostBaseFileRenameFile(rootPath, c)
}

func DeleteSpiderFile(c *gin.Context) {
	rootPath, err := getSpiderRootPath(c)
	if err != nil {
		HandleErrorForbidden(c, err)
		return
	}
	DeleteBaseFileFile(rootPath, c)
}

func PostSpiderCopyFile(c *gin.Context) {
	rootPath, err := getSpiderRootPath(c)
	if err != nil {
		HandleErrorForbidden(c, err)
		return
	}
	PostBaseFileCopyFile(rootPath, c)
}

func PostSpiderExport(c *gin.Context) {
	rootPath, err := getSpiderRootPath(c)
	if err != nil {
		HandleErrorForbidden(c, err)
		return
	}
	PostBaseFileExport(rootPath, c)
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

func getSpiderFsSvc(s *models.SpiderV2) (svc interfaces.FsServiceV2, err error) {
	workspacePath := viper.GetString("workspace")
	fsSvc := fs.NewFsServiceV2(filepath.Join(workspacePath, s.Id.Hex()))

	return fsSvc, nil
}

func GetSpiderFsSvcById(id primitive.ObjectID) (svc interfaces.FsServiceV2, err error) {
	return getSpiderFsSvcById(id)
}

func getSpiderFsSvcById(id primitive.ObjectID) (svc interfaces.FsServiceV2, err error) {
	s, err := service.NewModelServiceV2[models.SpiderV2]().GetById(id)
	if err != nil {
		log.Errorf("failed to get spider: %s", err.Error())
		trace.PrintError(err)
		return nil, err
	}
	return getSpiderFsSvc(s)
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

func UpsertSpiderDataCollection(s *models.SpiderV2) (err error) {
	return upsertSpiderDataCollection(s)
}

func getSpiderRootPath(c *gin.Context) (rootPath string, err error) {
	// spider id
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return "", err
	}

	// spider
	s, err := service.NewModelServiceV2[models.SpiderV2]().GetById(id)
	if err != nil {
		return "", err
	}

	// check git permission
	if !utils.IsPro() && !s.GitId.IsZero() {
		return "", errors.New("git is not allowed in the community version")
	}

	// if git id is zero, return spider id as root path
	if s.GitId.IsZero() {
		return id.Hex(), nil
	}

	return filepath.Join(s.GitId.Hex(), rootPath), nil
}
