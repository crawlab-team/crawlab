package model

import (
	"github.com/apex/log"
	"os"
	"runtime/debug"
)

// 获取本地日志
func GetLocalLog(logPath string) (fileBytes []byte, err error) {

	f, err := os.Open(logPath)
	if err != nil {
		log.Error(err.Error())
		debug.PrintStack()
		return nil, err
	}
	fi, err := f.Stat()
	if err != nil {
		log.Error(err.Error())
		debug.PrintStack()
		return nil, err
	}
	defer f.Close()

	const bufLen = 2 * 1024 * 1024
	logBuf := make([]byte, bufLen)

	off := int64(0)
	if fi.Size() > int64(len(logBuf)) {
		off = fi.Size() - int64(len(logBuf))
	}
	n, err := f.ReadAt(logBuf, off)

	//到文件结尾会有EOF标识
	if err != nil && err.Error() != "EOF" {
		log.Error(err.Error())
		debug.PrintStack()
		return nil, err
	}
	logBuf = logBuf[:n]
	return logBuf, nil
}
