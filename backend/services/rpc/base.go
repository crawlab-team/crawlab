package rpc

import (
	"crawlab/constants"
	"crawlab/database"
	"crawlab/entity"
	"crawlab/model"
	"crawlab/utils"
	"encoding/json"
	"fmt"
	"github.com/apex/log"
	"github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"runtime/debug"
)

type Service interface {
	ServerHandle() (entity.RpcMessage, error)
	ClientHandle() (interface{}, error)
}

func ClientFunc(msg entity.RpcMessage) func() (entity.RpcMessage, error) {
	return func() (replyMsg entity.RpcMessage, err error) {
		// 请求ID
		msg.Id = uuid.NewV4().String()

		// 发送RPC消息
		msgStr := utils.ObjectToString(msg)
		if err := database.RedisClient.LPush(fmt.Sprintf("rpc:%s", msg.NodeId), msgStr); err != nil {
			log.Errorf("RpcClientFunc error: " + err.Error())
			debug.PrintStack()
			return replyMsg, err
		}

		// 获取RPC回复消息
		dataStr, err := database.RedisClient.BRPop(fmt.Sprintf("rpc:%s:%s", msg.NodeId, msg.Id), msg.Timeout)
		if err != nil {
			log.Errorf("RpcClientFunc error: " + err.Error())
			debug.PrintStack()
			return replyMsg, err
		}

		// 反序列化消息
		if err := json.Unmarshal([]byte(dataStr), &replyMsg); err != nil {
			log.Errorf("RpcClientFunc error: " + err.Error())
			debug.PrintStack()
			return replyMsg, err
		}

		return
	}
}

func GetService(msg entity.RpcMessage) Service {
	if msg.Method == constants.RpcInstallLang {
		return &InstallLangService{msg: msg}
	} else if msg.Method == constants.RpcGetLang {
		return &GetLangService{msg: msg}
	}
	return nil
}

func InitRpcService() error {
	for i := 0; i < viper.GetInt("rpc.workers"); i++ {
		go func() {
			for {
				// 获取当前节点
				node, err := model.GetCurrentNode()
				if err != nil {
					log.Errorf(err.Error())
					debug.PrintStack()
					continue
				}

				// 获取获取消息队列信息
				dataStr, err := database.RedisClient.BRPop(fmt.Sprintf("rpc:%s", node.Id.Hex()), 0)
				if err != nil {
					if err != redis.ErrNil {
						log.Errorf(err.Error())
						debug.PrintStack()
					}
					continue
				}

				// 反序列化消息
				var msg entity.RpcMessage
				if err := json.Unmarshal([]byte(dataStr), &msg); err != nil {
					log.Errorf(err.Error())
					debug.PrintStack()
					continue
				}

				// 获取service
				service := GetService(msg)

				// 根据Method调用本地方法
				replyMsg, err := service.ServerHandle()
				if err != nil {
					log.Errorf(err.Error())
					debug.PrintStack()
				}

				// 发送返回消息
				if err := database.RedisClient.LPush(fmt.Sprintf("rpc:%s:%s", node.Id.Hex(), replyMsg.Id), utils.ObjectToString(replyMsg)); err != nil {
					log.Errorf(err.Error())
					debug.PrintStack()
					continue
				}
			}
		}()
	}
	return nil
}
