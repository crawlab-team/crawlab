package services

import (
	"crawlab/config"
	"fmt"
	"github.com/apex/log"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
	"os"
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

func TestGetLocalLog(t *testing.T) {
	//create a log file for test
	logPath := "../logs/crawlab/test.log"
	f, err := os.Create(logPath)
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())

	} else {
		_, err = f.Write([]byte("This is for test"))
	}

	Convey("Test GetLocalLog", t, func() {
		Convey("Test response", func() {
			logStr, err := GetLocalLog(logPath)
			log.Info(string(logStr))
			fmt.Println(err)
			So(err, ShouldEqual, nil)

		})
	})
	//delete the test log file
	os.Remove(logPath)

}
