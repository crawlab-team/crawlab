package services

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

type RpcMessage struct {
	Id      string            `json:"id"`
	Method  string            `json:"method"`
	Blocked bool              `json:"blocked"`
	Params  map[string]string `json:"params"`
	Result  string            `json:"result"`
}

// ========安装语言========

func RpcServerInstallLang(msg RpcMessage) RpcMessage {
	lang := GetRpcParam("lang", msg.Params)
	if lang == constants.Nodejs {
		output, _ := InstallNodejsLocalLang()
		msg.Result = output
	}
	return msg
}

func RpcClientInstallLang(nodeId string, lang string) (output string, err error) {
	params := map[string]string{}
	params["lang"] = lang

	// 发起 RPC 请求，获取服务端数据
	go func() {
		_, err := RpcClientFunc(nodeId, constants.RpcInstallLang, params, 600)()
		if err != nil {
			return
		}
	}()

	return
}

// ========./安装语言========

// ========获取语言========

func RpcServerGetLang(msg RpcMessage) RpcMessage {
	langName := GetRpcParam("lang", msg.Params)
	lang := GetLangFromLangNamePlain(langName)
	l := GetLangLocal(lang)
	lang.InstallStatus = l.InstallStatus

	// 序列化
	resultStr, _ := json.Marshal(lang)
	msg.Result = string(resultStr)
	return msg
}

func RpcClientGetLang(nodeId string, langName string) (lang entity.Lang, err error) {
	params := map[string]string{}
	params["lang"] = langName

	data, err := RpcClientFunc(nodeId, constants.RpcGetLang, params, 30)()
	if err != nil {
		return
	}

	// 反序列化结果
	if err := json.Unmarshal([]byte(data), &lang); err != nil {
		return lang, err
	}

	return
}

// ========./获取语言========

// ========安装依赖========

func RpcServerInstallDep(msg RpcMessage) RpcMessage {
	lang := GetRpcParam("lang", msg.Params)
	depName := GetRpcParam("dep_name", msg.Params)
	if lang == constants.Python {
		output, _ := InstallPythonLocalDep(depName)
		msg.Result = output
	}
	return msg
}

func RpcClientInstallDep(nodeId string, lang string, depName string) (output string, err error) {
	params := map[string]string{}
	params["lang"] = lang
	params["dep_name"] = depName

	data, err := RpcClientFunc(nodeId, constants.RpcInstallDep, params, 60)()
	if err != nil {
		return
	}

	output = data

	return
}

// ========./安装依赖========

// ========卸载依赖========

func RpcServerUninstallDep(msg RpcMessage) RpcMessage {
	lang := GetRpcParam("lang", msg.Params)
	depName := GetRpcParam("dep_name", msg.Params)
	if lang == constants.Python {
		output, _ := UninstallPythonLocalDep(depName)
		msg.Result = output
	}
	return msg
}

func RpcClientUninstallDep(nodeId string, lang string, depName string) (output string, err error) {
	params := map[string]string{}
	params["lang"] = lang
	params["dep_name"] = depName

	data, err := RpcClientFunc(nodeId, constants.RpcUninstallDep, params, 60)()
	if err != nil {
		return
	}

	output = data

	return
}

// ========./卸载依赖========

// ========获取已安装依赖列表========

func RpcServerGetInstalledDepList(nodeId string, msg RpcMessage) RpcMessage {
	lang := GetRpcParam("lang", msg.Params)
	if lang == constants.Python {
		depList, _ := GetPythonLocalInstalledDepList(nodeId)
		resultStr, _ := json.Marshal(depList)
		msg.Result = string(resultStr)
	} else if lang == constants.Nodejs {
		depList, _ := GetNodejsLocalInstalledDepList(nodeId)
		resultStr, _ := json.Marshal(depList)
		msg.Result = string(resultStr)
	}
	return msg
}

func RpcClientGetInstalledDepList(nodeId string, lang string) (list []entity.Dependency, err error) {
	params := map[string]string{}
	params["lang"] = lang

	data, err := RpcClientFunc(nodeId, constants.RpcGetInstalledDepList, params, 30)()
	if err != nil {
		return
	}

	// 反序列化结果
	if err := json.Unmarshal([]byte(data), &list); err != nil {
		return list, err
	}

	return
}

// ========./获取已安装依赖列表========

// RPC 客户端函数
func RpcClientFunc(nodeId string, method string, params map[string]string, timeout int) func() (string, error) {
	return func() (result string, err error) {
		// 请求ID
		id := uuid.NewV4().String()

		// 构造RPC消息
		msg := RpcMessage{
			Id:     id,
			Method: method,
			Params: params,
			Result: "",
		}

		// 发送RPC消息
		msgStr := ObjectToString(msg)
		if err := database.RedisClient.LPush(fmt.Sprintf("rpc:%s", nodeId), msgStr); err != nil {
			log.Errorf("RpcClientFunc error: " + err.Error())
			debug.PrintStack()
			return result, err
		}

		// 获取RPC回复消息
		dataStr, err := database.RedisClient.BRPop(fmt.Sprintf("rpc:%s", nodeId), timeout)
		if err != nil {
			log.Errorf("RpcClientFunc error: " + err.Error())
			debug.PrintStack()
			return result, err
		}

		// 反序列化消息
		if err := json.Unmarshal([]byte(dataStr), &msg); err != nil {
			log.Errorf("RpcClientFunc error: " + err.Error())
			debug.PrintStack()
			return result, err
		}

		return msg.Result, nil
	}
}

// 获取 RPC 参数
func GetRpcParam(key string, params map[string]string) string {
	return params[key]
}

// Object 转化为 String
func ObjectToString(params interface{}) string {
	bytes, _ := json.Marshal(params)
	return utils.BytesToString(bytes)
}

var IsRpcStopped = false

func StopRpcService() {
	IsRpcStopped = true
}

// 初始化 RPC 服务
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
				var msg RpcMessage
				if err := json.Unmarshal([]byte(dataStr), &msg); err != nil {
					log.Errorf(err.Error())
					debug.PrintStack()
					continue
				}

				// 根据Method调用本地方法
				var replyMsg RpcMessage
				if msg.Method == constants.RpcInstallDep {
					replyMsg = RpcServerInstallDep(msg)
				} else if msg.Method == constants.RpcUninstallDep {
					replyMsg = RpcServerUninstallDep(msg)
				} else if msg.Method == constants.RpcInstallLang {
					replyMsg = RpcServerInstallLang(msg)
				} else if msg.Method == constants.RpcGetInstalledDepList {
					replyMsg = RpcServerGetInstalledDepList(node.Id.Hex(), msg)
				} else if msg.Method == constants.RpcGetLang {
					replyMsg = RpcServerGetLang(msg)
				} else {
					continue
				}

				// 发送返回消息
				if err := database.RedisClient.LPush(fmt.Sprintf("rpc:%s", node.Id.Hex()), ObjectToString(replyMsg)); err != nil {
					log.Errorf(err.Error())
					debug.PrintStack()
					continue
				}

				// 如果停止RPC服务，则返回
				if IsRpcStopped {
					return
				}
			}
		}()
	}
	return nil
}
