package utils

import "unsafe"

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
