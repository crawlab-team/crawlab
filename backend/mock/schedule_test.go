package mock

import (
	"bytes"
	"crawlab/model"
	"crawlab/utils"
	"encoding/json"
	"github.com/globalsign/mgo/bson"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGetScheduleList(t *testing.T) {
	var resp Response
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/schedules", nil)
	app.ServeHTTP(w, req)
	err := json.Unmarshal(w.Body.Bytes(), &resp)
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
	err := json.Unmarshal(w.Body.Bytes(), &resp)
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

	err := json.Unmarshal(w.Body.Bytes(), &resp)
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
		NodeIds:  NodeIdss,
		Cron:     "***1*",
		EntryId:  10,
		// 前端展示
		SpiderName: "test scedule",

		CreateTs: time.Now(),
		UpdateTs: time.Now(),
	}

	var resp Response
	var mongoId = "5d429e6c19f7abede924fee2"
	body, _ := json.Marshal(newItem)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/schedules/"+mongoId, strings.NewReader(utils.BytesToString(body)))
	app.ServeHTTP(w, req)

	err := json.Unmarshal(w.Body.Bytes(), &resp)
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
		NodeIds:  NodeIdss,
		Cron:     "***1*",
		EntryId:  10,
		// 前端展示
		SpiderName: "test scedule",

		CreateTs: time.Now(),
		UpdateTs: time.Now(),
	}

	var resp Response
	body, _ := json.Marshal(newItem)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/schedules", bytes.NewReader(body))
	app.ServeHTTP(w, req)
	err := json.Unmarshal(w.Body.Bytes(), &resp)
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
