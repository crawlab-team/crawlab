package rpc

import (
	"crawlab/constants"
	"crawlab/entity"
	"crawlab/model"
	"crawlab/utils"
	"errors"
	"fmt"
	"github.com/globalsign/mgo/bson"
)

type CancelTaskService struct {
	msg entity.RpcMessage
}

func (s *CancelTaskService) ServerHandle() (entity.RpcMessage, error) {
	taskId := utils.GetRpcParam("task_id", s.msg.Params)
	nodeId := utils.GetRpcParam("node_id", s.msg.Params)
	if err := CancelTaskLocal(taskId, nodeId); err != nil {
		s.msg.Error = err.Error()
		return s.msg, err
	}
	s.msg.Result = "success"
	return s.msg, nil
}

func (s *CancelTaskService) ClientHandle() (o interface{}, err error) {
	// 发起 RPC 请求，获取服务端数据
	_, err = ClientFunc(s.msg)()
	if err != nil {
		return
	}

	return
}

func CancelTaskLocal(taskId string, nodeId string) error {
	if !utils.TaskExecChanMap.HasChanKey(taskId) {
		_ = model.UpdateTaskToAbnormal(bson.ObjectIdHex(nodeId))
		return errors.New(fmt.Sprintf("task id (%s) does not exist", taskId))
	}
	ch := utils.TaskExecChanMap.ChanBlocked(taskId)
	ch <- constants.TaskCancel
	return nil
}

func CancelTaskRemote(taskId string, nodeId string) (err error) {
	params := make(map[string]string)
	params["task_id"] = taskId
	params["node_id"] = nodeId
	s := GetService(entity.RpcMessage{
		NodeId:  nodeId,
		Method:  constants.RpcCancelTask,
		Params:  params,
		Timeout: 60,
	})
	_, err = s.ClientHandle()
	if err != nil {
		return
	}
	return
}
