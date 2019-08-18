package utils

import (
	"github.com/globalsign/mgo/bson"
	. "github.com/smartystreets/goconvey/convey"
	"strconv"
	"testing"
	"time"
)

func TestIsObjectIdNull(t *testing.T) {
	var id bson.ObjectId = "123455"
	Convey("Test Object ID is null or not", t, func() {
		res := IsObjectIdNull(id)
		So(res, ShouldEqual, false)
	})
}

func TestInterfaceToString(t *testing.T) {
	var valueBson bson.ObjectId = "12345"
	var valueString = "12345"
	var valueInt = 12345
	var valueTime = time.Now().Add(60 * time.Second)
	var valueOther = []string{"a", "b"}

	Convey("Test InterfaceToString", t, func() {
		resBson := InterfaceToString(valueBson)
		Convey("resBson should be string value", func() {
			So(resBson, ShouldEqual, valueBson.Hex())
		})
		resString := InterfaceToString(valueString)
		Convey("resString should be string value", func() {
			So(resString, ShouldEqual, valueString)
		})
		resInt := InterfaceToString(valueInt)
		Convey("resInt should be string value", func() {
			So(resInt, ShouldEqual, strconv.Itoa(valueInt))
		})
		resTime := InterfaceToString(valueTime)
		Convey("resTime should be string value", func() {
			So(resTime, ShouldEqual, valueTime.String())
		})
		resOther := InterfaceToString(valueOther)
		Convey("resOther should be empty string", func() {
			So(resOther, ShouldEqual, "")
		})
	})

}
