package msg_handler

import (
	"crawlab/constants"
	"crawlab/model"
)

type Handler interface {
	Handle() error
}

func GetMsgHandler(msg NodeMessage) Handler {
	if msg.Type == constants.MsgTypeGetLog {
		return &Log{
			msg: msg,
		}
	} else if msg.Type == constants.MsgTypeCancelTask {
		return &Task{
			msg: msg,
		}
	} else if msg.Type == constants.MsgTypeGetSystemInfo {
		return &SystemInfo{
			msg: msg,
		}
	}
	return nil
}

type NodeMessage struct {
	// 通信类别
	Type string `json:"type"`

	// 任务相关
	TaskId string `json:"task_id"` // 任务ID

	// 节点相关
	NodeId string `json:"node_id"` // 节点ID

	// 日志相关
	LogPath string `json:"log_path"` // 日志路径
	Log     string `json:"log"`      // 日志

	// 系统信息
	SysInfo model.SystemInfo `json:"sys_info"`

	// 错误相关
	Error string `json:"error"`
}
