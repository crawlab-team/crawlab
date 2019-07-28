package services

import (
	"crawlab/constants"
	"crawlab/database"
	"crawlab/model"
	"crawlab/utils"
	"encoding/json"
	"github.com/apex/log"
	"io/ioutil"
	"runtime/debug"
)

// 任务日志频道映射
var TaskLogChanMap = utils.NewChanMap()

// 获取本地日志
func GetLocalLog(logPath string) (fileBytes []byte, err error) {
	fileBytes, err = ioutil.ReadFile(logPath)
	if err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return fileBytes, err
	}
	return fileBytes, nil
}

// 获取远端日志
func GetRemoteLog(task model.Task) (logStr string, err error) {
	// 序列化消息
	msg := NodeMessage{
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
	if err := database.Publish(channel, string(msgBytes)); err != nil {
		log.Errorf(err.Error())
		return "", err
	}

	// 生成频道，等待获取log
	ch := TaskLogChanMap.ChanBlocked(task.Id)

	// 此处阻塞，等待结果
	logStr = <-ch

	return logStr, nil
}
