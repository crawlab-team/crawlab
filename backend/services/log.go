package services

import (
	"crawlab/constants"
	"crawlab/database"
	"crawlab/entity"
	"crawlab/lib/cron"
	"crawlab/model"
	"crawlab/utils"
	"encoding/json"
	"github.com/apex/log"
	"github.com/globalsign/mgo/bson"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"
)

// 任务日志频道映射
var TaskLogChanMap = utils.NewChanMap()

// 获取远端日志
func GetRemoteLog(task model.Task) (logStr string, err error) {
	// 序列化消息
	msg := entity.NodeMessage{
		Type:    constants.MsgTypeGetLog,
		LogPath: task.LogPath,
		TaskId:  task.Id,
	}
	msgBytes, err := json.Marshal(&msg)
	if err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return "", err
	}

	// 发布获取日志消息
	channel := "nodes:" + task.NodeId.Hex()
	if _, err := database.RedisClient.Publish(channel, utils.BytesToString(msgBytes)); err != nil {
		log.Errorf(err.Error())
		return "", err
	}

	// 生成频道，等待获取log
	ch := TaskLogChanMap.ChanBlocked(task.Id)

	select {
	case logStr = <-ch:
		log.Infof("get remote log")
		break
	case <-time.After(30 * time.Second):
		logStr = "get remote log timeout"
		break
	}

	return logStr, nil
}

// 定时删除日志
func DeleteLogPeriodically() {
	logDir := viper.GetString("log.path")
	if !utils.Exists(logDir) {
		log.Error("Can Not Set Delete Logs Periodically,No Log Dir")
		return
	}
	rd, err := ioutil.ReadDir(logDir)
	if err != nil {
		log.Error("Read Log Dir Failed")
		return
	}

	for _, fi := range rd {
		if fi.IsDir() {
			log.Info(filepath.Join(logDir, fi.Name()))
			os.RemoveAll(filepath.Join(logDir, fi.Name()))
			log.Info("Delete Log File Success")
		}
	}

}

// 删除本地日志
func RemoveLocalLog(path string) error {
	if err := model.RemoveFile(path); err != nil {
		log.Error("remove local file error: " + err.Error())
		return err
	}
	return nil
}

// 删除远程日志
func RemoveRemoteLog(task model.Task) error {
	msg := entity.NodeMessage{
		Type:    constants.MsgTypeRemoveLog,
		LogPath: task.LogPath,
		TaskId:  task.Id,
	}
	// 发布获取日志消息
	channel := "nodes:" + task.NodeId.Hex()
	if _, err := database.RedisClient.Publish(channel, utils.GetJson(msg)); err != nil {
		log.Errorf("publish redis error: %s", err.Error())
		debug.PrintStack()
		return err
	}
	return nil
}

// 删除日志文件
func RemoveLogByTaskId(id string) error {
	t, err := model.GetTask(id)
	if err != nil {
		log.Error("get task error:" + err.Error())
		return err
	}
	removeLog(t)

	return nil
}

func removeLog(t model.Task) {
	if err := RemoveLocalLog(t.LogPath); err != nil {
		log.Errorf("remove local log error: %s", err.Error())
		debug.PrintStack()
	}
	if err := RemoveRemoteLog(t); err != nil {
		log.Errorf("remove remote log error: %s", err.Error())
		debug.PrintStack()
	}
}

// 删除日志文件
func RemoveLogBySpiderId(id bson.ObjectId) error {
	tasks, err := model.GetTaskList(bson.M{"spider_id": id}, 0, constants.Infinite, "-create_ts")
	if err != nil {
		log.Errorf("get tasks error: %s", err.Error())
		debug.PrintStack()
	}
	for _, task := range tasks {
		removeLog(task)
	}
	return nil
}

// 初始化定时删除日志
func InitDeleteLogPeriodically() error {
	c := cron.New(cron.WithSeconds())
	if _, err := c.AddFunc(viper.GetString("log.deleteFrequency"), DeleteLogPeriodically); err != nil {
		return err
	}

	c.Start()
	return nil

}
