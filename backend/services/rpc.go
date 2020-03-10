package services

import (
	"crawlab/constants"
	"crawlab/database"
	"crawlab/entity"
	"crawlab/utils"
	"encoding/json"
	"fmt"
	"github.com/apex/log"
	uuid "github.com/satori/go.uuid"
	"runtime/debug"
)

type RpcMessage struct {
	Id      string            `json:"id"`
	Method  string            `json:"method"`
	Blocked bool              `json:"blocked"`
	Params  map[string]string `json:"params"`
	Result  string            `json:"result"`
}

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
