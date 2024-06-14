package notification

import (
	"net/http"
	"testing"
	"time"
)

func TestService_sendMobile(t *testing.T) {
	T.Setup(t)
	e := T.NewExpect(t)
	time.Sleep(1 * time.Second)

	data := map[string]interface{}{
		"task_id": T.TestTask.GetId().Hex(),
	}
	e.POST("/send/mobile").WithJSON(data).
		Expect().Status(http.StatusOK)
}
