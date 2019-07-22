package routes

import (
	"crawlab/model"
	"crawlab/services"
	"crawlab/utils"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
)

func GetSpiderList(c *gin.Context) {
	results, err := model.GetSpiderList(nil, 0, 0)
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    results,
	})
}

func GetSpider(c *gin.Context) {
	id := c.Param("id")

	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
	}

	result, err := model.GetSpider(bson.ObjectIdHex(id))
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
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

	if err := model.UpdateSpider(bson.ObjectIdHex(id), item); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

func PublishAllSpiders(c *gin.Context) {
	if err := services.PublishAllSpiders(); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

func PublishSpider(c *gin.Context) {
	id := c.Param("id")

	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
	}

	spider, err := model.GetSpider(bson.ObjectIdHex(id))
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	if err := services.PublishSpider(spider); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

func PutSpider(c *gin.Context) {
	// 从body中获取文件
	file, err := c.FormFile("file")
	if err != nil {
		debug.PrintStack()
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 如果不为zip文件，返回错误
	if !strings.HasSuffix(file.Filename, ".zip") {
		debug.PrintStack()
		HandleError(http.StatusBadRequest, c, errors.New("Not a valid zip file"))
		return
	}

	// 保存到本地临时文件
	randomId := uuid.NewV4()
	tmpFilePath := filepath.Join(viper.GetString("other.tmppath"), randomId.String()+".zip")
	if err := c.SaveUploadedFile(file, tmpFilePath); err != nil {
		debug.PrintStack()
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 读取临时文件
	tmpFile, err := os.OpenFile(tmpFilePath, os.O_RDONLY, 0777)
	if err != nil {
		debug.PrintStack()
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
	if err = tmpFile.Close(); err != nil {
		debug.PrintStack()
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 目标目录
	dstPath := filepath.Join(
		viper.GetString("spider.path"),
		strings.Replace(file.Filename, ".zip", "", 1),
	)

	// 如果目标目录已存在，删除目标目录
	if utils.Exists(dstPath) {
		if err := os.RemoveAll(dstPath); err != nil {
			debug.PrintStack()
			HandleError(http.StatusInternalServerError, c, err)
		}
	}

	// 将临时文件解压到爬虫目录
	if err := utils.DeCompress(tmpFile, dstPath); err != nil {
		debug.PrintStack()
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 删除临时文件
	if err = os.Remove(tmpFilePath); err != nil {
		debug.PrintStack()
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 更新爬虫
	services.UpdateSpiders()

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

func DeleteSpider(c *gin.Context) {
	id := c.Param("id")

	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
		return
	}

	// 获取该爬虫
	spider, err := model.GetSpider(bson.ObjectIdHex(id))
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 删除爬虫文件目录
	if err := os.RemoveAll(spider.Src); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 从数据库中删除该爬虫
	if err := model.RemoveSpider(bson.ObjectIdHex(id)); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

func GetSpiderTasks(c *gin.Context) {
	id := c.Param("id")

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
