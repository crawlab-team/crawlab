package msg_handler

import (
	"crawlab/constants"
	"crawlab/entity"
	"crawlab/model"
	"crawlab/utils"
	"github.com/apex/log"
	"runtime/debug"
	"time"
)

type Task struct {
	msg entity.NodeMessage
}

func (t *Task) Handle() error {
	log.Infof("received cancel task msg, task_id: %s", t.msg.TaskId)
	// 取消任务
	ch := utils.TaskExecChanMap.ChanBlocked(t.msg.TaskId)
	if ch != nil {
		ch <- constants.TaskCancel
	} else {
		log.Infof("chan is empty, update status to abnormal")
		// 节点可能被重启，找不到chan
		task, err := model.GetTask(t.msg.TaskId)
		if err != nil {
			log.Errorf("task not found, task_id: %s", t.msg.TaskId)
			debug.PrintStack()
			return err
		}
		task.Status = constants.StatusAbnormal
		task.FinishTs = time.Now()
		if err := task.Save(); err != nil {
			debug.PrintStack()
			log.Infof("cancel task error: %s", err.Error())
		}
	}
	return nil
}
