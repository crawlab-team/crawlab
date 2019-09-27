package spider_handler

import (
	"crawlab/config"
	"crawlab/database"
	"crawlab/model"
	"github.com/apex/log"
	"github.com/globalsign/mgo/bson"
	"runtime/debug"
	"testing"
)

var s SpiderSync

func init() {
	if err := config.InitConfig("../../conf/config.yml"); err != nil {
		log.Fatal("Init config failed")
	}
	log.Infof("初始化配置成功")

	// 初始化Mongodb数据库
	if err := database.InitMongo(); err != nil {
		log.Error("init mongodb error:" + err.Error())
		debug.PrintStack()
		panic(err)
	}
	log.Info("初始化Mongodb数据库成功")

	// 初始化Redis数据库
	if err := database.InitRedis(); err != nil {
		log.Error("init redis error:" + err.Error())
		debug.PrintStack()
		panic(err)
	}
	log.Info("初始化Redis数据库成功")

	s = SpiderSync{
		Spider: model.Spider{
			Id:     bson.ObjectIdHex("5d8d5e4b44500b000150009c"),
			Name:   "scrapy-pre_sale",
			FileId: bson.ObjectIdHex("5d8d5e4b44500b0001500098"),
			Src:    "/opt/crawlab/spiders/scrapy-pre_sale",
		},
	}
}

func TestSpiderSync_CreateMd5File(t *testing.T) {
	s.CreateMd5File("this is md5")
}

func TestSpiderSync_Download(t *testing.T) {
	s.Download()
}
