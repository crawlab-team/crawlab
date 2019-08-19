package database

import (
	"crawlab/config"
	"github.com/apex/log"
	"github.com/globalsign/mgo"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
	"reflect"
	"testing"
)

func init() {
	if err := config.InitConfig("../conf/config.yml"); err != nil {
		log.Fatal("Init config failed")
	}
	log.Infof("初始化配置成功")
	err := InitMongo()
	if err != nil {
		log.Fatal("Init mongodb failed")
	}

}

func TestGetDb(t *testing.T) {
	Convey("Test GetDb", t, func() {
		if err := config.InitConfig("../conf/config.yml"); err != nil {
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

func TestGetCol(t *testing.T) {
	var c = "nodes"
	var colActual *mgo.Collection
	Convey("Test GetCol", t, func() {
		s, col := GetCol(c)
		Convey("s should resemble Session.Copy", func() {
			So(s, ShouldResemble, Session.Copy())
			So(reflect.TypeOf(col), ShouldResemble, reflect.TypeOf(colActual))
		})
	})
}

func TestGetGridFs(t *testing.T) {
	var prefix = "files"
	var gfActual *mgo.GridFS

	Convey("Test GetGridFs", t, func() {
		s, gf := GetGridFs(prefix)
		Convey("s should be session.copy", func() {
			So(s, ShouldResemble, Session.Copy())
		})
		Convey("gf should be *mgo.GridFS", func() {
			So(reflect.TypeOf(gf), ShouldResemble, reflect.TypeOf(gfActual))
		})
	})
}
