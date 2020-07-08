package rpc

import (
	"crawlab/constants"
	"crawlab/database"
	"crawlab/entity"
	"crawlab/model"
	"crawlab/services/local_node"
	"crawlab/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/apex/log"
	"github.com/cenkalti/backoff/v4"
	"github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
	"runtime/debug"
)

// RPC服务基础类
type Service interface {
	ServerHandle() (entity.RpcMessage, error)
	ClientHandle() (interface{}, error)
}

// 客户端处理消息函数
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

		// 如果返回消息有错误，返回错误
		if replyMsg.Error != "" {
			return replyMsg, errors.New(replyMsg.Error)
		}

		return
	}
}

// 获取RPC服务
func GetService(msg entity.RpcMessage) Service {
	switch msg.Method {
	case constants.RpcInstallLang:
		return &InstallLangService{msg: msg}
	case constants.RpcInstallDep:
		return &InstallDepService{msg: msg}
	case constants.RpcUninstallDep:
		return &UninstallDepService{msg: msg}
	case constants.RpcGetLang:
		return &GetLangService{msg: msg}
	case constants.RpcGetInstalledDepList:
		return &GetInstalledDepsService{msg: msg}
	case constants.RpcCancelTask:
		return &CancelTaskService{msg: msg}
	case constants.RpcGetSystemInfoService:
		return &GetSystemInfoService{msg: msg}
	}
	return nil
}

// 处理RPC消息
func handleMsg(msgStr string, node *model.Node) {
	// 反序列化消息
	var msg entity.RpcMessage
	if err := json.Unmarshal([]byte(msgStr), &msg); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return
	}

	// 获取service
	service := GetService(msg)

	// 根据Method调用本地方法
	replyMsg, err := service.ServerHandle()
	if err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return
	}

	// 发送返回消息
	if err := database.RedisClient.LPush(fmt.Sprintf("rpc:%s:%s", node.Id.Hex(), replyMsg.Id), utils.ObjectToString(replyMsg)); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
	}
}

// 初始化服务端RPC服务
func InitRpcService() error {
	go func() {
		node := local_node.CurrentNode()
		for {
			// 获取当前节点
			//node, err := model.GetCurrentNode()
			//if err != nil {
			//	log.Errorf(err.Error())
			//	debug.PrintStack()
			//	continue
			//}
			b := backoff.NewExponentialBackOff()
			bp := backoff.WithMaxRetries(b, 10)
			var msgStr string
			var err error
			err = backoff.Retry(func() error {
				msgStr, err = database.RedisClient.BRPop(fmt.Sprintf("rpc:%s", node.Id.Hex()), 0)

				if err != nil && err != redis.ErrNil {
					log.WithError(err).Warnf("waiting for redis pool active connection. will after %f seconds try  again.", b.NextBackOff().Seconds())
					return err
				}
				return err
			}, bp)
			if err != nil {
				continue
			}
			// 处理消息
			go handleMsg(msgStr, node)
		}
	}()
	return nil
}
