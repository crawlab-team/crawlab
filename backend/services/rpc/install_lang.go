package rpc

import (
	"crawlab/constants"
	"crawlab/entity"
	"crawlab/utils"
	"errors"
	"fmt"
	"github.com/apex/log"
	"os/exec"
	"path"
	"runtime/debug"
)

type InstallLangService struct {
	msg entity.RpcMessage
}

func (s *InstallLangService) ServerHandle() (entity.RpcMessage, error) {
	lang := utils.GetRpcParam("lang", s.msg.Params)
	output, err := InstallLocalLang(lang)
	s.msg.Result = output
	if err != nil {
		s.msg.Error = err.Error()
		return s.msg, err
	}
	return s.msg, nil
}

func (s *InstallLangService) ClientHandle() (o interface{}, err error) {
	// 发起 RPC 请求，获取服务端数据
	go func() {
		_, err := ClientFunc(s.msg)()
		if err != nil {
			return
		}
	}()

	return
}

// 本地安装Node.js
func InstallNodejsLocalLang() (string, error) {
	cmd := exec.Command("/bin/sh", path.Join("scripts", "install-nodejs.sh"))
	output, err := cmd.Output()
	if err != nil {
		log.Error(err.Error())
		debug.PrintStack()
		return string(output), err
	}

	// TODO: check if Node.js is installed successfully

	return string(output), nil
}

// 本地安装Java
func InstallJavaLocalLang() (string, error) {
	cmd := exec.Command("/bin/sh", path.Join("scripts", "install-java.sh"))
	output, err := cmd.Output()
	if err != nil {
		log.Error(err.Error())
		debug.PrintStack()
		return string(output), err
	}

	// TODO: check if Java is installed successfully

	return string(output), nil
}

// 本地安装语言
func InstallLocalLang(lang string) (o string, err error) {
	if lang == constants.Nodejs {
		o, err = InstallNodejsLocalLang()
	} else if lang == constants.Java {
		o, err = InstallNodejsLocalLang()
	} else {
		return "", errors.New(fmt.Sprintf("%s is not implemented", lang))
	}
	return
}

// 远端安装语言
func InstallRemoteLang(nodeId string, lang string) (o string, err error) {
	params := make(map[string]string)
	params["lang"] = lang
	s := GetService(entity.RpcMessage{
		NodeId:  nodeId,
		Method:  constants.RpcInstallLang,
		Params:  params,
		Timeout: 60,
	})
	output, err := s.ClientHandle()
	o = output.(string)
	if err != nil {
		return
	}
	return
}
