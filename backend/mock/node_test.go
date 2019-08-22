package mock

import (
	"crawlab/model"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"ucloudBilling/ucloud/log"
)

var app *gin.Engine
// 本测试依赖MongoDB的服务，所以在测试之前需要启动MongoDB及相关服务
func init() {
	app = gin.Default()

	// mock Test
	// 节点相关的API
	app.GET("/ping", Ping)
	app.GET("/nodes", GetNodeList)               // 节点列表
	app.GET("/nodes/:id", GetNode)               // 节点详情
	app.POST("/nodes/:id", PostNode)             // 修改节点
	app.GET("/nodes/:id/tasks", GetNodeTaskList) // 节点任务列表
	app.GET("/nodes/:id/system", GetSystemInfo)  // 节点任务列表
	app.DELETE("/nodes/:id", DeleteNode)         // 删除节点
	//// 爬虫
	// 定时任务
	app.GET("/schedules", GetScheduleList)       // 定时任务列表
	app.GET("/schedules/:id", GetSchedule)       // 定时任务详情
	app.PUT("/schedules", PutSchedule)           // 创建定时任务
	app.POST("/schedules/:id", PostSchedule)     // 修改定时任务
	app.DELETE("/schedules/:id", DeleteSchedule) // 删除定时任务
}

//mock test, test data in ./mock
func TestGetNodeList(t *testing.T) {
	var resp Response
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/nodes", nil)
	app.ServeHTTP(w, req)
	err := json.Unmarshal([]byte(w.Body.String()), &resp)
	t.Log(resp.Data)
	if err != nil {
		t.Fatal("Unmarshal resp failed")
	}

	Convey("Test API GetNodeList", t, func() {
		Convey("Test response status", func() {
			So(resp.Status, ShouldEqual, "ok")
			So(resp.Message, ShouldEqual, "success")
		})
	})
}

func TestGetNode(t *testing.T) {
	var resp Response
	var mongoId = "5d429e6c19f7abede924fee2"
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/nodes/"+mongoId, nil)
	app.ServeHTTP(w, req)
	err := json.Unmarshal([]byte(w.Body.String()), &resp)
	if err != nil {
		t.Fatal("Unmarshal resp failed")
	}
	t.Log(resp.Data)
	Convey("Test API GetNode", t, func() {
		Convey("Test response status", func() {
			So(resp.Status, ShouldEqual, "ok")
			So(resp.Message, ShouldEqual, "success")
			So(resp.Data.(map[string]interface{})["_id"], ShouldEqual, bson.ObjectId(mongoId).Hex())
		})
	})
}

func TestPing(t *testing.T) {
	var resp Response
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	app.ServeHTTP(w, req)
	err := json.Unmarshal([]byte(w.Body.String()), &resp)
	if err != nil {
		t.Fatal("Unmarshal resp failed")
	}
	Convey("Test API ping", t, func() {
		Convey("Test response status", func() {
			So(resp.Status, ShouldEqual, "ok")
			So(resp.Message, ShouldEqual, "success")
		})
	})
}

func TestGetNodeTaskList(t *testing.T) {
	var resp Response
	var mongoId = "5d429e6c19f7abede924fee2"
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "nodes/"+mongoId+"/tasks", nil)
	app.ServeHTTP(w, req)
	err := json.Unmarshal([]byte(w.Body.String()), &resp)
	if err != nil {
		t.Fatal("Unmarshal resp failed")
	}
	Convey("Test API GetNodeTaskList", t, func() {
		Convey("Test response status", func() {
			So(resp.Status, ShouldEqual, "ok")
			So(resp.Message, ShouldEqual, "success")
		})
	})
}

func TestDeleteNode(t *testing.T) {
	var resp Response

	var mongoId = "5d429e6c19f7abede924fee2"
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "nodes/"+mongoId, nil)
	app.ServeHTTP(w, req)
	err := json.Unmarshal([]byte(w.Body.String()), &resp)
	if err != nil {
		t.Fatal("Unmarshal resp failed")
	}
	Convey("Test API DeleteNode", t, func() {
		Convey("Test response status", func() {
			So(resp.Status, ShouldEqual, "ok")
			So(resp.Message, ShouldEqual, "success")
		})
	})
}

func TestPostNode(t *testing.T) {
	var newItem = model.Node{
		Id:           bson.ObjectIdHex("5d429e6c19f7abede924fee2"),
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
	}

	var resp Response
	body, _ := json.Marshal(newItem)
	log.Info(strings.NewReader(string(body)))

	var mongoId = "5d429e6c19f7abede924fee2"
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "nodes/"+mongoId, strings.NewReader(string(body)))

	app.ServeHTTP(w, req)
	err := json.Unmarshal([]byte(w.Body.String()), &resp)
	t.Log(resp)
	if err != nil {
		t.Fatal("Unmarshal resp failed")
	}
	Convey("Test API PostNode", t, func() {
		Convey("Test response status", func() {
			So(resp.Status, ShouldEqual, "ok")
			So(resp.Message, ShouldEqual, "success")
		})
	})
}

func TestGetSystemInfo(t *testing.T) {
	var resp Response
	var mongoId = "5d429e6c19f7abede924fee2"
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "nodes/"+mongoId+"/system", nil)
	app.ServeHTTP(w, req)
	err := json.Unmarshal([]byte(w.Body.String()), &resp)
	if err != nil {
		t.Fatal("Unmarshal resp failed")
	}
	Convey("Test API GetSystemInfo", t, func() {
		Convey("Test response status", func() {
			So(resp.Status, ShouldEqual, "ok")
			So(resp.Message, ShouldEqual, "success")
		})
	})
}
