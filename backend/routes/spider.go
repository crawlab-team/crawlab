package routes

import (
	"crawlab/constants"
	"crawlab/database"
	"crawlab/entity"
	"crawlab/model"
	"crawlab/services"
	"crawlab/utils"
	"fmt"
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

// ======== 爬虫管理 ========

// @Summary Get spider list
// @Description Get spider list
// @Tags spider
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param page_num query string false "page num"
// @Param page_size query string false "page size"
// @Param keyword query string false "keyword"
// @Param project_id query string false "project_id"
// @Param type query string false "type"
// @Param sort_key query  string  false "sort_key"
// @Param sort_direction query string false "sort_direction"
// @Param owner_type query string false "owner_type"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /schedules [get]
func GetSpiderList(c *gin.Context) {
	pageNum := c.Query("page_num")
	pageSize := c.Query("page_size")
	keyword := c.Query("keyword")
	pid := c.Query("project_id")
	t := c.Query("type")
	sortKey := c.Query("sort_key")
	sortDirection := c.Query("sort_direction")
	ownerType := c.Query("owner_type")

	// 筛选-名称
	filter := bson.M{
		"name": bson.M{"$regex": bson.RegEx{Pattern: keyword, Options: "im"}},
	}

	// 筛选-类型
	if t != "" && t != "all" {
		filter["type"] = t
	}

	// 筛选-是否为长任务
	if t == "long-task" {
		delete(filter, "type")
		filter["is_long_task"] = true
	}

	// 筛选-项目
	if pid == "" {
		// do nothing
	} else if pid == constants.ObjectIdNull {
		filter["$or"] = []bson.M{
			{"project_id": bson.ObjectIdHex(pid)},
			{"project_id": bson.M{"$exists": false}},
		}
	} else {
		filter["project_id"] = bson.ObjectIdHex(pid)
	}

	// 筛选-用户
	if ownerType == constants.OwnerTypeAll {
		user := services.GetCurrentUser(c)
		if user.Role == constants.RoleNormal {
			filter["$or"] = []bson.M{
				{"user_id": services.GetCurrentUserId(c)},
				{"is_public": true},
			}
		}
	} else if ownerType == constants.OwnerTypeMe {
		filter["user_id"] = services.GetCurrentUserId(c)
	} else if ownerType == constants.OwnerTypePublic {
		filter["is_public"] = true
	}

	// 排序
	sortStr := "-_id"
	if sortKey != "" && sortDirection != "" {
		if sortDirection == constants.DESCENDING {
			sortStr = "-" + sortKey
		} else if sortDirection == constants.ASCENDING {
			sortStr = "+" + sortKey
		} else {
			HandleErrorF(http.StatusBadRequest, c, "invalid sort_direction")
			return
		}
	}

	// 分页
	page := &entity.Page{}
	page.GetPage(pageNum, pageSize)

	results, count, err := model.GetSpiderList(filter, page.Skip, page.Limit, sortStr)
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    bson.M{"list": results, "total": count},
	})
}

// @Summary Get spider by id
// @Description Get spider  by id
// @Tags spider
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id path string true "schedule id"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /spiders/{id} [get]
func GetSpider(c *gin.Context) {
	id := c.Param("id")

	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
		return
	}

	spider, err := model.GetSpider(bson.ObjectIdHex(id))
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    spider,
	})
}

// @Summary Post spider
// @Description Post spider
// @Tags spider
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id path string true "schedule id"
// @Param item body model.Spider true "spider item"
// @Success 200 json string Response
// @Failure 500 json string Response
// @Router /spiders/{id} [post]
func PostSpider(c *gin.Context) {
	id := c.Param("id")

	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
		return
	}

	var item model.Spider
	if err := c.ShouldBindJSON(&item); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	// UserId
	if !item.UserId.Valid() {
		item.UserId = bson.ObjectIdHex(constants.ObjectIdNull)
	}

	if err := model.UpdateSpider(bson.ObjectIdHex(id), item); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 更新 GitCron
	if err := services.GitCron.Update(); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 获取爬虫
	spider, err := model.GetSpider(bson.ObjectIdHex(id))
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 去重处理
	if err := services.UpdateSpiderDedup(spider); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

// @Summary Publish spider
// @Description Publish spider
// @Tags spider
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id path string true "schedule id"
// @Success 200 json string Response
// @Failure 500 json string Response
// @Router /spiders/{id}/publish [post]
func PublishSpider(c *gin.Context) {
	id := c.Param("id")

	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
		return
	}

	spider, err := model.GetSpider(bson.ObjectIdHex(id))
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	services.PublishSpider(spider)

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

// @Summary Put spider
// @Description Put spider
// @Tags spider
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param spider body model.Spider true "spider item"
// @Success 200 json string Response
// @Failure 500 json string Response
// @Router /spiders [put]
func PutSpider(c *gin.Context) {
	var spider model.Spider
	if err := c.ShouldBindJSON(&spider); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	// 爬虫名称不能为空
	if spider.Name == "" {
		HandleErrorF(http.StatusBadRequest, c, "spider name should not be empty")
		return
	}

	// 判断爬虫是否存在
	if spider := model.GetSpiderByName(spider.Name); spider.Name != "" {
		HandleErrorF(http.StatusBadRequest, c, fmt.Sprintf("spider for '%s' already exists", spider.Name))
		return
	}

	// 设置爬虫类别
	spider.Type = constants.Customized

	// 将FileId置空
	spider.FileId = bson.ObjectIdHex(constants.ObjectIdNull)

	// UserId
	spider.UserId = services.GetCurrentUserId(c)

	// 爬虫目录
	spiderDir := filepath.Join(viper.GetString("spider.path"), spider.Name)

	// 赋值到爬虫实例
	spider.Src = spiderDir

	// 移除已有爬虫目录
	if utils.Exists(spiderDir) {
		if err := os.RemoveAll(spiderDir); err != nil {
			HandleError(http.StatusInternalServerError, c, err)
			return
		}
	}

	// 生成爬虫目录
	if err := os.MkdirAll(spiderDir, 0777); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 如果为 Scrapy 项目，生成 Scrapy 项目
	if spider.IsScrapy {
		if err := services.CreateScrapyProject(spider); err != nil {
			HandleError(http.StatusInternalServerError, c, err)
			return
		}
	}

	// 添加爬虫到数据库
	if err := spider.Add(); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 同步到GridFS
	if err := services.UploadSpiderToGridFsFromMaster(spider); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 更新 GitCron
	if err := services.GitCron.Update(); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    spider,
	})
}

// @Summary Copy spider
// @Description Copy spider
// @Tags spider
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id path string true "schedule id"
// @Success 200 json string Response
// @Failure 500 json string Response
// @Router /spiders/{id}/copy [post]
func CopySpider(c *gin.Context) {
	type ReqBody struct {
		Name string `json:"name"`
	}

	id := c.Param("id")

	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
		return
	}

	var reqBody ReqBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	// 检查新爬虫名称是否存在
	// 如果存在，则返回错误
	s := model.GetSpiderByName(reqBody.Name)
	if s.Name != "" {
		HandleErrorF(http.StatusBadRequest, c, fmt.Sprintf("spider name '%s' already exists", reqBody.Name))
		return
	}

	// 被复制爬虫
	spider, err := model.GetSpider(bson.ObjectIdHex(id))
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// UserId
	spider.UserId = services.GetCurrentUserId(c)

	// 复制爬虫
	if err := services.CopySpider(spider, reqBody.Name); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

// @Summary Upload spider
// @Description Upload spider
// @Tags spider
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param file formData file true "spider file to upload"
// @Param name formData string true "spider name"
// @Param display_name formData string true "display name"
// @Param col formData string true "col"
// @Param cmd formData string true "cmd"
// @Success 200 json string Response
// @Failure 500 json string Response
// @Router /spiders [post]
func UploadSpider(c *gin.Context) {
	// 从body中获取文件
	uploadFile, err := c.FormFile("file")
	if err != nil {
		debug.PrintStack()
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 获取参数
	name := c.PostForm("name")
	displayName := c.PostForm("display_name")
	col := c.PostForm("col")
	cmd := c.PostForm("cmd")

	// 如果不为zip文件，返回错误
	if !strings.HasSuffix(uploadFile.Filename, ".zip") {
		HandleError(http.StatusBadRequest, c, errors.New("not a valid zip file"))
		return
	}

	// 以防tmp目录不存在
	tmpPath := viper.GetString("other.tmppath")
	if !utils.Exists(tmpPath) {
		if err := os.MkdirAll(tmpPath, os.ModePerm); err != nil {
			log.Error("mkdir other.tmppath dir error:" + err.Error())
			debug.PrintStack()
			HandleError(http.StatusBadRequest, c, errors.New("mkdir other.tmppath dir error"))
			return
		}
	}

	// 保存到本地临时文件
	randomId := uuid.NewV4()
	tmpFilePath := filepath.Join(tmpPath, randomId.String()+".zip")
	if err := c.SaveUploadedFile(uploadFile, tmpFilePath); err != nil {
		log.Error("save upload file error: " + err.Error())
		debug.PrintStack()
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 获取 GridFS 实例
	s, gf := database.GetGridFs("files")
	defer s.Close()

	// 判断文件是否已经存在
	var gfFile model.GridFs
	if err := gf.Find(bson.M{"filename": uploadFile.Filename}).One(&gfFile); err == nil {
		// 已经存在文件，则删除
		if err := gf.RemoveId(gfFile.Id); err != nil {
			log.Errorf("remove grid fs error: %s", err.Error())
			debug.PrintStack()
			HandleError(http.StatusInternalServerError, c, err)
			return
		}
	}

	// 上传到GridFs
	fid, err := services.RetryUploadToGridFs(uploadFile.Filename, tmpFilePath)
	if err != nil {
		log.Errorf("upload to grid fs error: %s", err.Error())
		debug.PrintStack()
		return
	}

	idx := strings.LastIndex(uploadFile.Filename, "/")
	targetFilename := uploadFile.Filename[idx+1:]

	// 判断爬虫是否存在
	spiderName := strings.Replace(targetFilename, ".zip", "", 1)
	if name != "" {
		spiderName = name
	}
	spider := model.GetSpiderByName(spiderName)
	if spider.Name == "" {
		// 保存爬虫信息
		srcPath := viper.GetString("spider.path")
		spider := model.Spider{
			Name:        spiderName,
			DisplayName: spiderName,
			Type:        constants.Customized,
			Src:         filepath.Join(srcPath, spiderName),
			FileId:      fid,
			ProjectId:   bson.ObjectIdHex(constants.ObjectIdNull),
			UserId:      services.GetCurrentUserId(c),
		}
		if name != "" {
			spider.Name = name
		}
		if displayName != "" {
			spider.DisplayName = displayName
		}
		if col != "" {
			spider.Col = col
		}
		if cmd != "" {
			spider.Cmd = cmd
		}
		if err := spider.Add(); err != nil {
			log.Error("add spider error: " + err.Error())
			debug.PrintStack()
			HandleError(http.StatusInternalServerError, c, err)
			return
		}
	} else {
		if name != "" {
			spider.Name = name
		}
		if displayName != "" {
			spider.DisplayName = displayName
		}
		if col != "" {
			spider.Col = col
		}
		if cmd != "" {
			spider.Cmd = cmd
		}
		// 更新file_id
		spider.FileId = fid
		if err := spider.Save(); err != nil {
			log.Error("add spider error: " + err.Error())
			debug.PrintStack()
			HandleError(http.StatusInternalServerError, c, err)
			return
		}
	}

	// 获取爬虫
	spider = model.GetSpiderByName(spiderName)

	// 发起同步
	services.PublishSpider(spider)

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    spider,
	})
}

// @Summary Upload spider by id
// @Description Upload spider by id
// @Tags spider
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param file formData file true "spider file to upload"
// @Param id path string true "spider id"
// @Success 200 json string Response
// @Failure 500 json string Response
// @Router /spiders/{id}/upload [post]
func UploadSpiderFromId(c *gin.Context) {
	// TODO: 与 UploadSpider 部分逻辑重复，需要优化代码
	// 爬虫ID
	spiderId := c.Param("id")
	if !bson.IsObjectIdHex(spiderId) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
		return
	}
	// 获取爬虫
	spider, err := model.GetSpider(bson.ObjectIdHex(spiderId))
	if err != nil {
		if err == mgo.ErrNotFound {
			HandleErrorF(http.StatusNotFound, c, "cannot find spider")
		} else {
			HandleError(http.StatusInternalServerError, c, err)
		}
		return
	}

	// 从body中获取文件
	uploadFile, err := c.FormFile("file")
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 如果不为zip文件，返回错误
	if !strings.HasSuffix(uploadFile.Filename, ".zip") {
		debug.PrintStack()
		HandleError(http.StatusBadRequest, c, errors.New("Not a valid zip file"))
		return
	}

	// 以防tmp目录不存在
	tmpPath := viper.GetString("other.tmppath")
	if !utils.Exists(tmpPath) {
		if err := os.MkdirAll(tmpPath, os.ModePerm); err != nil {
			log.Error("mkdir other.tmppath dir error:" + err.Error())
			debug.PrintStack()
			HandleError(http.StatusBadRequest, c, errors.New("Mkdir other.tmppath dir error"))
			return
		}
	}

	// 保存到本地临时文件
	randomId := uuid.NewV4()
	tmpFilePath := filepath.Join(tmpPath, randomId.String()+".zip")
	if err := c.SaveUploadedFile(uploadFile, tmpFilePath); err != nil {
		log.Error("save upload file error: " + err.Error())
		debug.PrintStack()
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 获取 GridFS 实例
	s, gf := database.GetGridFs("files")
	defer s.Close()

	// 判断文件是否已经存在
	var gfFile model.GridFs
	if err := gf.Find(bson.M{"filename": spider.Name}).One(&gfFile); err == nil {
		// 已经存在文件，则删除
		if err := gf.RemoveId(gfFile.Id); err != nil {
			log.Errorf("remove grid fs error: " + err.Error())
			debug.PrintStack()
			HandleError(http.StatusInternalServerError, c, err)
			return
		}
	}

	// 上传到GridFs
	fid, err := services.RetryUploadToGridFs(spider.Name, tmpFilePath)
	if err != nil {
		log.Errorf("upload to grid fs error: %s", err.Error())
		debug.PrintStack()
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 更新file_id
	spider.FileId = fid
	if err := spider.Save(); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return
	}

	// 发起同步
	services.PublishSpider(spider)

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

// @Summary Delete spider by id
// @Description Delete spider by id
// @Tags spider
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id path string true "spider id"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /spiders/{id} [delete]
func DeleteSpider(c *gin.Context) {
	id := c.Param("id")

	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
		return
	}

	if err := services.RemoveSpider(id); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 更新 GitCron
	if err := services.GitCron.Update(); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

// @Summary delete spider
// @Description delete spider
// @Tags spider
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Success 200 json string Response
// @Failure 500 json string Response
// @Router /spiders [post]
func DeleteSelectedSpider(c *gin.Context) {
	type ReqBody struct {
		SpiderIds []string `json:"spider_ids"`
	}

	var reqBody ReqBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		HandleErrorF(http.StatusBadRequest, c, "invalid request")
		return
	}

	for _, spiderId := range reqBody.SpiderIds {
		if err := services.RemoveSpider(spiderId); err != nil {
			HandleError(http.StatusInternalServerError, c, err)
			return
		}
	}

	// 更新 GitCron
	if err := services.GitCron.Update(); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

// @Summary cancel spider
// @Description cancel spider
// @Tags spider
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Success 200 json string Response
// @Failure 500 json string Response
// @Router /spiders-cancel [post]
func CancelSelectedSpider(c *gin.Context) {
	type ReqBody struct {
		SpiderIds []string `json:"spider_ids"`
	}

	var reqBody ReqBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		HandleErrorF(http.StatusBadRequest, c, "invalid request")
		return
	}

	for _, spiderId := range reqBody.SpiderIds {
		if err := services.CancelSpider(spiderId); err != nil {
			log.Errorf(err.Error())
			debug.PrintStack()
		}
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

// @Summary run spider
// @Description run spider
// @Tags spider
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Success 200 json string Response
// @Failure 500 json string Response
// @Router /spiders-run [post]
func RunSelectedSpider(c *gin.Context) {
	type TaskParam struct {
		SpiderId bson.ObjectId `json:"spider_id"`
		Param    string        `json:"param"`
	}
	type ReqBody struct {
		RunType    string          `json:"run_type"`
		NodeIds    []bson.ObjectId `json:"node_ids"`
		TaskParams []TaskParam     `json:"task_params"`
	}

	var reqBody ReqBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		HandleErrorF(http.StatusBadRequest, c, "invalid request")
		return
	}

	// 任务ID
	var taskIds []string

	// 遍历爬虫
	// TODO: 优化此部分代码，与 routes.PutTask 有重合部分
	for _, taskParam := range reqBody.TaskParams {
		if reqBody.RunType == constants.RunTypeAllNodes {
			// 所有节点
			nodes, err := model.GetNodeList(nil)
			if err != nil {
				HandleError(http.StatusInternalServerError, c, err)
				return
			}
			for _, node := range nodes {
				t := model.Task{
					SpiderId:   taskParam.SpiderId,
					NodeId:     node.Id,
					Param:      taskParam.Param,
					UserId:     services.GetCurrentUserId(c),
					RunType:    constants.RunTypeAllNodes,
					ScheduleId: bson.ObjectIdHex(constants.ObjectIdNull),
				}

				id, err := services.AddTask(t)
				if err != nil {
					HandleError(http.StatusInternalServerError, c, err)
					return
				}

				taskIds = append(taskIds, id)
			}
		} else if reqBody.RunType == constants.RunTypeRandom {
			// 随机
			t := model.Task{
				SpiderId:   taskParam.SpiderId,
				Param:      taskParam.Param,
				UserId:     services.GetCurrentUserId(c),
				RunType:    constants.RunTypeRandom,
				ScheduleId: bson.ObjectIdHex(constants.ObjectIdNull),
			}
			id, err := services.AddTask(t)
			if err != nil {
				HandleError(http.StatusInternalServerError, c, err)
				return
			}
			taskIds = append(taskIds, id)
		} else if reqBody.RunType == constants.RunTypeSelectedNodes {
			// 指定节点
			for _, nodeId := range reqBody.NodeIds {
				t := model.Task{
					SpiderId:   taskParam.SpiderId,
					NodeId:     nodeId,
					Param:      taskParam.Param,
					UserId:     services.GetCurrentUserId(c),
					RunType:    constants.RunTypeSelectedNodes,
					ScheduleId: bson.ObjectIdHex(constants.ObjectIdNull),
				}

				id, err := services.AddTask(t)
				if err != nil {
					HandleError(http.StatusInternalServerError, c, err)
					return
				}
				taskIds = append(taskIds, id)
			}
		} else {
			HandleErrorF(http.StatusInternalServerError, c, "invalid run_type")
			return
		}
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    taskIds,
	})
}

// @Summary Get task list
// @Description Get task list
// @Tags spider
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id path string true "spider id"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /spiders/{id}/tasks [get]
func GetSpiderTasks(c *gin.Context) {
	id := c.Param("id")
	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
		return
	}
	spider, err := model.GetSpider(bson.ObjectIdHex(id))
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	tasks, err := spider.GetTasks()
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    tasks,
	})
}

// @Summary Get spider stats
// @Description Get spider stats
// @Tags spider
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id path string true "spider id"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /spiders/{id}/stats [get]
func GetSpiderStats(c *gin.Context) {
	type Overview struct {
		TaskCount            int     `json:"task_count" bson:"task_count"`
		ResultCount          int     `json:"result_count" bson:"result_count"`
		SuccessCount         int     `json:"success_count" bson:"success_count"`
		SuccessRate          float64 `json:"success_rate"`
		TotalWaitDuration    float64 `json:"wait_duration" bson:"wait_duration"`
		TotalRuntimeDuration float64 `json:"runtime_duration" bson:"runtime_duration"`
		AvgWaitDuration      float64 `json:"avg_wait_duration"`
		AvgRuntimeDuration   float64 `json:"avg_runtime_duration"`
	}

	type Data struct {
		Overview Overview              `json:"overview"`
		Daily    []model.TaskDailyItem `json:"daily"`
	}

	id := c.Param("id")
	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
		return
	}
	spider, err := model.GetSpider(bson.ObjectIdHex(id))
	if err != nil {
		log.Errorf(err.Error())
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	s, col := database.GetCol("tasks")
	defer s.Close()

	// 起始日期
	startDate := time.Now().Add(-time.Hour * 24 * 30)
	endDate := time.Now()

	// match
	op1 := bson.M{
		"$match": bson.M{
			"spider_id": spider.Id,
			"create_ts": bson.M{
				"$gte": startDate,
				"$lt":  endDate,
			},
		},
	}

	// project
	op2 := bson.M{
		"$project": bson.M{
			"success_count": bson.M{
				"$cond": []interface{}{
					bson.M{
						"$eq": []string{
							"$status",
							constants.StatusFinished,
						},
					},
					1,
					0,
				},
			},
			"result_count":     "$result_count",
			"wait_duration":    "$wait_duration",
			"runtime_duration": "$runtime_duration",
		},
	}

	// group
	op3 := bson.M{
		"$group": bson.M{
			"_id":              nil,
			"task_count":       bson.M{"$sum": 1},
			"success_count":    bson.M{"$sum": "$success_count"},
			"result_count":     bson.M{"$sum": "$result_count"},
			"wait_duration":    bson.M{"$sum": "$wait_duration"},
			"runtime_duration": bson.M{"$sum": "$runtime_duration"},
		},
	}

	// run aggregation pipeline
	var overview Overview
	if err := col.Pipe([]bson.M{op1, op2, op3}).One(&overview); err != nil {
		if err == mgo.ErrNotFound {
			c.JSON(http.StatusOK, Response{
				Status:  "ok",
				Message: "success",
				Data: Data{
					Overview: overview,
					Daily:    []model.TaskDailyItem{},
				},
			})
			return
		}
		log.Errorf(err.Error())
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 后续处理
	successCount, _ := strconv.ParseFloat(strconv.Itoa(overview.SuccessCount), 64)
	taskCount, _ := strconv.ParseFloat(strconv.Itoa(overview.TaskCount), 64)
	overview.SuccessRate = successCount / taskCount
	overview.AvgWaitDuration = overview.TotalWaitDuration / taskCount
	overview.AvgRuntimeDuration = overview.TotalRuntimeDuration / taskCount

	items, err := model.GetDailyTaskStats(bson.M{"spider_id": spider.Id, "user_id": services.GetCurrentUserId(c)})
	if err != nil {
		log.Errorf(err.Error())
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data: Data{
			Overview: overview,
			Daily:    items,
		},
	})
}

// @Summary Get schedules
// @Description Get schedules
// @Tags spider
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id path string true "spider id"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /spiders/{id}/schedules [get]
func GetSpiderSchedules(c *gin.Context) {
	id := c.Param("id")

	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "spider_id is invalid")
		return
	}

	// 获取定时任务
	list, err := model.GetScheduleList(bson.M{"spider_id": bson.ObjectIdHex(id)})
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    list,
	})
}

// ======== ./爬虫管理 ========

// ======== 爬虫文件管理 ========

// @Summary Get spider dir
// @Description Get spider dir
// @Tags spider
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id path string true "spider id"
// @Param path query string true "path"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /spiders/{id}/dir [get]
func GetSpiderDir(c *gin.Context) {
	// 爬虫ID
	id := c.Param("id")
	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
		return
	}
	// 目录相对路径
	path := c.Query("path")

	// 获取爬虫
	spider, err := model.GetSpider(bson.ObjectIdHex(id))
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 获取目录下文件列表
	spiderPath := viper.GetString("spider.path")
	f, err := ioutil.ReadDir(filepath.Join(spiderPath, spider.Name, path))
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 遍历文件列表
	var fileList []model.File
	for _, file := range f {
		fileList = append(fileList, model.File{
			Name:  file.Name(),
			IsDir: file.IsDir(),
			Size:  file.Size(),
			Path:  filepath.Join(path, file.Name()),
		})
	}

	// 返回结果
	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    fileList,
	})
}

type SpiderFileReqBody struct {
	Path    string `json:"path"`
	Content string `json:"content"`
	NewPath string `json:"new_path"`
}

// @Summary Get spider file
// @Description Get spider file
// @Tags spider
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id path string true "spider id"
// @Param path query string true "path"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /spiders/{id}/file [get]
func GetSpiderFile(c *gin.Context) {
	// 爬虫ID
	id := c.Param("id")
	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
		return
	}
	// 文件相对路径
	path := c.Query("path")

	// 获取爬虫
	spider, err := model.GetSpider(bson.ObjectIdHex(id))
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 读取文件
	fileBytes, err := ioutil.ReadFile(filepath.Join(spider.Src, path))
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
	}

	// 返回结果
	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    utils.BytesToString(fileBytes),
	})
}

// @Summary Get spider dir
// @Description Get spider dir
// @Tags spider
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id path string true "spider id"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /spiders/{id}/file/tree [get]
func GetSpiderFileTree(c *gin.Context) {
	// 爬虫ID
	id := c.Param("id")
	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
		return
	}
	// 获取爬虫
	spider, err := model.GetSpider(bson.ObjectIdHex(id))
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 获取目录下文件列表
	spiderPath := viper.GetString("spider.path")
	spiderFilePath := filepath.Join(spiderPath, spider.Name)

	// 获取文件目录树
	fileNodeTree, err := services.GetFileNodeTree(spiderFilePath, 0)
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    fileNodeTree,
	})
}

// @Summary Post spider file
// @Description Post spider file
// @Tags spider
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id path string true "spider id"
// @Param reqBody body routes.SpiderFileReqBody true "path"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /spiders/{id}/file [post]
func PostSpiderFile(c *gin.Context) {
	// 爬虫ID
	id := c.Param("id")
	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
		return
	}
	// 文件相对路径
	var reqBody SpiderFileReqBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	// 获取爬虫
	spider, err := model.GetSpider(bson.ObjectIdHex(id))
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 写文件
	if err := ioutil.WriteFile(filepath.Join(spider.Src, reqBody.Path), []byte(reqBody.Content), os.ModePerm); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 同步到GridFS
	if err := services.UploadSpiderToGridFsFromMaster(spider); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

// @Summary Put spider file
// @Description Put spider file
// @Tags spider
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id path string true "spider id"
// @Param reqBody body routes.SpiderFileReqBody true "path"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /spiders/{id}/file [post]
func PutSpiderFile(c *gin.Context) {
	spiderId := c.Param("id")
	if !bson.IsObjectIdHex(spiderId) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
		return
	}
	var reqBody SpiderFileReqBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}
	spider, err := model.GetSpider(bson.ObjectIdHex(spiderId))
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 文件路径
	filePath := path.Join(spider.Src, reqBody.Path)

	// 如果文件已存在，则报错
	if utils.Exists(filePath) {
		HandleErrorF(http.StatusInternalServerError, c, fmt.Sprintf(`%s already exists`, filePath))
		return
	}

	// 写入文件
	if err := ioutil.WriteFile(filePath, []byte(reqBody.Content), 0777); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 同步到GridFS
	if err := services.UploadSpiderToGridFsFromMaster(spider); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

// @Summary Post spider dir
// @Description Post spider dir
// @Tags spider
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id path string true "spider id"
// @Param reqBody body routes.SpiderFileReqBody true "path"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /spiders/{id}/file [put]
func PutSpiderDir(c *gin.Context) {
	spiderId := c.Param("id")
	if !bson.IsObjectIdHex(spiderId) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
		return
	}
	var reqBody SpiderFileReqBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}
	spider, err := model.GetSpider(bson.ObjectIdHex(spiderId))
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 文件路径
	filePath := path.Join(spider.Src, reqBody.Path)

	// 如果文件已存在，则报错
	if utils.Exists(filePath) {
		HandleErrorF(http.StatusInternalServerError, c, fmt.Sprintf(`%s already exists`, filePath))
		return
	}

	// 创建文件夹
	if err := os.MkdirAll(filePath, 0777); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 同步到GridFS
	if err := services.UploadSpiderToGridFsFromMaster(spider); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

// @Summary Delete spider file
// @Description Delete spider file
// @Tags spider
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id path string true "spider id"
// @Param reqBody body routes.SpiderFileReqBody true "path"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /spiders/{id}/file [delete]
func DeleteSpiderFile(c *gin.Context) {
	spiderId := c.Param("id")
	if !bson.IsObjectIdHex(spiderId) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
		return
	}
	var reqBody SpiderFileReqBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}
	spider, err := model.GetSpider(bson.ObjectIdHex(spiderId))
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
	filePath := path.Join(spider.Src, reqBody.Path)
	if err := os.RemoveAll(filePath); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 同步到GridFS
	if err := services.UploadSpiderToGridFsFromMaster(spider); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

// @Summary Rename spider file
// @Description Rename spider file
// @Tags spider
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id path string true "spider id"
// @Param reqBody body routes.SpiderFileReqBody true "path"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /spiders/{id}/file/rename [post]
func RenameSpiderFile(c *gin.Context) {
	spiderId := c.Param("id")

	if !bson.IsObjectIdHex(spiderId) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
		return
	}
	var reqBody SpiderFileReqBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		HandleError(http.StatusBadRequest, c, err)
	}
	spider, err := model.GetSpider(bson.ObjectIdHex(spiderId))
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 原文件路径
	filePath := path.Join(spider.Src, reqBody.Path)
	newFilePath := path.Join(path.Join(path.Dir(filePath), reqBody.NewPath))

	// 如果新文件已存在，则报错
	if utils.Exists(newFilePath) {
		HandleErrorF(http.StatusInternalServerError, c, fmt.Sprintf(`%s already exists`, newFilePath))
		return
	}

	// 重命名
	if err := os.Rename(filePath, newFilePath); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 删除原文件
	if err := os.RemoveAll(filePath); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
	}

	// 同步到GridFS
	if err := services.UploadSpiderToGridFsFromMaster(spider); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

// ======== 爬虫文件管理 ========

// ======== Scrapy 部分 ========

// @Summary Get scrapy spider file
// @Description Get scrapy spider file
// @Tags spider
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id path string true "spider id"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /spiders/{id}/scrapy/spiders [get]
func GetSpiderScrapySpiders(c *gin.Context) {
	id := c.Param("id")

	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "spider_id is invalid")
		return
	}

	spider, err := model.GetSpider(bson.ObjectIdHex(id))
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	spiderNames, err := services.GetScrapySpiderNames(spider)
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    spiderNames,
	})
}

// @Summary Put scrapy spider file
// @Description Put scrapy spider file
// @Tags spider
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id path string true "spider id"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /spiders/{id}/scrapy/spiders [put]
func PutSpiderScrapySpiders(c *gin.Context) {
	type ReqBody struct {
		Name     string `json:"name"`
		Domain   string `json:"domain"`
		Template string `json:"template"`
	}

	id := c.Param("id")
	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
		return
	}
	var reqBody ReqBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		HandleErrorF(http.StatusBadRequest, c, "invalid request")
		return
	}

	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "spider_id is invalid")
		return
	}

	spider, err := model.GetSpider(bson.ObjectIdHex(id))
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	if err := services.CreateScrapySpider(spider, reqBody.Name, reqBody.Domain, reqBody.Template); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

// @Summary Get scrapy spider settings
// @Description Get scrapy spider settings
// @Tags spider
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id path string true "spider id"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /spiders/{id}/scrapy/settings [get]
func GetSpiderScrapySettings(c *gin.Context) {
	id := c.Param("id")

	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "spider_id is invalid")
		return
	}

	spider, err := model.GetSpider(bson.ObjectIdHex(id))
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	data, err := services.GetScrapySettings(spider)
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    data,
	})
}

// @Summary Get scrapy spider file
// @Description Get scrapy spider file
// @Tags spider
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id path string true "spider id"
// @Param reqData body []entity.ScrapySettingParam true "req data"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /spiders/{id}/scrapy/settings [post]
func PostSpiderScrapySettings(c *gin.Context) {
	id := c.Param("id")

	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "spider_id is invalid")
		return
	}

	var reqData []entity.ScrapySettingParam
	if err := c.ShouldBindJSON(&reqData); err != nil {
		HandleErrorF(http.StatusBadRequest, c, "invalid request")
		return
	}

	spider, err := model.GetSpider(bson.ObjectIdHex(id))
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	if err := services.SaveScrapySettings(spider, reqData); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

// @Summary Get scrapy spider items
// @Description Get scrapy spider items
// @Tags spider
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id path string true "spider id"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /spiders/{id}/scrapy/items [get]
func GetSpiderScrapyItems(c *gin.Context) {
	id := c.Param("id")

	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "spider_id is invalid")
		return
	}

	spider, err := model.GetSpider(bson.ObjectIdHex(id))
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	data, err := services.GetScrapyItems(spider)
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    data,
	})
}

// @Summary Post scrapy spider items
// @Description Post scrapy spider items
// @Tags spider
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id path string true "spider id"
// @Param reqData body 	[]entity.ScrapyItem true "req data"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /spiders/{id}/scrapy/items [post]
func PostSpiderScrapyItems(c *gin.Context) {
	id := c.Param("id")

	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "spider_id is invalid")
		return
	}

	var reqData []entity.ScrapyItem
	if err := c.ShouldBindJSON(&reqData); err != nil {
		HandleErrorF(http.StatusBadRequest, c, "invalid request")
		return
	}

	spider, err := model.GetSpider(bson.ObjectIdHex(id))
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	if err := services.SaveScrapyItems(spider, reqData); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

// @Summary Get scrapy spider pipelines
// @Description Get scrapy spider pipelines
// @Tags spider
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id path string true "spider id"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /spiders/{id}/scrapy/pipelines [get]
func GetSpiderScrapyPipelines(c *gin.Context) {
	id := c.Param("id")

	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "spider_id is invalid")
		return
	}

	spider, err := model.GetSpider(bson.ObjectIdHex(id))
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	data, err := services.GetScrapyPipelines(spider)
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    data,
	})
}

// @Summary Get scrapy spider file path
// @Description Get scrapy spider file path
// @Tags spider
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id path string true "spider id"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /spiders/{id}/scrapy/spider/filepath [get]
func GetSpiderScrapySpiderFilepath(c *gin.Context) {
	id := c.Param("id")
	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
		return
	}
	spiderName := c.Query("spider_name")
	if spiderName == "" {
		HandleErrorF(http.StatusBadRequest, c, "spider_name is empty")
		return
	}

	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "spider_id is invalid")
		return
	}

	spider, err := model.GetSpider(bson.ObjectIdHex(id))
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	data, err := services.GetScrapySpiderFilepath(spider, spiderName)
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    data,
	})
}

// ======== ./Scrapy 部分 ========

// ======== Git 部分 ========

// @Summary Post  spider  sync git
// @Description Post  spider  sync git
// @Tags spider
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id path string true "spider id"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /spiders/{id}/git/sync [post]
func PostSpiderSyncGit(c *gin.Context) {
	id := c.Param("id")

	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "spider_id is invalid")
		return
	}

	spider, err := model.GetSpider(bson.ObjectIdHex(id))
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	if err := services.SyncSpiderGit(spider); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

// @Summary Post  spider  reset git
// @Description Post  spider  reset git
// @Tags spider
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id path string true "spider id"
// @Success 200 json string Response
// @Failure 400 json string Response
// @Router /spiders/{id}/git/reset [post]
func PostSpiderResetGit(c *gin.Context) {
	id := c.Param("id")

	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "spider_id is invalid")
		return
	}

	spider, err := model.GetSpider(bson.ObjectIdHex(id))
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	if err := services.ResetSpiderGit(spider); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

// ======== ./Git 部分 ========
