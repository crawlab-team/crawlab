package services

import (
	"crawlab/constants"
	"crawlab/database"
	"crawlab/lib/cron"
	"crawlab/model"
	"crawlab/utils"
	"encoding/json"
	"github.com/apex/log"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime/debug"
)

// 任务日志频道映射
var TaskLogChanMap = utils.NewChanMap()

// 获取本地日志
func GetLocalLog(logPath string) (fileBytes []byte, err error) {

	f, err := os.Open(logPath)
	if err != nil {
		log.Error(err.Error())
		debug.PrintStack()
		return nil, err
	}
	fi, err := f.Stat()
	if err != nil {
		log.Error(err.Error())
		debug.PrintStack()
		return nil, err
	}
	defer f.Close()

	const bufLen = 2 * 1024 * 1024
	logBuf := make([]byte, bufLen)

	off := int64(0)
	if fi.Size() > int64(len(logBuf)) {
		off = fi.Size() - int64(len(logBuf))
	}
	n, err := f.ReadAt(logBuf, off)

	//到文件结尾会有EOF标识
	if err != nil && err.Error() != "EOF" {
		log.Error(err.Error())
		debug.PrintStack()
		return nil, err
	}
	logBuf = logBuf[:n]
	return logBuf, nil
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

func InitDeleteLogPeriodically() error {
	c := cron.New(cron.WithSeconds())
	if _, err := c.AddFunc(viper.GetString("log.deleteFrequency"), DeleteLogPeriodically); err != nil {
		return err
	}

	c.Start()
	return nil

}
