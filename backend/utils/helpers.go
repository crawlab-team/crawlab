package utils

import (
	"context"
	"crawlab/database"
	"crawlab/entity"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/apex/log"
	"github.com/gomodule/redigo/redis"
	"math/rand"
	"runtime/debug"
	"time"
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

	return SubWithContext(ctx, channel, consume)
}

func SubWithContext(ctx context.Context, channel string, consume database.ConsumeFunc) error {
	if err := database.RedisClient.Subscribe(ctx, consume, channel); err != nil {
		log.Errorf("subscribe redis error: %s", err.Error())
		debug.PrintStack()
		return err
	}
	return nil
}

//生成随机字符串
func RandomString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// 生成32位MD5
func MD5(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

//第一版的加密方式
func EncryptPasswordV1(password string) string {
	return MD5(password)
}
//第二版的加密方式
func EncryptPasswordV2(password, salt string) string {
	return MD5(salt + MD5(password+salt))
}
