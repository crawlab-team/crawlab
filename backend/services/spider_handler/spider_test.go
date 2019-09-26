package spider_handler

import (
	"crawlab/config"
	"github.com/apex/log"
	"testing"
)

func init() {
	if err := config.InitConfig("../../conf/config.yml"); err != nil {
		log.Fatal("Init config failed")
	}
	log.Infof("初始化配置成功")
}

func TestSpiderSync_CreateMd5File(t *testing.T) {
	s := SpiderSync{}
	s.CreateMd5File("asssss", "gongyu_abc")
}
