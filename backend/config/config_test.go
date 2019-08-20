package config

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestInitConfig(t *testing.T) {
	Convey("Test InitConfig func", t, func() {
		x := InitConfig("")

		Convey("The value should be nil", func() {
			So(x, ShouldEqual, nil)
		})
	})
}
