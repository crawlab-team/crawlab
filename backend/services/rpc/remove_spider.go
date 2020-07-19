package rpc

import (
	"crawlab/constants"
	"crawlab/entity"
	"crawlab/model"
	"crawlab/utils"
	"github.com/globalsign/mgo/bson"
	"github.com/spf13/viper"
	"path/filepath"
)

type RemoveSpiderService struct {
	msg entity.RpcMessage
}

func (s *RemoveSpiderService) ServerHandle() (entity.RpcMessage, error) {
	spiderId := utils.GetRpcParam("spider_id", s.msg.Params)
	if err := RemoveSpiderServiceLocal(spiderId); err != nil {
		s.msg.Error = err.Error()
		return s.msg, err
	}
	s.msg.Result = "success"
	return s.msg, nil
}

func (s *RemoveSpiderService) ClientHandle() (o interface{}, err error) {
	// 发起 RPC 请求，获取服务端数据
	_, err = ClientFunc(s.msg)()
	if err != nil {
		return
	}

	return
}

func RemoveSpiderServiceLocal(spiderId string) error {
	// 移除本地的爬虫目录
	spider, err := model.GetSpider(bson.ObjectIdHex(spiderId))
	if err != nil {
		return err
	}
	path := filepath.Join(viper.GetString("spider.path"), spider.Name)
	utils.RemoveFiles(path)
	return nil
}

func RemoveSpiderServiceRemote(spiderId string, nodeId string) (err error) {
	params := make(map[string]string)
	params["spider_id"] = spiderId
	s := GetService(entity.RpcMessage{
		NodeId:  nodeId,
		Method:  constants.RpcRemoveSpider,
		Params:  params,
		Timeout: 60,
	})
	_, err = s.ClientHandle()
	if err != nil {
		return
	}
	return
}
