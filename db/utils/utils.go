package utils

import "io"

func Close(c io.Closer) {
	err := c.Close()
	if err != nil {
		//log.WithError(err).Error("关闭资源文件失败。")
	}
}

func ContainsString(list []string, item string) bool {
	for _, d := range list {
		if d == item {
			return true
		}
	}
	return false
}
