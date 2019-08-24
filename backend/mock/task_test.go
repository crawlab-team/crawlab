package mock

import (
	"crawlab/model"
	"encoding/json"
	"github.com/globalsign/mgo/bson"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGetTaskList(t *testing.T) {
	//var teskListRequestFrom = TaskListRequestData{
	//	PageNum:  2,
	//	PageSize: 10,
	//	NodeId:   "434221grfsf",
	//	SpiderId: "fdfewqrftea",
	//}

	var resp ListResponse
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tasks?PageNum=2&PageSize=10&NodeId=342dfsff&SpiderId=f8dsf", nil)
	app.ServeHTTP(w, req)
	err := json.Unmarshal([]byte(w.Body.String()), &resp)
	if err != nil {
		t.Fatal("Unmarshal resp failed")
	}

	Convey("Test API GetNodeList", t, func() {
		Convey("Test response status", func() {
			So(resp.Status, ShouldEqual, "ok")
			So(resp.Message, ShouldEqual, "success")
			So(resp.Total, ShouldEqual, 2)
		})
	})
}

func TestGetTask(t *testing.T) {
	var resp Response
	var taskId = "1234"
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tasks/"+taskId, nil)
	app.ServeHTTP(w, req)
	err := json.Unmarshal([]byte(w.Body.String()), &resp)
	if err != nil {
		t.Fatal("Unmarshal resp failed")
	}
	Convey("Test API GetTask", t, func() {
		Convey("Test response status", func() {
			So(resp.Status, ShouldEqual, "ok")
			So(resp.Message, ShouldEqual, "success")
		})
	})
}

func TestPutTask(t *testing.T) {
	var newItem = model.Task{
		Id:              "1234",
		SpiderId:        bson.ObjectIdHex("5d429e6c19f7abede924fee2"),
		StartTs:         time.Now(),
		FinishTs:        time.Now(),
		Status:          "online",
		NodeId:          bson.ObjectIdHex("5d429e6c19f7abede924fee2"),
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
	}

	var resp Response
	body, _ := json.Marshal(&newItem)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/tasks", strings.NewReader(string(body)))
	app.ServeHTTP(w, req)
	err := json.Unmarshal([]byte(w.Body.String()), &resp)
	if err != nil {
		t.Fatal("unmarshal resp failed")
	}
	Convey("Test API PutTask", t, func() {
		Convey("Test response status", func() {
			So(resp.Status, ShouldEqual, "ok")
			So(resp.Message, ShouldEqual, "success")
		})
	})
}

func TestDeleteTask(t *testing.T) {
	taskId := "1234"
	var resp Response
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/tasks/"+taskId, nil)
	app.ServeHTTP(w, req)
	err := json.Unmarshal([]byte(w.Body.String()), &resp)
	if err != nil {
		t.Fatal("unmarshal resp failed")
	}
	Convey("Test API DeleteTask", t, func() {
		Convey("Test response status", func() {
			So(resp.Status, ShouldEqual, "ok")
			So(resp.Message, ShouldEqual, "success")
		})
	})
}

func TestGetTaskResults(t *testing.T) {
	//var teskListResultFrom = TaskResultsRequestData{
	//	PageNum:  2,
	//	PageSize: 1,
	//}
	taskId := "1234"

	var resp ListResponse
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tasks/"+taskId+"/results?PageNum=2&PageSize=1", nil)
	app.ServeHTTP(w, req)
	err := json.Unmarshal([]byte(w.Body.String()), &resp)
	if err != nil {
		t.Fatal("Unmarshal resp failed")
	}

	Convey("Test API GetNodeList", t, func() {
		Convey("Test response status", func() {
			So(resp.Status, ShouldEqual, "ok")
			So(resp.Message, ShouldEqual, "success")
			So(resp.Total, ShouldEqual, 2)
		})
	})
}
