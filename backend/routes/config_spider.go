package routes

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
)

type Field struct {
	Name  string `yaml:"name" json:"name"`
	Css   string `yaml:"css" json:"css"`
	Xpath string `yaml:"xpath" json:"xpath"`
	Attr  string `yaml:"attr" json:"attr"`
	Stage string `yaml:"stage" json:"stage"`
}

type Stage struct {
	List   bool    `yaml:"list" json:"list"`
	Css    string  `yaml:"css" json:"css"`
	Xpath  string  `yaml:"xpath" json:"xpath"`
	Fields []Field `yaml:"fields" json:"fields"`
}

type ConfigSpiderData struct {
	Version  string           `yaml:"version" json:"version"`
	StartUrl string           `yaml:"startUrl" json:"start_url"`
	Stages   map[string]Stage `yaml:"stages" json:"stages"`
}

func PutConfigSpider(c *gin.Context) {
	// 构造配置数据
	data := ConfigSpiderData{}

	// 读取YAML文件
	yamlFile, err := ioutil.ReadFile("./template/Spiderfile")
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
