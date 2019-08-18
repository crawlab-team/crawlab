package utils

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestEncryptPassword(t *testing.T) {
	var passwd = "test"
	Convey("Test EncryptPassword", t, func() {
		res := EncryptPassword(passwd)
		t.Log(res)
	})
}
