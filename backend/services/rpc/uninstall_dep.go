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

type UninstallDepService struct {
	msg entity.RpcMessage
}

func (s *UninstallDepService) ServerHandle() (entity.RpcMessage, error) {
	lang := utils.GetRpcParam("lang", s.msg.Params)
	depName := utils.GetRpcParam("dep_name", s.msg.Params)
	if err := UninstallDepLocal(lang, depName); err != nil {
		s.msg.Error = err.Error()
		return s.msg, err
	}
	s.msg.Result = "success"
	return s.msg, nil
}

func (s *UninstallDepService) ClientHandle() (o interface{}, err error) {
	// 发起 RPC 请求，获取服务端数据
	_, err = ClientFunc(s.msg)()
	if err != nil {
		return
	}

	return
}

func UninstallDepLocal(lang string, depName string) error {
	if lang == constants.Python {
		output, err := UninstallPythonDepLocal(depName)
		if err != nil {
			log.Debugf(output)
			return err
		}
	} else if lang == constants.Nodejs {
		output, err := UninstallNodejsDepLocal(depName)
		if err != nil {
			log.Debugf(output)
			return err
		}
	} else {
		return errors.New(fmt.Sprintf("%s is not implemented", lang))
	}
	return nil
}

func UninstallPythonDepLocal(depName string) (string, error) {
	cmd := exec.Command("pip", "uninstall", "-y", depName)
	outputBytes, err := cmd.Output()
	if err != nil {
		log.Errorf(string(outputBytes))
		log.Errorf(err.Error())
		debug.PrintStack()
		return fmt.Sprintf("error: %s", err.Error()), err
	}
	return string(outputBytes), nil
}

func UninstallNodejsDepLocal(depName string) (string, error) {
	cmd := exec.Command("npm", "uninstall", depName, "-g")
	outputBytes, err := cmd.Output()
	if err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return fmt.Sprintf("error: %s", err.Error()), err
	}
	return string(outputBytes), nil
}

func UninstallDepRemote(nodeId string, lang string, depName string) (err error) {
	params := make(map[string]string)
	params["lang"] = lang
	params["dep_name"] = depName
	s := GetService(entity.RpcMessage{
		NodeId:  nodeId,
		Method:  constants.RpcUninstallDep,
		Params:  params,
		Timeout: 300,
	})
	_, err = s.ClientHandle()
	if err != nil {
		return
	}
	return
}
