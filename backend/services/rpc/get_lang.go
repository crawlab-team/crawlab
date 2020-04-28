package rpc

import (
	"crawlab/constants"
	"crawlab/entity"
	"crawlab/utils"
	"encoding/json"
)

type GetLangService struct {
	msg entity.RpcMessage
}

func (s *GetLangService) ServerHandle() (entity.RpcMessage, error) {
	langName := utils.GetRpcParam("lang", s.msg.Params)
	lang := utils.GetLangFromLangNamePlain(langName)
	l := GetLangLocal(lang)
	lang.InstallStatus = l.InstallStatus

	// 序列化
	resultStr, _ := json.Marshal(lang)
	s.msg.Result = string(resultStr)
	return s.msg, nil
}

func (s *GetLangService) ClientHandle() (o interface{}, err error) {
	// 发起 RPC 请求，获取服务端数据
	s.msg, err = ClientFunc(s.msg)()
	if err != nil {
		return o, err
	}

	var output entity.Lang
	if err := json.Unmarshal([]byte(s.msg.Result), &output); err != nil {
		return o, err
	}
	o = output

	return
}

func GetLangLocal(lang entity.Lang) entity.Lang {
	// 检查是否存在执行路径
	for _, p := range lang.ExecutablePaths {
		if utils.Exists(p) {
			lang.InstallStatus = constants.InstallStatusInstalled
			return lang
		}
	}

	// 检查是否正在安装
	if utils.Exists(lang.LockPath) {
		lang.InstallStatus = constants.InstallStatusInstalling
		return lang
	}

	// 检查其他语言是否在安装
	if utils.Exists("/tmp/install.lock") {
		lang.InstallStatus = constants.InstallStatusInstallingOther
		return lang
	}

	lang.InstallStatus = constants.InstallStatusNotInstalled
	return lang
}

func GetLangRemote(nodeId string, lang entity.Lang) (l entity.Lang, err error) {
	params := make(map[string]string)
	params["lang"] = lang.ExecutableName
	s := GetService(entity.RpcMessage{
		NodeId:  nodeId,
		Method:  constants.RpcGetLang,
		Params:  params,
		Timeout: 60,
	})
	o, err := s.ClientHandle()
	if err != nil {
		return
	}
	l = o.(entity.Lang)
	return
}
