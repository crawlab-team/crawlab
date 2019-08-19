package routes

import (
	"crawlab/config"
	"crawlab/database"
	"encoding/json"
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
	"net/http"
	"net/http/httptest"
	"runtime/debug"
	"testing"
)

var app *gin.Engine
// 本测试依赖MongoDB的服务，所以在测试之前需要启动MongoDB及相关服务
func init() {
	app = gin.Default()

	// 初始化配置
	if err := config.InitConfig("../conf/config.yml"); err != nil {
		panic(err)
	}
	log.Info("初始化配置成功")

	// 初始化日志设置
	logLevel := viper.GetString("log.level")
	if logLevel != "" {
		log.SetLevelFromString(logLevel)
	}
	log.Info("初始化日志设置成功")

	// 初始化Mongodb数据库
	if err := database.InitMongo(); err != nil {
		debug.PrintStack()
		panic(err)
	}
	log.Info("初始化Mongodb数据库成功")

	// 初始化Redis数据库
	if err := database.InitRedis(); err != nil {
		debug.PrintStack()
		panic(err)
	}
	log.Info("初始化Redis数据库成功")

	// 路由
	// 节点相关的API
	app.GET("/ping", Ping)
	app.GET("/nodes", GetNodeList)               // 节点列表
	app.GET("/nodes/:id", GetNode)               // 节点详情
	app.POST("/nodes/:id", PostNode)             // 修改节点
	app.GET("/nodes/:id/tasks", GetNodeTaskList) // 节点任务列表
	app.GET("/nodes/:id/system", GetSystemInfo)  // 节点任务列表
	app.DELETE("/nodes/:id", DeleteNode)         // 删除节点
	// 爬虫
}

//先测试GetNodeList得到一个节点的ID，再继续测试下面的API
func TestGetNodeList(t *testing.T) {
	var resp Response
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/nodes", nil)
	app.ServeHTTP(w, req)
	err := json.Unmarshal([]byte(w.Body.String()), &resp)
	if err != nil {
		t.Fatal("Unmarshal resp failed")
	}
	t.Log(resp.Data)
	Convey("Test API GetNodeList", t, func() {
		Convey("Test response status", func() {
			So(resp.Status, ShouldEqual, "ok")
			So(resp.Message, ShouldEqual, "success")
		})
	})
}

//依赖MongoDB中的数据,_id=5d429e6c19f7abede924fee2,实际测试时需替换
func TestGetNode(t *testing.T) {
	var resp Response
	var mongoId = "5d5a658319f7ab423585b0b0"
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
			So(resp.Data.(map[string]interface{})["_id"], ShouldEqual, mongoId)
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
	var mongoId = "5d5a658319f7ab423585b0b0"
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
	var mongoId = "5d5a658319f7ab423585b0b0"
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
