package msg_handler

import (
	"crawlab/constants"
	"crawlab/entity"
	"crawlab/model"
	"crawlab/utils"
)

type SystemInfo struct {
	msg entity.NodeMessage
}

func (s *SystemInfo) Handle() error {
	// 获取环境信息
	sysInfo, err := model.GetLocalSystemInfo()
	if err != nil {
		return err
	}
	msgSd := entity.NodeMessage{
		Type:    constants.MsgTypeGetSystemInfo,
		NodeId:  s.msg.NodeId,
		SysInfo: sysInfo,
	}
	if err := utils.Pub(constants.ChannelMasterNode, msgSd); err != nil {
		return err
	}
	return nil
}
