package model

import (
	"crawlab/config"
	"crawlab/constants"
	"crawlab/database"
	"github.com/apex/log"
	. "github.com/smartystreets/goconvey/convey"
	"runtime/debug"
	"testing"
)

func TestAddNode(t *testing.T) {
	Convey("Test AddNode", t, func() {
		if err := config.InitConfig("../conf/config.yml"); err != nil {
			log.Error("init config error:" + err.Error())
			panic(err)
		}
		log.Info("初始化配置成功")

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

		var node = Node{
			Key:      "c4:b3:01:bd:b5:e7",
			Name:     "10.27.238.101",
			Ip:       "10.27.238.101",
			Port:     "8000",
			Mac:      "c4:b3:01:bd:b5:e7",
			Status:   constants.StatusOnline,
			IsMaster: true,
		}
		if err := node.Add(); err != nil {
			log.Error("add node error:" + err.Error())
			panic(err)
		}
	})
}
