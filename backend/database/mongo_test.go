package database

import (
	"crawlab/config"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
	"testing"
)

func TestGetDb(t *testing.T) {
	Convey("Test GetDb", t, func() {
		if err := config.InitConfig(""); err != nil {
			t.Fatal("Init config failed")
		}
		t.Log("初始化配置成功")
		err := InitMongo()
		if err != nil {
			t.Fatal("Init mongodb failed")
		}
		s, db := GetDb()
		Convey("The value should be Session.Copy", func() {
			So(s, ShouldResemble, Session.Copy())
		})
		Convey("The value should be reference of database", func() {
			So(db, ShouldResemble, s.DB(viper.GetString("mongo.db")))
		})
	})
}
