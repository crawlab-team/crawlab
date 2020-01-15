package notification

import (
	"errors"
	"github.com/apex/log"
	"github.com/imroc/req"
	"runtime/debug"
)

func SendDingTalkNotification(webhook string, title string, content string) error {
	type ResBody struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}

	// 请求头
	header := req.Header{
		"Content-Type": "application/json; charset=utf-8",
	}

	// 请求数据
	data := req.Param{
		"msgtype": "markdown",
		"markdown": req.Param{
			"title": title,
			"text":  content,
		},
		"at": req.Param{
			"atMobiles": []string{},
			"isAtAll":   false,
		},
	}

	// 发起请求
	res, err := req.Post(webhook, header, req.BodyJSON(&data))
	if err != nil {
		log.Errorf("dingtalk notification error: " + err.Error())
		debug.PrintStack()
		return err
	}

	// 解析响应
	var resBody ResBody
	if err := res.ToJSON(&resBody); err != nil {
		log.Errorf("dingtalk notification error: " + err.Error())
		debug.PrintStack()
		return err
	}

	// 判断响应是否报错
	if resBody.ErrCode != 0 {
		log.Errorf("dingtalk notification error: " + resBody.ErrMsg)
		debug.PrintStack()
		return errors.New(resBody.ErrMsg)
	}

	return nil
}
