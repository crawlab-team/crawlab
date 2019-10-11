package mock

import (
	"bytes"
	"crawlab/model"
	"encoding/json"
	"github.com/globalsign/mgo/bson"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetSpiderList(t *testing.T) {
	var resp Response
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/spiders", nil)
	app.ServeHTTP(w, req)
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatal("unmarshal resp faild")
	}
	Convey("Test API GetSpiderList", t, func() {
		Convey("Test response status", func() {
			So(resp.Status, ShouldEqual, "ok")
			So(resp.Message, ShouldEqual, "success")
		})
	})
}

func TestGetSpider(t *testing.T) {
	var resp Response
	var spiderId = "5d429e6c19f7abede924fee2"
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/spiders/"+spiderId, nil)
	app.ServeHTTP(w, req)
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatal("unmarshal resp failed")
	}
	Convey("Test API GetSpider", t, func() {
		Convey("Test response status", func() {
			So(resp.Status, ShouldEqual, "ok")
			So(resp.Message, ShouldEqual, "success")
		})
	})
}

func TestPostSpider(t *testing.T) {
	var spider = model.Spider{
		Id:          bson.ObjectIdHex("5d429e6c19f7abede924fee2"),
		Name:        "For test",
		DisplayName: "test",
		Type:        "test",
		Col:         "test",
		Site:        "www.baidu.com",
		Envs:        nil,
		Src:         "/app/spider",
		Cmd:         "scrapy crawl test",
		LastRunTs:   time.Now(),
		CreateTs:    time.Now(),
		UpdateTs:    time.Now(),
	}
	var resp Response
	var spiderId = "5d429e6c19f7abede924fee2"
	w := httptest.NewRecorder()
	body, _ := json.Marshal(spider)
	req, _ := http.NewRequest("POST", "/spiders/"+spiderId, bytes.NewReader(body))
	app.ServeHTTP(w, req)
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatal("unmarshal resp failed")
	}
	Convey("Test API PostSpider", t, func() {
		Convey("Test response status", func() {
			So(resp.Status, ShouldEqual, "ok")
			So(resp.Message, ShouldEqual, "success")
		})
	})

}

func TestGetSpiderDir(t *testing.T) {
	var spiderId = "5d429e6c19f7abede924fee2"
	var resp Response
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/spiders/"+spiderId+"/dir", nil)
	app.ServeHTTP(w, req)
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatal("unmarshal resp failed")
	}
	Convey("Test API GetSpiderDir", t, func() {
		Convey("Test response status", func() {
			So(resp.Status, ShouldEqual, "ok")
			So(resp.Message, ShouldEqual, "success")
		})
	})

}

func TestGetSpiderTasks(t *testing.T) {
	var spiderId = "5d429e6c19f7abede924fee2"
	var resp Response
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/spiders/"+spiderId+"/tasks", nil)
	app.ServeHTTP(w, req)
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatal("unmarshal resp failed")
	}
	Convey("Test API GetSpiderTasks", t, func() {
		Convey("Test response status", func() {
			So(resp.Status, ShouldEqual, "ok")
			So(resp.Message, ShouldEqual, "success")
		})
	})
}

func TestDeleteSpider(t *testing.T) {
	var spiderId = "5d429e6c19f7abede924fee2"
	var resp Response
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/spiders/"+spiderId, nil)
	app.ServeHTTP(w, req)
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatal("unmarshal resp failed")
	}
	Convey("Test API DeleteSpider", t, func() {
		Convey("Test response status", func() {
			So(resp.Status, ShouldEqual, "ok")
			So(resp.Message, ShouldEqual, "success")
		})
	})
}
