package notification

import (
	"errors"
	"fmt"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab/core/models/models/v2"
	"github.com/crawlab-team/crawlab/trace"
	"github.com/imroc/req"
	"regexp"
	"strings"
)

type ResBody struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func SendIMNotification(ch *models.NotificationChannelV2, title, content string) error {
	// TODO: compatibility with different IM providers
	switch ch.Provider {
	case ChannelIMProviderLark:
		return sendImLark(ch, title, content)
	case ChannelIMProviderSlack:
		return sendImSlack(ch, title, content)
	case ChannelIMProviderDingtalk:
		return sendImDingTalk(ch, title, content)
	case ChannelIMProviderWechatWork:
		return sendImWechatWork(ch, title, content)
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

func convertMarkdownToSlack(markdown string) string {
	// Convert bold text
	reBold := regexp.MustCompile(`\*\*(.*?)\*\*`)
	slack := reBold.ReplaceAllString(markdown, `*$1*`)

	// Convert italic text
	reItalic := regexp.MustCompile(`\*(.*?)\*`)
	slack = reItalic.ReplaceAllString(slack, `_$1_`)

	// Convert links
	reLink := regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)
	slack = reLink.ReplaceAllString(slack, `<$2|$1>`)

	// Convert inline code
	reInlineCode := regexp.MustCompile("`(.*?)`")
	slack = reInlineCode.ReplaceAllString(slack, "`$1`")

	// Convert unordered list
	slack = strings.ReplaceAll(slack, "- ", "• ")

	// Convert ordered list
	reOrderedList := regexp.MustCompile(`^\d+\. `)
	slack = reOrderedList.ReplaceAllStringFunc(slack, func(s string) string {
		return strings.Replace(s, ". ", ". ", 1)
	})

	// Convert blockquote
	reBlockquote := regexp.MustCompile(`^> (.*)`)
	slack = reBlockquote.ReplaceAllString(slack, `> $1`)

	return slack
}

func sendImLark(ch *models.NotificationChannelV2, title, content string) error {
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
	data := req.Param{
		"blocks": []req.Param{
			{"type": "header", "text": req.Param{"type": "plain_text", "text": title}},
			{"type": "section", "text": req.Param{"type": "mrkdwn", "text": convertMarkdownToSlack(content)}},
		},
	}
	return performIMRequest(ch.WebhookUrl, data)
}

func sendImDingTalk(ch *models.NotificationChannelV2, title string, content string) error {
	data := req.Param{
		"msgtype": "markdown",
		"markdown": req.Param{
			"title": title,
			"text":  fmt.Sprintf("# %s\n\n%s", title, content),
		},
	}
	return performIMRequest(ch.WebhookUrl, data)
}

func sendImWechatWork(ch *models.NotificationChannelV2, title string, content string) error {
	data := req.Param{
		"msgtype": "markdown",
		"markdown": req.Param{
			"content": fmt.Sprintf("# %s\n\n%s", title, content),
		},
	}
	return performIMRequest(ch.WebhookUrl, data)
}
