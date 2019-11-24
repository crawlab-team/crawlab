package routes

import (
	"crawlab/constants"
	"crawlab/entity"
	"crawlab/model"
	"crawlab/services"
	"crawlab/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
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

	// 创建爬虫目录
	spiderDir := filepath.Join(viper.GetString("spider.path"), spider.Name)
	if utils.Exists(spiderDir) {
		if err := os.RemoveAll(spiderDir); err != nil {
			HandleError(http.StatusInternalServerError, c, err)
			return
		}
	}
	if err := os.MkdirAll(spiderDir, 0777); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
	spider.Src = spiderDir

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

// 上传可配置爬虫Spiderfile
func UploadConfigSpider(c *gin.Context) {
	id := c.Param("id")

	// 获取爬虫
	var spider model.Spider
	spider, err := model.GetSpider(bson.ObjectIdHex(id))
	if err != nil {
		HandleErrorF(http.StatusBadRequest, c, fmt.Sprintf("cannot find spider (id: %s)", id))
	}

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

	// 爬虫目录
	spiderDir := filepath.Join(viper.GetString("spider.path"), spider.Name)

	// 爬虫Spiderfile文件路径
	sfPath := filepath.Join(spiderDir, filename)

	// 创建（如果不存在）或打开Spiderfile（如果存在）
	var f *os.File
	if utils.Exists(sfPath) {
		f, err = os.OpenFile(sfPath, os.O_WRONLY, 0777)
		if err != nil {
			HandleError(http.StatusInternalServerError, c, err)
		}
	} else {
		f, err = os.Create(sfPath)
		if err != nil {
			HandleError(http.StatusInternalServerError, c, err)
		}
	}

	// 将上传的文件拷贝到爬虫Spiderfile文件
	_, err = io.Copy(f, file)
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 关闭Spiderfile文件
	_ = f.Close()

	// 构造配置数据
	configData := entity.ConfigSpiderData{}

	// 读取YAML文件
	yamlFile, err := ioutil.ReadFile(sfPath)
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 反序列化
	if err := yaml.Unmarshal(yamlFile, &configData); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 删除已有的爬虫文件
	for _, fInfo := range utils.ListDir(spiderDir) {
		// 不删除Spiderfile
		if fInfo.Name() == filename {
			continue
		}

		// 删除其他文件
		if err := os.RemoveAll(filepath.Join(spiderDir, fInfo.Name())); err != nil {
			HandleError(http.StatusInternalServerError, c, err)
			return
		}
	}

	// 拷贝爬虫文件
	tplDir := "./template/scrapy"
	for _, fInfo := range utils.ListDir(tplDir) {
		// 跳过Spiderfile
		if fInfo.Name() == "Spiderfile" {
			continue
		}

		srcPath := filepath.Join(tplDir, fInfo.Name())
		if fInfo.IsDir() {
			dirPath := filepath.Join(spiderDir, fInfo.Name())
			if err := utils.CopyDir(srcPath, dirPath); err != nil {
				HandleError(http.StatusInternalServerError, c, err)
				return
			}
		} else {
			if err := utils.CopyFile(srcPath, filepath.Join(spiderDir, fInfo.Name())); err != nil {
				HandleError(http.StatusInternalServerError, c, err)
				return
			}
		}
	}

	// 更改爬虫文件
	if err := services.GenerateConfigSpiderFiles(spider, configData); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// TODO: 上传到GridFS

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    configData,
	})
}
