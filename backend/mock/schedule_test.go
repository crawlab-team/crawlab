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
	"ucloudBilling/ucloud/log"
)

func TestGetScheduleList(t *testing.T) {
	var resp Response
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/schedules", nil)
	app.ServeHTTP(w, req)
	err := json.Unmarshal([]byte(w.Body.String()), &resp)
	if err != nil {
		t.Fatal("Unmarshal resp failed")
	}
	t.Log(resp.Data)
	Convey("Test API GetScheduleList", t, func() {
		Convey("Test response status", func() {
			So(resp.Status, ShouldEqual, "ok")
			So(resp.Message, ShouldEqual, "success")
		})
	})
}

func TestGetSchedule(t *testing.T) {
	var mongoId = "5d429e6c19f7abede924fee2"
	var resp Response
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/schedules/"+mongoId, nil)
	app.ServeHTTP(w, req)
	err := json.Unmarshal([]byte(w.Body.String()), &resp)
	if err != nil {
		t.Fatal("Unmarshal resp failed")
	}
	Convey("Test API GetSchedule", t, func() {
		Convey("Test response status", func() {
			So(resp.Status, ShouldEqual, "ok")
			So(resp.Message, ShouldEqual, "success")
			So(resp.Data.(map[string]interface{})["_id"], ShouldEqual, bson.ObjectId(mongoId).Hex())
		})
	})
}

func TestDeleteSchedule(t *testing.T) {
	var mongoId = "5d429e6c19f7abede924fee2"
	var resp Response
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/schedules/"+mongoId, nil)
	app.ServeHTTP(w, req)

	err := json.Unmarshal([]byte(w.Body.String()), &resp)
	log.Info(w.Body.String())
	if err != nil {
		t.Fatal("Unmarshal resp failed")
	}

	Convey("Test DeleteSchedule", t, func() {
		Convey("Test resp status", func() {
			So(resp.Status, ShouldEqual, "ok")
		})
	})
}

func TestPostSchedule(t *testing.T) {
	var newItem = model.Schedule{
		Id:       bson.ObjectIdHex("5d429e6c19f7abede924fee2"),
		Name:     "test schedule",
		SpiderId: bson.ObjectIdHex("5d429e6c19f7abede924fee2"),
		NodeId:   bson.ObjectIdHex("5d429e6c19f7abede924fee2"),
		Cron:     "***1*",
		EntryId:  10,
		// 前端展示
		SpiderName: "test scedule",
		NodeName:   "测试节点",

		CreateTs: time.Now(),
		UpdateTs: time.Now(),
	}

	var resp Response
	var mongoId = "5d429e6c19f7abede924fee2"
	body,_ := json.Marshal(newItem)
	log.Info(strings.NewReader(string(body)))
	w := httptest.NewRecorder()
	req,_ := http.NewRequest("POST", "/schedules/"+mongoId,strings.NewReader(string(body)))
	app.ServeHTTP(w, req)
	err := json.Unmarshal([]byte(w.Body.String()),&resp)
	t.Log(resp)
	if err != nil {
		t.Fatal("unmarshal resp failed")
	}
	Convey("Test API PostSchedule", t, func() {
		Convey("Test response status", func() {
			So(resp.Status, ShouldEqual, "ok")
			So(resp.Message, ShouldEqual, "success")
		})
	})

}

func TestPutSchedule(t *testing.T) {
	var newItem = model.Schedule{
		Id:       bson.ObjectIdHex("5d429e6c19f7abede924fee2"),
		Name:     "test schedule",
		SpiderId: bson.ObjectIdHex("5d429e6c19f7abede924fee2"),
		NodeId:   bson.ObjectIdHex("5d429e6c19f7abede924fee2"),
		Cron:     "***1*",
		EntryId:  10,
		// 前端展示
		SpiderName: "test scedule",
		NodeName:   "测试节点",

		CreateTs: time.Now(),
		UpdateTs: time.Now(),
	}

	var resp Response
	body,_ := json.Marshal(newItem)
	log.Info(strings.NewReader(string(body)))
	w := httptest.NewRecorder()
	req,_ := http.NewRequest("PUT", "/schedules",strings.NewReader(string(body)))
	app.ServeHTTP(w, req)
	err := json.Unmarshal([]byte(w.Body.String()),&resp)
	t.Log(resp)
	if err != nil {
		t.Fatal("unmarshal resp failed")
	}
	Convey("Test API PutSchedule", t, func() {
		Convey("Test response status", func() {
			So(resp.Status, ShouldEqual, "ok")
			So(resp.Message, ShouldEqual, "success")
		})
	})

}
