package msg_handler

import (
	"crawlab/model"
	"crawlab/utils"
	"github.com/globalsign/mgo/bson"
	"github.com/spf13/viper"
	"path/filepath"
)

type Spider struct {
	SpiderId string
}

func (s *Spider) Handle() error {
	// 移除本地的爬虫目录
	spider, err := model.GetSpider(bson.ObjectIdHex(s.SpiderId))
	if err != nil {
		return err
	}
	path := filepath.Join(viper.GetString("spider.path"), spider.Name)
	utils.RemoveFiles(path)
	return nil
}
