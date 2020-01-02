package services

import (
	"crawlab/constants"
	"crawlab/database"
	"crawlab/entity"
	"crawlab/model"
	"encoding/json"
	"fmt"
	"github.com/apex/log"
	uuid "github.com/satori/go.uuid"
	"runtime/debug"
)

type RpcMessage struct {
	Id     string `json:"id"`
	Method string `json:"method"`
	Params string `json:"params"`
	Result string `json:"result"`
}

func RpcServerInstallLang(msg RpcMessage) RpcMessage {
	// install dep rpc
	return msg
}

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

	data, err := RpcClientFunc(nodeId, params, 10)()
	if err != nil {
		return
	}

	output = data.(string)

	return
}

func RpcServerGetDepList(nodeId string, msg RpcMessage) RpcMessage {
	lang := GetRpcParam("lang", msg.Params)
	searchDepName := GetRpcParam("search_dep_name", msg.Params)
	if lang == constants.Python {
		depList, _ := GetPythonLocalDepList(nodeId, searchDepName)
		resultStr, _ := json.Marshal(depList)
		msg.Result = string(resultStr)
	}
	return msg
}

func RpcClientGetDepList(nodeId string, lang string, searchDepName string) (list []entity.Dependency, err error) {
	params := map[string]string{}
	params["lang"] = lang
	params["search_dep_name"] = searchDepName

	data, err := RpcClientFunc(nodeId, params, 30)()
	if err != nil {
		return
	}

	list = data.([]entity.Dependency)

	return
}

func RpcServerGetInstalledDepList(nodeId string, msg RpcMessage) RpcMessage {
	lang := GetRpcParam("lang", msg.Params)
	if lang == constants.Python {
		depList, _ := GetPythonLocalInstalledDepList(nodeId)
		resultStr, _ := json.Marshal(depList)
		msg.Result = string(resultStr)
	}
	return msg
}

func RpcClientGetInstalledDepList(nodeId string, lang string) (list []entity.Dependency, err error) {
	params := map[string]string{}
	params["lang"] = lang

	data, err := RpcClientFunc(nodeId, params, 10)()
	if err != nil {
		return
	}

	list = data.([]entity.Dependency)

	return
}

func RpcClientFunc(nodeId string, params interface{}, timeout int) func() (interface{}, error) {
	return func() (data interface{}, err error) {
		// 请求ID
		id := uuid.NewV4().String()

		// 构造RPC消息
		msg := RpcMessage{
			Id:     id,
			Method: constants.RpcGetDepList,
			Params: ObjectToString(params),
			Result: "",
		}

		// 发送RPC消息
		if err := database.RedisClient.LPush(fmt.Sprintf("rpc:%s", nodeId), ObjectToString(msg)); err != nil {
			return data, err
		}

		// 获取RPC回复消息
		dataStr, err := database.RedisClient.BRPop(fmt.Sprintf("rpc:%s", nodeId), timeout)
		if err != nil {
			return data, err
		}

		// 反序列化消息
		if err := json.Unmarshal([]byte(dataStr), &msg); err != nil {
			return data, err
		}

		// 反序列化列表
		if err := json.Unmarshal([]byte(msg.Result), &data); err != nil {
			return data, err
		}

		return data, err
	}
}

func GetRpcParam(key string, params interface{}) string {
	var paramsObj map[string]string
	if err := json.Unmarshal([]byte(params.(string)), &paramsObj); err != nil {
		return ""
	}
	return paramsObj[key]
}

func ObjectToString(params interface{}) string {
	str, _ := json.Marshal(params)
	return string(str)
}

var IsRpcStopped = false

func StopRpcService() {
	IsRpcStopped = true
}

func InitRpcService() error {
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
			dataStr, err := database.RedisClient.BRPop(fmt.Sprintf("rpc:%s", node.Id.Hex()), 300)
			if err != nil {
				log.Errorf(err.Error())
				debug.PrintStack()
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
			} else if msg.Method == constants.RpcInstallLang {
				replyMsg = RpcServerInstallLang(msg)
			} else if msg.Method == constants.RpcGetDepList {
				replyMsg = RpcServerGetDepList(node.Id.Hex(), msg)
			} else if msg.Method == constants.RpcGetInstalledDepList {
				replyMsg = RpcServerGetInstalledDepList(node.Id.Hex(), msg)
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
	return nil
}
