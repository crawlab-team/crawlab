package services

import (
	"crawlab/config"
	"crawlab/utils"
	"fmt"
	"github.com/apex/log"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
	"os"
	"testing"
)

func TestDeleteLogPeriodically(t *testing.T) {
	Convey("Test DeleteLogPeriodically", t, func() {
		err := config.InitConfig("../conf/config.yml")
		So(err, ShouldBeNil)
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
	defer utils.Close(f)
	if err != nil {
		fmt.Println(err)

	} else {
		_, err = f.WriteString("This is for test")
		fmt.Println(err)
	}

	//delete the test log file
	_ = os.Remove(logPath)

}
