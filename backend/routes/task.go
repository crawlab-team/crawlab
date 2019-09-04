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
	uuid "github.com/satori/go.uuid"
	"net/http"
)

type TaskListRequestData struct {
	PageNum  int    `form:"page_num"`
	PageSize int    `form:"page_size"`
	NodeId   string `form:"node_id"`
	SpiderId string `form:"spider_id"`
}

type TaskResultsRequestData struct {
	PageNum  int `form:"page_num"`
	PageSize int `form:"page_size"`
}

func GetTaskList(c *gin.Context) {
	// 绑定数据
	data := TaskListRequestData{}
	if err := c.ShouldBindQuery(&data); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}
	if data.PageNum == 0 {
		data.PageNum = 1
	}
	if data.PageSize == 0 {
		data.PageNum = 10
	}

	// 过滤条件
	query := bson.M{}
	if data.NodeId != "" {
		query["node_id"] = bson.ObjectIdHex(data.NodeId)
	}
	if data.SpiderId != "" {
		query["spider_id"] = bson.ObjectIdHex(data.SpiderId)
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
	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    result,
	})
}

func PutTask(c *gin.Context) {
	// 生成任务ID
	id := uuid.NewV4()

	// 绑定数据
	var t model.Task
	if err := c.ShouldBindJSON(&t); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}
	t.Id = id.String()
	t.Status = constants.StatusPending

	// 如果没有传入node_id，则置为null
	if t.NodeId.Hex() == "" {
		t.NodeId = bson.ObjectIdHex(constants.ObjectIdNull)
	}

	// 将任务存入数据库
	if err := model.AddTask(t); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	// 加入任务队列
	if err := services.AssignTask(t); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

func DeleteTask(c *gin.Context) {
	id := c.Param("id")

	if err := model.RemoveTask(id); err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

func GetTaskLog(c *gin.Context) {
	id := c.Param("id")

	logStr, err := services.GetTaskLog(id)
	if err != nil {
		HandleError(http.StatusInternalServerError, c, err)
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    logStr,
	})
}

func GetTaskResults(c *gin.Context) {
	id := c.Param("id")

	// 绑定数据
	data := TaskResultsRequestData{}
	if err := c.ShouldBindQuery(&data); err != nil {
		HandleError(http.StatusBadRequest, c, err)
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

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}
