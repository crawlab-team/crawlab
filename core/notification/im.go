package notification

import (
	"errors"
	"github.com/crawlab-team/crawlab/core/models/models/v2"
	"github.com/crawlab-team/crawlab/trace"
	"github.com/imroc/req"
	"strings"
)

type ResBody struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func SendIMNotification(s *models.NotificationSettingV2, ch *models.NotificationChannelV2, content string) error {
	// TODO: compatibility with different IM providers

	// request header
	header := req.Header{
		"Content-Type": "application/json; charset=utf-8",
	}

	// request data
	data := req.Param{
		"msgtype": "markdown",
		"markdown": req.Param{
			"title":   s.Title,
			"text":    content,
			"content": content,
		},
		"at": req.Param{
			"atMobiles": []string{},
			"isAtAll":   false,
		},
		"text": content,
	}
	if strings.Contains(strings.ToLower(ch.WebhookUrl), "feishu") {
		data = req.Param{
			"msg_type": "text",
			"content": req.Param{
				"text": content,
			},
		}
	}

	// perform request
	res, err := req.Post(ch.WebhookUrl, header, req.BodyJSON(&data))
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
