package notification

import (
	"errors"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab/core/models/models/v2"
	"github.com/crawlab-team/crawlab/trace"
	"github.com/imroc/req"
	"strings"
)

type ResBody struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func SendIMNotification(s *models.NotificationSettingV2, ch *models.NotificationChannelV2, title, content string) error {
	// TODO: compatibility with different IM providers
	switch ch.Provider {
	case ChannelIMProviderLark:
		return sendImLark(ch, title, content)
	case ChannelIMProviderSlack:
		return sendImSlack(ch, title, content)
	}

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

func getIMRequestHeader() req.Header {
	return req.Header{
		"Content-Type": "application/json; charset=utf-8",
	}
}

func performIMRequest(webhookUrl string, data req.Param) error {
	// perform request
	res, err := req.Post(webhookUrl, getIMRequestHeader(), req.BodyJSON(&data))
	if err != nil {
		log.Errorf("IM request error: %v", err)
		return err
	}

	// parse response
	var resBody ResBody
	if err := res.ToJSON(&resBody); err != nil {
		log.Errorf("Parsing IM response error: %v", err)
		return err
	}

	// validate response code
	if resBody.ErrCode != 0 {
		log.Errorf("IM response error: %v", resBody.ErrMsg)
		return errors.New(resBody.ErrMsg)
	}

	return nil
}

func sendImLark(ch *models.NotificationChannelV2, title, content string) error {
	// request header
	data := req.Param{
		"msg_type": "interactive",
		"card": req.Param{
			"header": req.Param{
				"title": req.Param{
					"tag":     "plain_text",
					"content": title,
				},
			},
			"elements": []req.Param{
				{
					"tag":     "markdown",
					"content": content,
				},
			},
		},
	}
	return performIMRequest(ch.WebhookUrl, data)
}

func sendImSlack(ch *models.NotificationChannelV2, title, content string) error {
	// request header
	data := req.Param{
		"blocks": []req.Param{
			{"type": "header", "text": req.Param{"type": "plain_text", "text": title}},
			{"type": "section", "text": req.Param{"type": "mrkdwn", "text": content}},
		},
	}
	return performIMRequest(ch.WebhookUrl, data)
}
