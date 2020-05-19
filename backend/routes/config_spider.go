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
	"strings"
)

// 添加可配置爬虫

// @Summary Put config spider
// @Description Put config spider
// @Tags config spider
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param spider body  model.Spider true "spider item"
// @Success 200 json string Response
// @Failure 500 json string Response
// @Router /config_spiders [put]
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

	// 模版名不能为空
	if spider.Template == "" {
		HandleErrorF(http.StatusBadRequest, c, "spider template should not be empty")
		return
	}

	// 判断爬虫是否存在
	if spider := model.GetSpiderByName(spider.Name); spider.Name != "" {
		HandleErrorF(http.StatusBadRequest, c, fmt.Sprintf("spider for '%s' already exists", spider.Name))
		return
	}

	// 设置爬虫类别
	spider.Type = constants.Configurable

	// 将FileId置空
	spider.FileId = bson.ObjectIdHex(constants.ObjectIdNull)

	// UserId
	spider.UserId = services.GetCurrentUserId(c)

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

	// 复制Spiderfile模版
	contentByte, err := ioutil.ReadFile("./template/spiderfile/Spiderfile." + spider.Template)
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
	f, err := os.Create(filepath.Join(spider.Src, "Spiderfile"))
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
	defer f.Close()
	if _, err := f.Write(contentByte); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

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

// @Summary Upload config spider
// @Description Upload config spider
// @Tags config spider
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param spider body  model.Spider true "spider item"
// @Param id path string true "spider id"
// @Success 200 json string Response
// @Failure 500 json string Response
// @Router /config_spiders/{id}/upload [post]
func UploadConfigSpider(c *gin.Context) {
	id := c.Param("id")

	// 获取爬虫
	var spider model.Spider
	spider, err := model.GetSpider(bson.ObjectIdHex(id))
	if err != nil {
		HandleErrorF(http.StatusBadRequest, c, fmt.Sprintf("cannot find spider (id: %s)", id))
		return
	}

	// UserId
	spider.UserId = services.GetCurrentUserId(c)

	// 获取上传文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	// 文件名称必须为Spiderfile
	filename := header.Filename
	if filename != "Spiderfile" && filename != "Spiderfile.yaml" && filename != "Spiderfile.yml" {
		HandleErrorF(http.StatusBadRequest, c, "filename must be 'Spiderfile(.yaml|.yml)'")
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
			return
		}
	} else {
		f, err = os.Create(sfPath)
		if err != nil {
			HandleError(http.StatusInternalServerError, c, err)
			return
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

	// 根据序列化后的数据处理爬虫文件
	if err := services.ProcessSpiderFilesFromConfigData(spider, configData); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

// @Summary Post config spider
// @Description Post config spider
// @Tags config spider
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id path string true "spider id"
// @Success 200 json string Response
// @Failure 500 json string Response
// @Router /config_spiders/{id}/spiderfile [post]
func PostConfigSpiderSpiderfile(c *gin.Context) {
	type Body struct {
		Content string `json:"content"`
	}

	id := c.Param("id")

	// 文件内容
	var reqBody Body
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}
	content := reqBody.Content

	// 获取爬虫
	var spider model.Spider
	spider, err := model.GetSpider(bson.ObjectIdHex(id))
	if err != nil {
		HandleErrorF(http.StatusBadRequest, c, fmt.Sprintf("cannot find spider (id: %s)", id))
		return
	}

	// UserId
	if !spider.UserId.Valid() {
		spider.UserId = bson.ObjectIdHex(constants.ObjectIdNull)
	}

	// 反序列化
	var configData entity.ConfigSpiderData
	if err := yaml.Unmarshal([]byte(content), &configData); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	// 校验configData
	if err := services.ValidateSpiderfile(configData); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 写文件
	if err := ioutil.WriteFile(filepath.Join(spider.Src, "Spiderfile"), []byte(content), os.ModePerm); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 根据序列化后的数据处理爬虫文件
	if err := services.ProcessSpiderFilesFromConfigData(spider, configData); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

// @Summary Post config spider config
// @Description Post config spider config
// @Tags config spider
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param spider body  model.Spider true "spider item"
// @Param id path string true "spider id"
// @Success 200 json string Response
// @Failure 500 json string Response
// @Router /config_spiders/{id}/config [post]
func PostConfigSpiderConfig(c *gin.Context) {
	id := c.Param("id")

	// 获取爬虫
	var spider model.Spider
	spider, err := model.GetSpider(bson.ObjectIdHex(id))
	if err != nil {
		HandleErrorF(http.StatusBadRequest, c, fmt.Sprintf("cannot find spider (id: %s)", id))
		return
	}

	// UserId
	if !spider.UserId.Valid() {
		spider.UserId = bson.ObjectIdHex(constants.ObjectIdNull)
	}

	// 反序列化配置数据
	var configData entity.ConfigSpiderData
	if err := c.ShouldBindJSON(&configData); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	// 校验configData
	if err := services.ValidateSpiderfile(configData); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 替换Spiderfile文件
	if err := services.GenerateSpiderfileFromConfigData(spider, configData); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 根据序列化后的数据处理爬虫文件
	if err := services.ProcessSpiderFilesFromConfigData(spider, configData); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

// @Summary Get config spider
// @Description Get config spider
// @Tags config spider
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id path string true "spider id"
// @Success 200 json string Response
// @Failure 500 json string Response
// @Router /config_spiders/{id}/config [get]
func GetConfigSpiderConfig(c *gin.Context) {
	id := c.Param("id")

	// 校验ID
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

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    spider.Config,
	})
}

// 获取模版名称列表

// @Summary Get config spider template list
// @Description Get config spider template list
// @Tags config spider
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id path string true "spider id"
// @Success 200 json string Response
// @Failure 500 json string Response
// @Router /config_spiders_templates [get]
func GetConfigSpiderTemplateList(c *gin.Context) {
	var data []string
	for _, fInfo := range utils.ListDir("./template/spiderfile") {
		templateName := strings.Replace(fInfo.Name(), "Spiderfile.", "", -1)
		data = append(data, templateName)
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    data,
	})
}
