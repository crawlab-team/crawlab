package mock

import (
	"encoding/json"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetHomeStats(t *testing.T) {
	var resp Response
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/stats/home", nil)
	app.ServeHTTP(w, req)
	err := json.Unmarshal([]byte(w.Body.String()), &resp)
	fmt.Println(resp.Data)
	if err != nil {
		t.Fatal("Unmarshal resp failed")
	}

	Convey("Test API GetHomeStats", t, func() {
		Convey("Test response status", func() {
			So(resp.Status, ShouldEqual, "ok")
			So(resp.Message, ShouldEqual, "success")
		})
	})
}