package routes

import (
	"bytes"
	"crawlab/constants"
	"crawlab/model"
	"crawlab/services"
	"crawlab/utils"
	"encoding/csv"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"net/http"
)

type TaskListRequestData struct {
	PageNum  int    `form:"page_num"`
	PageSize int    `form:"page_size"`
	NodeId   string `form:"node_id"`
	SpiderId string `form:"spider_id"`
	Status   string `form:"status"`
}

type TaskResultsRequestData struct {
	PageNum  int `form:"page_num"`
	PageSize int `form:"page_size"`
}

func GetTaskList(c *gin.Context) {
	// 绑定数据
	data := TaskListRequestData{}
	if err := c.ShouldBindQuery(&data); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
	if data.PageNum == 0 {
		data.PageNum = 1
	}
	if data.PageSize == 0 {
		data.PageSize = 10
	}

	// 过滤条件
	query := bson.M{}
	if data.NodeId != "" {
		query["node_id"] = bson.ObjectIdHex(data.NodeId)
	}
	if data.SpiderId != "" {
		query["spider_id"] = bson.ObjectIdHex(data.SpiderId)
	}
	//新增根据任务状态获取task列表
	if data.Status != "" {
		query["status"] = data.Status
	}

	// 获取任务列表
	tasks, err := model.GetTaskList(query, (data.PageNum-1)*data.PageSize, data.PageSize, "-create_ts")
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 获取总任务数
	total, err := model.GetTaskListTotal(query)
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, ListResponse{
		Status:  "ok",
		Message: "success",
		Total:   total,
		Data:    tasks,
	})
}

func GetTask(c *gin.Context) {
	id := c.Param("id")

	result, err := model.GetTask(id)
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
	HandleSuccessData(c, result)
}

func PutTask(c *gin.Context) {
	type TaskRequestBody struct {
		SpiderId bson.ObjectId   `json:"spider_id"`
		RunType  string          `json:"run_type"`
		NodeIds  []bson.ObjectId `json:"node_ids"`
		Param    string          `json:"param"`
	}

	// 绑定数据
	var reqBody TaskRequestBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

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
				SpiderId:   reqBody.SpiderId,
				NodeId:     node.Id,
				Param:      reqBody.Param,
				UserId:     services.GetCurrentUserId(c),
				RunType:    constants.RunTypeAllNodes,
				ScheduleId: bson.ObjectIdHex(constants.ObjectIdNull),
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
			SpiderId:   reqBody.SpiderId,
			Param:      reqBody.Param,
			UserId:     services.GetCurrentUserId(c),
			RunType:    constants.RunTypeRandom,
			ScheduleId: bson.ObjectIdHex(constants.ObjectIdNull),
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
				SpiderId:   reqBody.SpiderId,
				NodeId:     nodeId,
				Param:      reqBody.Param,
				UserId:     services.GetCurrentUserId(c),
				RunType:    constants.RunTypeSelectedNodes,
				ScheduleId: bson.ObjectIdHex(constants.ObjectIdNull),
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

func DeleteTaskByStatus(c *gin.Context) {
	status := c.Query("status")

	//删除相应的日志文件
	if err := services.RemoveLogByTaskStatus(status); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	//删除该状态下的task
	if err := model.RemoveTaskByStatus(status); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	HandleSuccess(c)
}

// 删除多个任务
func DeleteSelectedTask(c *gin.Context) {
	ids := make(map[string][]string)
	if err := c.ShouldBindJSON(&ids); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
	list := ids["ids"]
	for _, id := range list {
		if err := services.RemoveLogByTaskId(id); err != nil {
			HandleError(http.StatusInternalServerError, c, err)
			return
		}
		if err := model.RemoveTask(id); err != nil {
			HandleError(http.StatusInternalServerError, c, err)
			return
		}
	}
	HandleSuccess(c)
}

// 删除单个任务
func DeleteTask(c *gin.Context) {
	id := c.Param("id")

	// 删除日志文件
	if err := services.RemoveLogByTaskId(id); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
	// 删除task
	if err := model.RemoveTask(id); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
	HandleSuccess(c)
}

func GetTaskLog(c *gin.Context) {
	id := c.Param("id")
	logStr, err := services.GetTaskLog(id)
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
	HandleSuccessData(c, logStr)
}

func GetTaskResults(c *gin.Context) {
	id := c.Param("id")

	// 绑定数据
	data := TaskResultsRequestData{}
	if err := c.ShouldBindQuery(&data); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 获取任务
	task, err := model.GetTask(id)
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 获取结果
	results, total, err := task.GetResults(data.PageNum, data.PageSize)
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, ListResponse{
		Status:  "ok",
		Message: "success",
		Data:    results,
		Total:   total,
	})
}

func DownloadTaskResultsCsv(c *gin.Context) {
	id := c.Param("id")

	// 获取任务
	task, err := model.GetTask(id)
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 获取结果
	results, _, err := task.GetResults(1, constants.Infinite)
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 字段列表
	var columns []string
	if len(results) == 0 {
		columns = []string{}
	} else {
		item := results[0].(bson.M)
		for key := range item {
			columns = append(columns, key)
		}
	}

	// 缓冲
	bytesBuffer := &bytes.Buffer{}

	// 写入UTF-8 BOM，避免使用Microsoft Excel打开乱码
	bytesBuffer.WriteString("\xEF\xBB\xBF")

	writer := csv.NewWriter(bytesBuffer)

	// 写入表头
	if err := writer.Write(columns); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 写入内容
	for _, result := range results {
		// 将result转换为[]string
		item := result.(bson.M)
		var values []string
		for _, col := range columns {
			value := utils.InterfaceToString(item[col])
			values = append(values, value)
		}

		// 写入数据
		if err := writer.Write(values); err != nil {
			HandleError(http.StatusInternalServerError, c, err)
			return
		}
	}

	// 此时才会将缓冲区数据写入
	writer.Flush()

	// 设置下载的文件名
	c.Writer.Header().Set("Content-Disposition", "attachment;filename=data.csv")

	// 设置文件类型以及输出数据
	c.Data(http.StatusOK, "text/csv", bytesBuffer.Bytes())
}

func CancelTask(c *gin.Context) {
	id := c.Param("id")

	if err := services.CancelTask(id); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}
	HandleSuccess(c)
}
