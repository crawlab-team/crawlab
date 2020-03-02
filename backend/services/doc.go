package services

import (
	"github.com/apex/log"
	"github.com/imroc/req"
	"runtime/debug"
)

func GetDocs() (data string, err error) {
	// 获取远端数据
	res, err := req.Get("https://docs.crawlab.cn/search_plus_index.json")
	if err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return data, err
	}

	// 反序列化
	data, err = res.ToString()
	if err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return data, err
	}

	return data, nil
}
