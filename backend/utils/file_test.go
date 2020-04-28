package utils

import (
	"archive/zip"
	. "github.com/smartystreets/goconvey/convey"
	"io"
	"log"
	"os"
	"runtime/debug"
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
	err := os.Mkdir("testCompress", os.ModePerm)
	if err != nil {
		t.Error("create testCompress failed")
	}
	var pathString = "testCompress"
	var files []*os.File
	var disPath = "testCompress"
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
	_ = os.RemoveAll("testCompress")

}
func Zip(zipFile string, fileList []string) error {
	// 创建 zip 包文件
	fw, err := os.Create(zipFile)
	if err != nil {
		log.Fatal()
	}
	defer Close(fw)

	// 实例化新的 zip.Writer
	zw := zip.NewWriter(fw)
	defer Close(zw)

	for _, fileName := range fileList {
		fr, err := os.Open(fileName)
		if err != nil {
			return err
		}
		fi, err := fr.Stat()
		if err != nil {
			return err
		}
		// 写入文件的头信息
		fh, err := zip.FileInfoHeader(fi)
		if err != nil {
			return err
		}
		w, err := zw.CreateHeader(fh)
		if err != nil {
			return err
		}
		// 写入文件内容
		_, err = io.Copy(w, fr)
		if err != nil {
			return err
		}
	}
	return nil
}

func TestDeCompress(t *testing.T) {
	err := os.Mkdir("testDeCompress", os.ModePerm)
	if err != nil {
		t.Error(err)

	}
	err = Zip("demo.zip", []string{})
	if err != nil {
		t.Error("create zip file failed")
	}
	tmpFile, err := os.OpenFile("demo.zip", os.O_RDONLY, 0777)
	if err != nil {
		debug.PrintStack()
		t.Error("open demo.zip failed")
	}
	var dstPath = "./testDeCompress"
	Convey("Test DeCopmress func", t, func() {

		err := DeCompress(tmpFile, dstPath)
		So(err, ShouldEqual, nil)
	})
	_ = os.RemoveAll("testDeCompress")
	_ = os.Remove("demo.zip")

}
