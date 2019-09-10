package msg_handler

import (
	"crawlab/constants"
	"crawlab/database"
	"crawlab/model"
	"crawlab/utils"
	"encoding/json"
	"github.com/apex/log"
	"runtime/debug"
)

type SystemInfo struct {
	msg NodeMessage
}

func (s *SystemInfo) Handle() error {
	// 获取环境信息
	sysInfo, err := model.GetLocalSystemInfo()
	if err != nil {
		return err
	}
	msgSd := NodeMessage{
		Type:    constants.MsgTypeGetSystemInfo,
		NodeId:  s.msg.NodeId,
		SysInfo: sysInfo,
	}
	msgSdBytes, err := json.Marshal(&msgSd)
	if err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return err
	}
	if _, err := database.RedisClient.Publish("nodes:master", utils.BytesToString(msgSdBytes)); err != nil {
		log.Errorf(err.Error())
		return err
	}
	return nil
}
