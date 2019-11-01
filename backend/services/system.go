package services

import (
	"crawlab/constants"
	"crawlab/database"
	"crawlab/entity"
	"crawlab/model"
	"crawlab/utils"
	"encoding/json"
)

var SystemInfoChanMap = utils.NewChanMap()

func GetRemoteSystemInfo(id string) (sysInfo entity.SystemInfo, err error) {
	// 发送消息
	msg := entity.NodeMessage{
		Type:   constants.MsgTypeGetSystemInfo,
		NodeId: id,
	}

	// 序列化
	msgBytes, _ := json.Marshal(&msg)
	if _, err := database.RedisClient.Publish("nodes:"+id, utils.BytesToString(msgBytes)); err != nil {
		return entity.SystemInfo{}, err
	}

	// 通道
	ch := SystemInfoChanMap.ChanBlocked(id)

	// 等待响应，阻塞
	sysInfoStr := <-ch

	// 反序列化
	if err := json.Unmarshal([]byte(sysInfoStr), &sysInfo); err != nil {
		return sysInfo, err
	}

	return sysInfo, nil
}

func GetSystemInfo(id string) (sysInfo entity.SystemInfo, err error) {
	if IsMasterNode(id) {
		sysInfo, err = model.GetLocalSystemInfo()
	} else {
		sysInfo, err = GetRemoteSystemInfo(id)
	}
	return
}
