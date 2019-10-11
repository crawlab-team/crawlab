package utils

import (
	"context"
	"crawlab/database"
	"crawlab/entity"
	"encoding/json"
	"github.com/apex/log"
	"github.com/gomodule/redigo/redis"
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

func Pub(channel string, msg entity.NodeMessage) error {
	if _, err := database.RedisClient.Publish(channel, GetJson(msg)); err != nil {
		log.Errorf("publish redis error: %s", err.Error())
		debug.PrintStack()
		return err
	}
	return nil
}

func Sub(channel string, consume database.ConsumeFunc) error {
	ctx := context.Background()
	if err := database.RedisClient.Subscribe(ctx, consume, channel); err != nil {
		log.Errorf("subscribe redis error: %s", err.Error())
		debug.PrintStack()
		return err
	}
	return nil
}
