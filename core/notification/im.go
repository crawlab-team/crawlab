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
	switch ch.Provider {
	case ChannelIMProviderLark:
		return sendIMLark(ch, title, content)
	case ChannelIMProviderDingtalk:
		return sendIMDingTalk(ch, title, content)
	case ChannelIMProviderWechatWork:
		return sendIMWechatWork(ch, title, content)
	case ChannelIMProviderSlack:
		return sendIMSlack(ch, title, content)
	case ChannelIMProviderTelegram:
		return sendIMTelegram(ch, title, content)
	case ChannelIMProviderDiscord:
		return sendIMDiscord(ch, title, content)
	case ChannelIMProviderMSTeams:
		return sendIMMSTeams(ch, title, content)
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

func performIMRequest(webhookUrl string, data req.Param) (res *req.Resp, err error) {
	// perform request
	res, err = req.Post(webhookUrl, getIMRequestHeader(), req.BodyJSON(&data))
	if err != nil {
		log.Errorf("IM request error: %v", err)
		return nil, err
	}

	// get response
	response := res.Response()

	// check status code
	if response.StatusCode >= 400 {
		log.Errorf("IM response status code: %d", res.Response().StatusCode)
		return nil, errors.New(fmt.Sprintf("IM error response %d: %s", response.StatusCode, res.String()))
	}

	return res, nil
}

func performIMRequestWithJson[T any](webhookUrl string, data req.Param) (resBody T, err error) {
	res, err := performIMRequest(webhookUrl, data)
	if err != nil {
		return resBody, err
	}

	// parse response
	if err := res.ToJSON(&resBody); err != nil {
		log.Warnf("Parsing IM response error: %v", err)
		resText, err := res.ToString()
		if err != nil {
			log.Warnf("Converting response to string error: %v", err)
			return resBody, err
		}
		log.Infof("IM response: %s", resText)
		return resBody, nil
	}

	return resBody, nil
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

func convertMarkdownToTelegram(markdownText string) string {
	// Combined regex to handle bold and italic
	re := regexp.MustCompile(`(?m)(\*\*)(.*)(\*\*)|(__)(.*)(__)|(\*)(.*)(\*)|(_)(.*)(_)`)
	markdownText = re.ReplaceAllStringFunc(markdownText, func(match string) string {
		groups := re.FindStringSubmatch(match)
		if groups[1] != "" || groups[4] != "" {
			// Handle bold
			return "*" + match[2:len(match)-2] + "*"
		} else if groups[6] != "" || groups[9] != "" {
			// Handle italic
			return "_" + match[1:len(match)-1] + "_"
		} else {
			// No match
			return match
		}
	})

	// Convert unordered list
	re = regexp.MustCompile(`(?m)^- (.*)`)
	markdownText = re.ReplaceAllString(markdownText, "• $1")

	// Escape characters
	escapeChars := []string{"#", "-", "."}
	for _, c := range escapeChars {
		markdownText = strings.ReplaceAll(markdownText, c, "\\"+c)
	}

	return markdownText
}

func sendIMLark(ch *models.NotificationChannelV2, title, content string) error {
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
	resBody, err := performIMRequestWithJson[ResBody](ch.WebhookUrl, data)
	if err != nil {
		return err
	}
	if resBody.ErrCode != 0 {
		return errors.New(resBody.ErrMsg)
	}
	return nil
}

func sendIMDingTalk(ch *models.NotificationChannelV2, title string, content string) error {
	data := req.Param{
		"msgtype": "markdown",
		"markdown": req.Param{
			"title": title,
			"text":  fmt.Sprintf("# %s\n\n%s", title, content),
		},
	}
	resBody, err := performIMRequestWithJson[ResBody](ch.WebhookUrl, data)
	if err != nil {
		return err
	}
	if resBody.ErrCode != 0 {
		return errors.New(resBody.ErrMsg)
	}
	return nil
}

func sendIMWechatWork(ch *models.NotificationChannelV2, title string, content string) error {
	data := req.Param{
		"msgtype": "markdown",
		"markdown": req.Param{
			"content": fmt.Sprintf("# %s\n\n%s", title, content),
		},
	}
	resBody, err := performIMRequestWithJson[ResBody](ch.WebhookUrl, data)
	if err != nil {
		return err
	}
	if resBody.ErrCode != 0 {
		return errors.New(resBody.ErrMsg)
	}
	return nil
}

func sendIMSlack(ch *models.NotificationChannelV2, title, content string) error {
	data := req.Param{
		"blocks": []req.Param{
			{"type": "header", "text": req.Param{"type": "plain_text", "text": title}},
			{"type": "section", "text": req.Param{"type": "mrkdwn", "text": convertMarkdownToSlack(content)}},
		},
	}
	_, err := performIMRequest(ch.WebhookUrl, data)
	if err != nil {
		return err
	}
	return nil
}

func sendIMTelegram(ch *models.NotificationChannelV2, title string, content string) error {
	type ResBody struct {
		Ok          bool   `json:"ok"`
		Description string `json:"description"`
	}

	// chat id
	chatId := ch.TelegramChatId
	if !strings.HasPrefix("@", ch.TelegramChatId) {
		chatId = fmt.Sprintf("@%s", ch.TelegramChatId)
	}

	// webhook url
	webhookUrl := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", ch.TelegramBotToken)

	// original Markdown text
	text := fmt.Sprintf("**%s**\n\n%s", title, content)

	// convert to Telegram MarkdownV2
	text = convertMarkdownToTelegram(text)

	// request data
	data := req.Param{
		"chat_id":    chatId,
		"text":       text,
		"parse_mode": "MarkdownV2",
	}

	// perform request
	_, err := performIMRequest(webhookUrl, data)
	if err != nil {
		return err
	}
	return nil
}

func sendIMDiscord(ch *models.NotificationChannelV2, title string, content string) error {
	data := req.Param{
		"embeds": []req.Param{
			{
				"title":       title,
				"description": content,
			},
		},
	}
	_, err := performIMRequest(ch.WebhookUrl, data)
	if err != nil {
		return err
	}
	return nil
}

func sendIMMSTeams(ch *models.NotificationChannelV2, title string, content string) error {
	data := req.Param{
		"type": "message",
		"attachments": []req.Param{{
			"contentType": "application/vnd.microsoft.card.adaptive",
			"contentUrl":  nil,
			"content": req.Param{
				"$schema": "https://adaptivecards.io/schemas/adaptive-card.json",
				"type":    "AdaptiveCard",
				"version": "1.2",
				"body": []req.Param{
					{
						"type": "TextBlock",
						"text": fmt.Sprintf("**%s**", title),
						"size": "Large",
					},
					{
						"type": "TextBlock",
						"text": content,
					},
				},
			},
		}},
	}
	_, err := performIMRequest(ch.WebhookUrl, data)
	if err != nil {
		return err
	}
	return nil
}
