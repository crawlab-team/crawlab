package utils

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEncryptPassword(t *testing.T) {
	var passwd = "test"
	Convey("Test EncryptPassword", t, func() {
		res := EncryptPassword(passwd)
		t.Log(res)
	})
}
