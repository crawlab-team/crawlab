package notification

import (
	"errors"
	"github.com/crawlab-team/go-trace"
	"github.com/imroc/req"
	"strings"
)

type ResBody struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func SendMobileNotification(webhook string, title string, content string) error {
	// request header
	header := req.Header{
		"Content-Type": "application/json; charset=utf-8",
	}

	// request data
	data := req.Param{
		"msgtype": "markdown",
		"markdown": req.Param{
			"title":   title,
			"text":    content,
			"content": content,
		},
		"at": req.Param{
			"atMobiles": []string{},
			"isAtAll":   false,
		},
		"text": content,
	}
	if strings.Contains(strings.ToLower(webhook), "feishu") {
		data = req.Param{
			"msg_type": "text",
			"content": req.Param{
				"text": content,
			},
		}
	}

	// perform request
	res, err := req.Post(webhook, header, req.BodyJSON(&data))
	if err != nil {
		return trace.TraceError(err)
	}

	// parse response
	var resBody ResBody
	if err := res.ToJSON(&resBody); err != nil {
		return trace.TraceError(err)
	}

	// validate response code
	if resBody.ErrCode != 0 {
		return errors.New(resBody.ErrMsg)
	}

	return nil
}
