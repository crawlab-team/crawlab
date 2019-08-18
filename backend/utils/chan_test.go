package utils

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewChanMap(t *testing.T) {
	mapTest := make(map[string]chan string)
	chanTest := make(chan string)
	test := "test"

	Convey("Call NewChanMap to generate ChanMap", t, func() {
		mapTest[test] = chanTest
		chanMapTest := ChanMap{mapTest}
		chanMap := NewChanMap()
		chanMap.m[test] = chanTest

		Convey(test, func() {
			So(chanMap, ShouldResemble, &chanMapTest)
		})

	})
}

func TestChan(t *testing.T) {
	mapTest := make(map[string]chan string)
	chanTest := make(chan string)
	mapTest["test"] = chanTest
	chanMapTest := ChanMap{mapTest}

	Convey("Test Chan use exist key", t, func() {
		ch1 := chanMapTest.Chan(
			"test")
		Convey("ch1 should equal chanTest", func() {
			So(ch1, ShouldEqual, chanTest)
		})

	})
	Convey("Test Chan use no-exist key", t, func() {
		ch2 := chanMapTest.Chan("test2")
		Convey("ch2 should equal chanMapTest.m[test2]", func() {

			So(chanMapTest.m["test2"], ShouldEqual, ch2)
		})
		Convey("Cap of chanMapTest.m[test2] should equal 10", func() {
			So(10, ShouldEqual, cap(chanMapTest.m["test2"]))
		})
	})
}

func TestChanBlocked(t *testing.T) {
	mapTest := make(map[string]chan string)
	chanTest := make(chan string)
	mapTest["test"] = chanTest
	chanMapTest := ChanMap{mapTest}

	Convey("Test Chan use exist key", t, func() {
		ch1 := chanMapTest.ChanBlocked(
			"test")
		Convey("ch1 should equal chanTest", func() {
			So(ch1, ShouldEqual, chanTest)
		})

	})
	Convey("Test Chan use no-exist key", t, func() {
		ch2 := chanMapTest.ChanBlocked("test2")
		Convey("ch2 should equal chanMapTest.m[test2]", func() {

			So(chanMapTest.m["test2"], ShouldEqual, ch2)
		})
		Convey("Cap of chanMapTest.m[test2] should equal 10", func() {
			So(0, ShouldEqual, cap(chanMapTest.m["test2"]))
		})
	})
}
