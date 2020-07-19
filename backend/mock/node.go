package mock

import (
	"crawlab/entity"
	"crawlab/model"
	"crawlab/services"
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"net/http"
	"time"
)

var NodeList = []model.Node{
	{
		Id:           bson.ObjectId("5d429e6c19f7abede924fee2"),
		Ip:           "10.32.35.15",
		Name:         "test1",
		Status:       "online",
		Port:         "8081",
		Mac:          "ac:12:df:12:fd",
		Description:  "For test1",
		IsMaster:     true,
		UpdateTs:     time.Now(),
		CreateTs:     time.Now(),
		UpdateTsUnix: time.Now().Unix(),
	},
	{
		Id:           bson.ObjectId("5d429e6c19f7abede924fe22"),
		Ip:           "10.32.35.12",
		Name:         "test2",
		Status:       "online",
		Port:         "8082",
		Mac:          "ac:12:df:12:vh",
		Description:  "For test2",
		IsMaster:     true,
		UpdateTs:     time.Now(),
		CreateTs:     time.Now(),
		UpdateTsUnix: time.Now().Unix(),
	},
}

var TaskList = []model.Task{
	{
		Id:              "1234",
		SpiderId:        bson.ObjectId("5d429e6c19f7abede924fee2"),
		StartTs:         time.Now(),
		FinishTs:        time.Now(),
		Status:          "进行中",
		NodeId:          bson.ObjectId("5d429e6c19f7abede924fee2"),
		LogPath:         "./log",
		Cmd:             "scrapy crawl test",
		Error:           "",
		ResultCount:     0,
		WaitDuration:    10.0,
		RuntimeDuration: 10,
		TotalDuration:   20,
		SpiderName:      "test",
		NodeName:        "test",
		CreateTs:        time.Now(),
		UpdateTs:        time.Now(),
	},
	{
		Id:              "5678",
		SpiderId:        bson.ObjectId("5d429e6c19f7abede924fee2"),
		StartTs:         time.Now(),
		FinishTs:        time.Now(),
		Status:          "进行中",
		NodeId:          bson.ObjectId("5d429e6c19f7abede924fee2"),
		LogPath:         "./log",
		Cmd:             "scrapy crawl test2",
		Error:           "",
		ResultCount:     0,
		WaitDuration:    10.0,
		RuntimeDuration: 10,
		TotalDuration:   20,
		SpiderName:      "test",
		NodeName:        "test",
		CreateTs:        time.Now(),
		UpdateTs:        time.Now(),
	},
}

var dataList = []services.Data{
	{
		Mac:          "ac:12:fc:fd:ds:dd",
		Ip:           "192.10.2.1",
		Master:       true,
		UpdateTs:     time.Now(),
		UpdateTsUnix: time.Now().Unix(),
	},
	{
		Mac:          "22:12:fc:fd:ds:dd",
		Ip:           "182.10.2.2",
		Master:       true,
		UpdateTs:     time.Now(),
		UpdateTsUnix: time.Now().Unix(),
	},
}

var executeble = []entity.Executable{
	{
		Path:        "/test",
		FileName:    "test.py",
		DisplayName: "test.py",
	},
}
var systemInfo = entity.SystemInfo{ARCH: "x86",
	OS:          "linux",
	Hostname:    "test",
	NumCpu:      4,
	Executables: executeble,
}

func GetNodeList(c *gin.Context) {
	nodes := NodeList

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    nodes,
	})
}

func GetNode(c *gin.Context) {
	var result model.Node
	id := c.Param("id")
	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
		return
	}
	for _, node := range NodeList {
		if node.Id == bson.ObjectId(id) {
			result = node
		}
	}

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    result,
	})
}

func Ping(c *gin.Context) {
	data := dataList[0]

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    data,
	})
}

func PostNode(c *gin.Context) {
	id := c.Param("id")
	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
		return
	}
	var oldItem model.Node
	for _, node := range NodeList {
		if node.Id == bson.ObjectId(id) {
			oldItem = node
		}

	}
	log.Info(id)
	var newItem model.Node
	if err := c.ShouldBindJSON(&newItem); err != nil {
		HandleError(http.StatusBadRequest, c, err)
		return
	}
	newItem.Id = oldItem.Id

	log.Info("Post Node success")

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

func GetNodeTaskList(c *gin.Context) {

	tasks := TaskList

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    tasks,
	})
}

func DeleteNode(c *gin.Context) {
	id := bson.ObjectId("5d429e6c19f7abede924fee2")

	for _, node := range NodeList {
		if node.Id == id {
			log.Infof("Delete a node")
		}
	}
	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
	})
}

func GetSystemInfo(c *gin.Context) {
	id := c.Param("id")
	if !bson.IsObjectIdHex(id) {
		HandleErrorF(http.StatusBadRequest, c, "invalid id")
		return
	}
	sysInfo := systemInfo

	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Message: "success",
		Data:    sysInfo,
	})
}
