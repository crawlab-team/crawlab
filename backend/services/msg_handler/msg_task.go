package msg_handler

import (
	"crawlab/constants"
	"crawlab/utils"
)

type Task struct {
	msg NodeMessage
}

func (t *Task) Handle() error {
	// 取消任务
	ch := utils.TaskExecChanMap.ChanBlocked(t.msg.TaskId)
	ch <- constants.TaskCancel
	return nil
}
