package rpc

import (
	"crawlab/constants"
	"crawlab/entity"
	"crawlab/model"
	"encoding/json"
)

type GetSystemInfoService struct {
	msg entity.RpcMessage
}

func (s *GetSystemInfoService) ServerHandle() (entity.RpcMessage, error) {
	sysInfo, err := GetSystemInfoServiceLocal()
	if err != nil {
		s.msg.Error = err.Error()
		return s.msg, err
	}

	// 序列化
	resultStr, _ := json.Marshal(sysInfo)
	s.msg.Result = string(resultStr)
	return s.msg, nil
}

func (s *GetSystemInfoService) ClientHandle() (o interface{}, err error) {
	// 发起 RPC 请求，获取服务端数据
	s.msg, err = ClientFunc(s.msg)()
	if err != nil {
		return o, err
	}

	var output entity.SystemInfo
	if err := json.Unmarshal([]byte(s.msg.Result), &output); err != nil {
		return o, err
	}
	o = output

	return
}

func GetSystemInfoServiceLocal() (sysInfo entity.SystemInfo, err error) {
	// 获取环境信息
	sysInfo, err = model.GetLocalSystemInfo()
	if err != nil {
		return sysInfo, err
	}
	return sysInfo, nil
}

func GetSystemInfoServiceRemote(nodeId string) (sysInfo entity.SystemInfo, err error) {
	params := make(map[string]string)
	params["node_id"] = nodeId
	s := GetService(entity.RpcMessage{
		NodeId:  nodeId,
		Method:  constants.RpcGetSystemInfoService,
		Params:  params,
		Timeout: 60,
	})
	o, err := s.ClientHandle()
	if err != nil {
		return
	}
	sysInfo = o.(entity.SystemInfo)
	return
}
