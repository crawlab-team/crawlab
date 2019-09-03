package utils

import (
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func TestExists(t *testing.T) {
	var pathString = "../config"
	var wrongPathString = "test"

	Convey("Test path or file is Exists or not", t, func() {
		res := Exists(pathString)
		Convey("The result should be true", func() {
			So(res, ShouldEqual, true)
		})
		wrongRes := Exists(wrongPathString)
		Convey("The result should be false", func() {
			So(wrongRes, ShouldEqual, false)
		})
	})
}

func TestIsDir(t *testing.T) {
	var pathString = "../config"
	var fileString = "../config/config.go"
	var wrongString = "test"

	Convey("Test path is folder or not", t, func() {
		res := IsDir(pathString)
		So(res, ShouldEqual, true)
		fileRes := IsDir(fileString)
		So(fileRes, ShouldEqual, false)
		wrongRes := IsDir(wrongString)
		So(wrongRes, ShouldEqual, false)
	})
}

func TestCompress(t *testing.T) {
	var pathString = "../utils"
	var files []*os.File
	var disPath = "../utils/test"
	file, err := os.Open(pathString)
	if err != nil {
		t.Error("open source path failed")
	}
	files = append(files, file)
	Convey("Verify dispath is valid path", t, func() {
		er := Compress(files, disPath)
		Convey("err should be nil", func() {
			So(er, ShouldEqual, nil)
		})
	})

}

// 测试之前需存在有效的test(.zip)文件
func TestDeCompress(t *testing.T) {
	var tmpFilePath = "./test"
	tmpFile, err := os.OpenFile(tmpFilePath, os.O_RDONLY, 0777)
	if err != nil {
		t.Fatal("open zip file failed")
	}
	var dstPath = "./testDeCompress"
	Convey("Test DeCopmress func", t, func() {

		err := DeCompress(tmpFile, dstPath)
		So(err, ShouldEqual, nil)
	})

}
