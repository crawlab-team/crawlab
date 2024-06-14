package utils

import (
	"github.com/crawlab-team/go-trace"
	"io"
	"reflect"
	"unsafe"
)

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func Close(c io.Closer) {
	err := c.Close()
	if err != nil {
		trace.PrintError(err)
	}
}

func Contains(array interface{}, val interface{}) (fla bool) {
	fla = false
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		{
			s := reflect.ValueOf(array)
			for i := 0; i < s.Len(); i++ {
				if reflect.DeepEqual(val, s.Index(i).Interface()) {
					fla = true
					return
				}
			}
		}
	}
	return
}
