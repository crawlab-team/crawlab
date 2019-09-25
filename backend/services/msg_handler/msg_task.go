package msg_handler

import (
	"crawlab/constants"
	"crawlab/model"
	"crawlab/utils"
	"github.com/apex/log"
	"runtime/debug"
	"time"
)

type Task struct {
	msg NodeMessage
}

func (t *Task) Handle() error {
	// 取消任务
	ch := utils.TaskExecChanMap.ChanBlocked(t.msg.TaskId)
	if ch != nil {
		ch <- constants.TaskCancel
	} else {
		// 节点可能被重启，找不到chan
		t, _ := model.GetTask(t.msg.TaskId)
		t.Status = constants.StatusCancelled
		t.FinishTs = time.Now()
		if err := t.Save(); err != nil {
			debug.PrintStack()
			log.Infof("cancel task error: %s", err.Error())
		}
	}
	return nil
}
