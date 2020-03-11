package rpc

import (
	"crawlab/constants"
	"crawlab/entity"
	"crawlab/utils"
	"encoding/json"
	"os/exec"
	"regexp"
	"runtime/debug"
	"strings"
)

type GetInstalledDepsService struct {
	msg entity.RpcMessage
}

func (s *GetInstalledDepsService) ServerHandle() (entity.RpcMessage, error) {
	lang := utils.GetRpcParam("lang", s.msg.Params)
	deps, err := GetInstalledDepsLocal(lang)
	if err != nil {
		s.msg.Error = err.Error()
		return s.msg, err
	}
	resultStr, _ := json.Marshal(deps)
	s.msg.Result = string(resultStr)
	return s.msg, nil
}

func (s *GetInstalledDepsService) ClientHandle() (o interface{}, err error) {
	// 发起 RPC 请求，获取服务端数据
	s.msg, err = ClientFunc(s.msg)()
	if err != nil {
		return o, err
	}

	// 反序列化
	var output []entity.Dependency
	if err := json.Unmarshal([]byte(s.msg.Result), &output); err != nil {
		return o, err
	}
	o = output

	return
}

// 获取本地已安装依赖列表
func GetInstalledDepsLocal(lang string) (deps []entity.Dependency, err error) {
	if lang == constants.Python {
		deps, err = GetPythonInstalledDepListLocal()
	} else if lang == constants.Nodejs {
		deps, err = GetNodejsInstalledDepListLocal()
	}
	return deps, err
}

// 获取Python本地已安装依赖列表
func GetPythonInstalledDepListLocal() ([]entity.Dependency, error) {
	var list []entity.Dependency

	cmd := exec.Command("pip", "freeze")
	outputBytes, err := cmd.Output()
	if err != nil {
		debug.PrintStack()
		return list, err
	}

	for _, line := range strings.Split(string(outputBytes), "\n") {
		arr := strings.Split(line, "==")
		if len(arr) < 2 {
			continue
		}
		dep := entity.Dependency{
			Name:      strings.ToLower(arr[0]),
			Version:   arr[1],
			Installed: true,
		}
		list = append(list, dep)
	}

	return list, nil
}

// 获取Node.js本地已安装依赖列表
func GetNodejsInstalledDepListLocal() ([]entity.Dependency, error) {
	var list []entity.Dependency

	cmd := exec.Command("npm", "ls", "-g", "--depth", "0")
	outputBytes, _ := cmd.Output()

	regex := regexp.MustCompile("\\s(.*)@(.*)")
	for _, line := range strings.Split(string(outputBytes), "\n") {
		arr := regex.FindStringSubmatch(line)
		if len(arr) < 3 {
			continue
		}
		dep := entity.Dependency{
			Name:      strings.ToLower(arr[1]),
			Version:   arr[2],
			Installed: true,
		}
		list = append(list, dep)
	}

	return list, nil
}

func GetInstalledDepsRemote(nodeId string, lang string) (deps []entity.Dependency, err error) {
	params := make(map[string]string)
	params["lang"] = lang
	s := GetService(entity.RpcMessage{
		NodeId:  nodeId,
		Method:  constants.RpcGetInstalledDepList,
		Params:  params,
		Timeout: 60,
	})
	o, err := s.ClientHandle()
	if err != nil {
		return
	}
	deps = o.([]entity.Dependency)
	return
}
