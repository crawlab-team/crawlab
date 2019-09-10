package services

import (
	"crawlab/constants"
	"crawlab/database"
	"crawlab/lib/cron"
	"crawlab/model"
	"crawlab/services/msg_handler"
	"crawlab/utils"
	"encoding/json"
	"github.com/apex/log"
	"github.com/globalsign/mgo/bson"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime/debug"
)

// 任务日志频道映射
var TaskLogChanMap = utils.NewChanMap()

// 获取远端日志
func GetRemoteLog(task model.Task) (logStr string, err error) {
	// 序列化消息
	msg := msg_handler.NodeMessage{
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

	// 此处阻塞，等待结果
	logStr = <-ch

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
	msg := msg_handler.NodeMessage{
		Type:    constants.MsgTypeRemoveLog,
		LogPath: task.LogPath,
		TaskId:  task.Id,
	}
	msgBytes, err := json.Marshal(&msg)
	if err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return err
	}
	// 发布获取日志消息
	channel := "nodes:" + task.NodeId.Hex()
	if _, err := database.RedisClient.Publish(channel, utils.BytesToString(msgBytes)); err != nil {
		log.Errorf(err.Error())
		return err
	}
	return nil
}

// 删除日志文件
func RemoveLogBySpiderId(id bson.ObjectId) error {
	tasks, err := model.GetTaskList(bson.M{"spider_id": id}, 0, constants.Infinite, "-create_ts")
	if err != nil {
		log.Error("get tasks error:" + err.Error())
	}
	for _, task := range tasks {
		if err := RemoveLocalLog(task.LogPath); err != nil {
			log.Error("remove local log error:" + err.Error())
		}
		if err := RemoveRemoteLog(task); err != nil {
			log.Error("remove remote log error:" + err.Error())
		}
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
