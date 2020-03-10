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
	output, err := InstallLangLocal(lang)
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

// 本地安装语言
func InstallLangLocal(lang string) (o string, err error) {
	l := utils.GetLangFromLangNamePlain(lang)
	if l.Name == "" || l.InstallScript == "" {
		return "", errors.New(fmt.Sprintf("%s is not implemented", lang))
	}
	cmd := exec.Command("/bin/sh", path.Join("scripts", l.InstallScript))
	output, err := cmd.Output()
	if err != nil {
		log.Error(err.Error())
		debug.PrintStack()
		return string(output), err
	}
	return
}

// 远端安装语言
func InstallLangRemote(nodeId string, lang string) (o string, err error) {
	params := make(map[string]string)
	params["lang"] = lang
	s := GetService(entity.RpcMessage{
		NodeId:  nodeId,
		Method:  constants.RpcInstallLang,
		Params:  params,
		Timeout: 60,
	})
	_, err = s.ClientHandle()
	if err != nil {
		return
	}
	return
}
