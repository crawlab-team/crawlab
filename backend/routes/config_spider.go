package routes

import (
	"crawlab/constants"
	"crawlab/entity"
	"crawlab/model"
	"crawlab/utils"
	"fmt"
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
)

// 添加可配置爬虫
func PutConfigSpider(c *gin.Context) {
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
	if spider := model.GetSpiderByName(spider.Name); spider != nil {
		HandleErrorF(http.StatusBadRequest, c, fmt.Sprintf("spider for '%s' already exists", spider.Name))
		return
	}

	// 设置爬虫类别
	spider.Type = constants.Configurable

	// 将FileId置空
	spider.FileId = bson.ObjectIdHex(constants.ObjectIdNull)

	// 添加爬虫到数据库
	if err := spider.Add(); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    spider,
	})
}

// 更改可配置爬虫
func PostConfigSpider(c *gin.Context) {
	PostSpider(c)
}

func UploadConfigSpider(c *gin.Context) {
	// 获取上传文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	// 文件名称必须为Spiderfile
	filename := header.Filename
	if filename != "Spiderfile" {
		HandleErrorF(http.StatusBadRequest, c, "filename must be 'Spiderfile'")
		return
	}

	// 以防tmp目录不存在
	tmpPath := viper.GetString("other.tmppath")
	if !utils.Exists(tmpPath) {
		if err := os.MkdirAll(tmpPath, os.ModePerm); err != nil {
			log.Error("mkdir other.tmppath dir error:" + err.Error())
			debug.PrintStack()
			HandleError(http.StatusBadRequest, c, err)
			return
		}
	}

	//创建文件
	randomId := uuid.NewV4()
	tmpFilePath := filepath.Join(tmpPath, "Spiderfile."+randomId.String())
	out, err := os.Create(tmpFilePath)
	if err != nil {
	}
	_, err = io.Copy(out, file)
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
	}
	_ = out.Close()

	// 构造配置数据
	data := entity.ConfigSpiderData{}

	// 读取YAML文件
	yamlFile, err := ioutil.ReadFile(tmpFilePath)
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 反序列化
	if err := yaml.Unmarshal(yamlFile, &data); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// TODO: 生成爬虫文件

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    data,
	})
}
