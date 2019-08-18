package database

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
	"testing"
)

func TestGetDb(t *testing.T) {
	Convey("Test GetDb", t, func() {
		s, db := GetDb()
		Convey("The value should be Session.Copy", func() {
			So(s, ShouldEqual, Session.Copy())
		})
		Convey("The value should be reference of database", func() {
			So(db, ShouldEqual, s.DB(viper.GetString("mongo.db")))
		})
	})
}
