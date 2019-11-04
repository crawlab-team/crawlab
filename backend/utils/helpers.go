package utils

import (
	"crawlab/entity"
	"encoding/json"
	"github.com/apex/log"
	"github.com/gomodule/redigo/redis"
	"io"
	"runtime/debug"
	"unsafe"
)

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func GetJson(message entity.NodeMessage) string {
	msgBytes, err := json.Marshal(&message)
	if err != nil {
		log.Errorf("node message to json error: %s", err.Error())
		debug.PrintStack()
		return ""
	}
	return BytesToString(msgBytes)
}

func GetMessage(message redis.Message) *entity.NodeMessage {
	msg := entity.NodeMessage{}
	if err := json.Unmarshal(message.Data, &msg); err != nil {
		log.Errorf("message byte to object error: %s", err.Error())
		debug.PrintStack()
		return nil
	}
	return &msg
}

func Close(c io.Closer) {
	err := c.Close()
	if err != nil {
		//log.WithError(err).Error("关闭资源文件失败。")
	}
}
