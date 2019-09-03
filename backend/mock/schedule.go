package mock

import (
	"crawlab/constants"
	"crawlab/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"net/http"
	"time"
)

var scheduleList = []model.Schedule{
	{
		Id:       bson.ObjectId("5d429e6c19f7abede924fee2"),
		Name:     "test schedule",
		SpiderId: "123",
		NodeId:   bson.ObjectId("5d429e6c19f7abede924fee2"),
		Cron:     "***1*",
		EntryId:  10,
		// 前端展示
		SpiderName: "test scedule",
		NodeName:   "测试节点",

		CreateTs: time.Now(),
		UpdateTs: time.Now(),
	},
	{
		Id:       bson.ObjectId("xx429e6c19f7abede924fee2"),
		Name:     "test schedule2",
		SpiderId: "234",
		NodeId:   bson.ObjectId("5d429e6c19f7abede924fee2"),
		Cron:     "***1*",
		EntryId:  10,
		// 前端展示
		SpiderName: "test scedule2",
		NodeName:   "测试节点",

		CreateTs: time.Now(),
		UpdateTs: time.Now(),
	},
}

func GetScheduleList(c *gin.Context) {
	results := scheduleList

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    results,
	})
}

func GetSchedule(c *gin.Context) {
	id := c.Param("id")

	var result model.Schedule
	for _, sch := range scheduleList {
		if sch.Id == bson.ObjectId(id) {
			result = sch
		}
	}
	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    result,
	})
}

func PostSchedule(c *gin.Context) {
	id := c.Param("id")
	var oldItem model.Schedule
	for _, sch := range scheduleList {
		if sch.Id == bson.ObjectId(id) {
			oldItem = sch
		}

	}

	var newItem model.Schedule
	if err := c.ShouldBindJSON(&newItem); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}
	newItem.Id = oldItem.Id

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

func PutSchedule(c *gin.Context) {
	var item model.Schedule

	// 绑定数据模型
	if err := c.ShouldBindJSON(&item); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}

	// 如果node_id为空，则置为空ObjectId
	if item.NodeId == "" {
		item.NodeId = bson.ObjectIdHex(constants.ObjectIdNull)
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

func DeleteSchedule(c *gin.Context) {
	id := bson.ObjectIdHex("5d429e6c19f7abede924fee2")
	for _, sch := range scheduleList {
		if sch.Id == id {
			fmt.Println("delete a schedule")
		}
	}
	fmt.Println(id)
	fmt.Println("update schedule")
	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}
