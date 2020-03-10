package rpc

import (
	"crawlab/constants"
	"crawlab/entity"
	"crawlab/utils"
	"errors"
	"fmt"
	"github.com/apex/log"
	"os/exec"
	"runtime/debug"
)

type InstallDepService struct {
	msg entity.RpcMessage
}

func (s *InstallDepService) ServerHandle() (entity.RpcMessage, error) {
	lang := utils.GetRpcParam("lang", s.msg.Params)
	depName := utils.GetRpcParam("dep_name", s.msg.Params)
	if err := InstallDepLocal(lang, depName); err != nil {
		return entity.RpcMessage{}, err
	}
	s.msg.Result = "success"
	return s.msg, nil
}

func (s *InstallDepService) ClientHandle() (o interface{}, err error) {
	// 发起 RPC 请求，获取服务端数据
	_, err = ClientFunc(s.msg)()
	if err != nil {
		return
	}

	return
}

func InstallDepLocal(lang string, depName string) error {
	if lang == constants.Python {
		_, err := InstallPythonDepLocal(depName)
		if err != nil {
			return err
		}
	} else if lang == constants.Nodejs {
		_, err := InstallNodejsDepLocal(depName)
		if err != nil {
			return err
		}
	} else {
		return errors.New(fmt.Sprintf("%s is not implemented", lang))
	}
	return nil
}

// 安装Python本地依赖
func InstallPythonDepLocal(depName string) (string, error) {
	// 依赖镜像URL
	url := "https://pypi.tuna.tsinghua.edu.cn/simple"

	cmd := exec.Command("pip", "install", depName, "-i", url)
	outputBytes, err := cmd.Output()
	if err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return fmt.Sprintf("error: %s", err.Error()), err
	}
	return string(outputBytes), nil
}

func InstallNodejsDepLocal(depName string) (string, error) {
	// 依赖镜像URL
	url := "https://registry.npm.taobao.org"

	cmd := exec.Command("npm", "install", depName, "-g", "--registry", url)
	outputBytes, err := cmd.Output()
	if err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return fmt.Sprintf("error: %s", err.Error()), err
	}
	return string(outputBytes), nil
}

func InstallDepRemote(nodeId string, lang string, depName string) (err error) {
	params := make(map[string]string)
	params["lang"] = lang
	params["dep_name"] = depName
	s := GetService(entity.RpcMessage{
		NodeId:  nodeId,
		Method:  constants.RpcInstallDep,
		Params:  params,
		Timeout: 300,
	})
	_, err = s.ClientHandle()
	if err != nil {
		return
	}
	return
}
