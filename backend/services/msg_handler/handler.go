package msg_handler

import (
	"crawlab/constants"
	"crawlab/entity"
)

type Handler interface {
	Handle() error
}

func GetMsgHandler(msg entity.NodeMessage) Handler {
	if msg.Type == constants.MsgTypeGetLog || msg.Type == constants.MsgTypeRemoveLog {
		// 日志相关
		return &Log{
			msg: msg,
		}
	} else if msg.Type == constants.MsgTypeCancelTask {
		// 任务相关
		return &Task{
			msg: msg,
		}
	} else if msg.Type == constants.MsgTypeGetSystemInfo {
		// 系统信息相关
		return &SystemInfo{
			msg: msg,
		}
	} else if msg.Type == constants.MsgTypeRemoveSpider {
		// 爬虫相关
		return &Spider{
			SpiderId: msg.SpiderId,
		}
	}
	return nil
}
