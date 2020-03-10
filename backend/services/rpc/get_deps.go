package rpc

import (
	"crawlab/entity"
	"crawlab/utils"
	"encoding/json"
)

type GetDepsService struct {
	msg entity.RpcMessage
}

func (s *GetDepsService) ServerHandle() (entity.RpcMessage, error) {
	langName := utils.GetRpcParam("lang", s.msg.Params)
	lang := utils.GetLangFromLangNamePlain(langName)
	l := GetLangLocal(lang)
	lang.InstallStatus = l.InstallStatus

	// 序列化
	resultStr, _ := json.Marshal(lang)
	s.msg.Result = string(resultStr)
	return s.msg, nil
}

func (s *GetDepsService) ClientHandle() (o interface{}, err error) {
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
