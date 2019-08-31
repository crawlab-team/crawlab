package services

import (
	"crawlab/config"
	"github.com/apex/log"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
	"testing"
)

func TestDeleteLogPeriodically(t *testing.T) {
	Convey("Test DeleteLogPeriodically", t, func() {
		if err := config.InitConfig("../conf/config.yml"); err != nil {
			log.Error("init config error:" + err.Error())
			panic(err)
		}
		log.Info("初始化配置成功")
		logDir := viper.GetString("log.path")
		log.Info(logDir)
		DeleteLogPeriodically()
	})
}
