package utils

import (
	. "github.com/smartystreets/goconvey/convey"
	"sync"
	"testing"
)

func TestNewChanMap(t *testing.T) {
	mapTest := sync.Map{}
	chanTest := make(chan string)
	test := "test"

	Convey("Call NewChanMap to generate ChanMap", t, func() {
		mapTest.Store("test", chanTest)
		chanMapTest := ChanMap{mapTest}
		chanMap := NewChanMap()
		chanMap.m.Store("test", chanTest)

		Convey(test, func() {
			v1, ok := chanMap.m.Load("test")
			So(ok, ShouldBeTrue)
			v2, ok := chanMapTest.m.Load("test")
			So(ok, ShouldBeTrue)
			So(v1, ShouldResemble, v2)
		})
	})
}

func TestChan(t *testing.T) {
	mapTest := sync.Map{}
	chanTest := make(chan string)
	mapTest.Store("test", chanTest)
	chanMapTest := ChanMap{mapTest}

	Convey("Test Chan use exist key", t, func() {
		ch1 := chanMapTest.Chan("test")
		Convey("ch1 should equal chanTest", func() {
			So(ch1, ShouldEqual, chanTest)
		})
	})
	Convey("Test Chan use no-exist key", t, func() {
		ch2 := chanMapTest.Chan("test2")
		Convey("ch2 should equal chanMapTest.m[test2]", func() {
			v, ok := chanMapTest.m.Load("test2")
			So(ok, ShouldBeTrue)
			So(v, ShouldEqual, ch2)
		})
		Convey("Cap of chanMapTest.m[test2] should equal 10", func() {
			So(10, ShouldEqual, cap(ch2))
		})
	})
}

func TestChanBlocked(t *testing.T) {
	mapTest := sync.Map{}
	chanTest := make(chan string)
	mapTest.Store("test", chanTest)
	chanMapTest := ChanMap{mapTest}

	Convey("Test Chan use exist key", t, func() {
		ch1 := chanMapTest.ChanBlocked("test")
		Convey("ch1 should equal chanTest", func() {
			So(ch1, ShouldEqual, chanTest)
		})
	})
	Convey("Test Chan use no-exist key", t, func() {
		ch2 := chanMapTest.ChanBlocked("test2")
		Convey("ch2 should equal chanMapTest.m[test2]", func() {
			v, ok := chanMapTest.m.Load("test2")
			So(ok, ShouldBeTrue)
			So(v, ShouldEqual, ch2)
		})
		Convey("Cap of chanMapTest.m[test2] should equal 10", func() {
			So(0, ShouldEqual, cap(ch2))
		})
	})
}
