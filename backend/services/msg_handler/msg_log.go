package msg_handler

import (
	"crawlab/constants"
	"crawlab/database"
	"crawlab/model"
	"crawlab/utils"
	"encoding/json"
	"github.com/apex/log"
	"runtime/debug"
)

type Log struct {
	msg NodeMessage
}

func (g *Log) Handle() error {
	if g.msg.Type == constants.MsgTypeGetLog {
		return g.get()
	} else if g.msg.Type == constants.MsgTypeRemoveLog {
		return g.remove()
	}
	return nil
}

func (g *Log) get() error {
	// 发出的消息
	msgSd := NodeMessage{
		Type:   constants.MsgTypeGetLog,
		TaskId: g.msg.TaskId,
	}
	// 获取本地日志
	logStr, err := model.GetLocalLog(g.msg.LogPath)
	log.Info(utils.BytesToString(logStr))
	if err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		msgSd.Error = err.Error()
		msgSd.Log = err.Error()
	} else {
		msgSd.Log = utils.BytesToString(logStr)
	}

	// 序列化
	msgSdBytes, err := json.Marshal(&msgSd)
	if err != nil {
		return err
	}

	// 发布消息给主节点
	log.Info("publish get log msg to master")
	if _, err := database.RedisClient.Publish("nodes:master", utils.BytesToString(msgSdBytes)); err != nil {
		return err
	}
	return nil
}

func (g *Log) remove() error {
	return model.RemoveFile(g.msg.LogPath)
}
