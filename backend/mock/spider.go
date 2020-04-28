package mock

import (
	"crawlab/constants"
	"crawlab/model"
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var SpiderList = []model.Spider{
	{
		Id:          bson.ObjectId("5d429e6c19f7abede924fee2"),
		Name:        "For test",
		DisplayName: "test",
		Type:        "test",
		Col:         "test",
		Site:        "www.baidu.com",
		Envs:        nil,
		Src:         "../app/spiders",
		Cmd:         "scrapy crawl test",
		LastRunTs:   time.Now(),
		CreateTs:    time.Now(),
		UpdateTs:    time.Now(),
		UserId:      constants.ObjectIdNull,
	},
}

func GetSpiderList(c *gin.Context) {

	// mock get spider list from database
	results := SpiderList

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    results,
	})
}

func GetSpider(c *gin.Context) {
	id := c.Param("id")
	var result model.Spider

	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
	}

	for _, spider := range SpiderList {
		if spider.Id == bson.ObjectId(id) {
			result = spider
		}
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    result,
	})
}

func PostSpider(c *gin.Context) {
	id := c.Param("id")

	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
	}

	var item model.Spider
	if err := c.ShouldBindJSON(&item); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	log.Info("modify the item")

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}
func GetSpiderDir(c *gin.Context) {
	// 爬虫ID
	id := c.Param("id")

	// 目录相对路径
	path := c.Query("path")
	var spi model.Spider

	// 获取爬虫
	for _, spider := range SpiderList {
		if spider.Id == bson.ObjectId(id) {
			spi = spider
		}
	}

	// 获取目录下文件列表
	f, err := ioutil.ReadDir(filepath.Join(spi.Src, path))
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

func GetSpiderTasks(c *gin.Context) {
	id := c.Param("id")

	var spider model.Spider
	for _, spi := range SpiderList {
		if spi.Id == bson.ObjectId(id) {
			spider = spi
		}
	}

	var tasks model.Task
	for _, task := range TaskList {
		if task.SpiderId == spider.Id {
			tasks = task
		}
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    tasks,
	})
}

func DeleteSpider(c *gin.Context) {
	id := c.Param("id")

	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
		return
	}

	// 获取该爬虫,get this spider
	var spider model.Spider
	for _, spi := range SpiderList {
		if spi.Id == bson.ObjectId(id) {
			spider = spi
		}
	}

	// 删除爬虫文件目录,delete the spider dir
	if err := os.RemoveAll(spider.Src); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 从数据库中删除该爬虫,delete this spider from database

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}
