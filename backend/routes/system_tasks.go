package routes

import (
	"crawlab/constants"
	"crawlab/model"
	"crawlab/services"
	"crawlab/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"net/http"
)

func GetSystemScripts(c *gin.Context) {
	HandleSuccessData(c, utils.GetSystemScripts())
}

func PutSystemTask(c *gin.Context) {
	type TaskRequestBody struct {
		RunType string          `json:"run_type"`
		NodeIds []bson.ObjectId `json:"node_ids"`
		Script  string          `json:"script"`
	}

	// 绑定数据
	var reqBody TaskRequestBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	// 校验脚本参数不为空
	if reqBody.Script == "" {
		HandleErrorF(http.StatusBadRequest, c, "script cannot be empty")
		return
	}

	// 校验脚本参数是否存在
	var allScripts = utils.GetSystemScripts()
	if !utils.StringArrayContains(allScripts, reqBody.Script) {
		HandleErrorF(http.StatusBadRequest, c, "script does not exist")
		return
	}

	// 获取执行命令
	cmd := fmt.Sprintf("sh %s", utils.GetSystemScriptPath(reqBody.Script))

	// 任务ID
	var taskIds []string

	if reqBody.RunType == constants.RunTypeAllNodes {
		// 所有节点
		nodes, err := model.GetNodeList(nil)
		if err != nil {
			HandleError(http.StatusInternalServerError, c, err)
			return
		}
		for _, node := range nodes {
			t := model.Task{
				SpiderId:   bson.ObjectIdHex(constants.ObjectIdNull),
				NodeId:     node.Id,
				UserId:     services.GetCurrentUserId(c),
				RunType:    constants.RunTypeAllNodes,
				ScheduleId: bson.ObjectIdHex(constants.ObjectIdNull),
				Type:       constants.TaskTypeSystem,
				Cmd:        cmd,
			}
			id, err := services.AddTask(t)
			if err != nil {
				HandleError(http.StatusInternalServerError, c, err)
				return
			}
			taskIds = append(taskIds, id)
		}
	} else if reqBody.RunType == constants.RunTypeRandom {
		// 随机
		t := model.Task{
			SpiderId:   bson.ObjectIdHex(constants.ObjectIdNull),
			UserId:     services.GetCurrentUserId(c),
			RunType:    constants.RunTypeRandom,
			ScheduleId: bson.ObjectIdHex(constants.ObjectIdNull),
			Type:       constants.TaskTypeSystem,
			Cmd:        cmd,
		}
		id, err := services.AddTask(t)
		if err != nil {
			HandleError(http.StatusInternalServerError, c, err)
			return
		}
		taskIds = append(taskIds, id)
	} else if reqBody.RunType == constants.RunTypeSelectedNodes {
		// 指定节点
		for _, nodeId := range reqBody.NodeIds {
			t := model.Task{
				SpiderId:   bson.ObjectIdHex(constants.ObjectIdNull),
				NodeId:     nodeId,
				UserId:     services.GetCurrentUserId(c),
				RunType:    constants.RunTypeSelectedNodes,
				ScheduleId: bson.ObjectIdHex(constants.ObjectIdNull),
				Type:       constants.TaskTypeSystem,
				Cmd:        cmd,
			}
			id, err := services.AddTask(t)
			if err != nil {
				HandleError(http.StatusInternalServerError, c, err)
				return
			}
			taskIds = append(taskIds, id)
		}
	} else {
		HandleErrorF(http.StatusInternalServerError, c, "invalid run_type")
		return
	}

	HandleSuccessData(c, taskIds)
}
