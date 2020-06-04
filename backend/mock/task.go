package mock

import (
	"bytes"
	"crawlab/constants"
	"crawlab/model"
	"crawlab/utils"
	"encoding/csv"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"github.com/satori/go.uuid"
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
	tasks := TaskList

	// 获取总任务数
	total := len(TaskList)

	c.JSON(http.StatusOK, ListResponse{
		Status:  "ok",
		Message: "success",
		Total:   total,
		Data:    tasks,
	})
}

func GetTask(c *gin.Context) {
	id := c.Param("id")
	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
		return
	}
	var result model.Task
	for _, task := range TaskList {
		if task.Id == id {
			result = task
		}
	}
	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    result,
	})
}

func PutTask(c *gin.Context) {
	// 生成任务ID,generate task ID
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

	// 将任务存入数据库,put the task into database
	fmt.Println("put the task into database")

	// 加入任务队列, put the task into task queue
	fmt.Println("put the task into task queue")

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

func DeleteTask(c *gin.Context) {
	id := c.Param("id")
	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
		return
	}
	for _, task := range TaskList {
		if task.Id == id {
			fmt.Println("delete the task")
		}
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

func GetTaskResults(c *gin.Context) {
	id := c.Param("id")
	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
		return
	}
	// 绑定数据
	data := TaskResultsRequestData{}
	if err := c.ShouldBindQuery(&data); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	// 获取任务
	var task model.Task
	for _, ta := range TaskList {
		if ta.Id == id {
			task = ta
		}
	}

	fmt.Println(task)
	// 获取结果
	var results interface{}
	total := len(TaskList)

	c.JSON(http.StatusOK, ListResponse{
		Status:  "ok",
		Message: "success",
		Data:    results,
		Total:   total,
	})
}

func DownloadTaskResultsCsv(c *gin.Context) {
	id := c.Param("id")
	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
		return
	}
	// 获取任务
	var task model.Task
	for _, ta := range TaskList {
		if ta.Id == id {
			task = ta
		}
	}
	fmt.Println(task)

	// 获取结果
	var results []interface {
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
