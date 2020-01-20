package dingrobot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Roboter is the interface implemented by Robot that can send multiple types of messages.
type Roboter interface {
	SendText(content string, atMobiles []string, isAtAll bool) error
	SendLink(title, text, messageURL, picURL string) error
	SendMarkdown(title, text string, atMobiles []string, isAtAll bool) error
	SendActionCard(title, text, singleTitle, singleURL, btnOrientation, hideAvatar string) error
}

// Robot represents a dingtalk custom robot that can send messages to groups.
type Robot struct {
	Webhook string
}

// NewRobot returns a roboter that can send messages.
func NewRobot(webhook string) Roboter {
	return Robot{Webhook: webhook}
}

// SendText send a text type message.
func (r Robot) SendText(content string, atMobiles []string, isAtAll bool) error {
	return r.send(&textMessage{
		MsgType: msgTypeText,
		Text: textParams{
			Content: content,
		},
		At: atParams{
			AtMobiles: atMobiles,
			IsAtAll:   isAtAll,
		},
	})
}

// SendLink send a link type message.
func (r Robot) SendLink(title, text, messageURL, picURL string) error {
	return r.send(&linkMessage{
		MsgType: msgTypeLink,
		Link: linkParams{
			Title:      title,
			Text:       text,
			MessageURL: messageURL,
			PicURL:     picURL,
		},
	})
}

// SendMarkdown send a markdown type message.
func (r Robot) SendMarkdown(title, text string, atMobiles []string, isAtAll bool) error {
	return r.send(&markdownMessage{
		MsgType: msgTypeMarkdown,
		Markdown: markdownParams{
			Title: title,
			Text:  text,
		},
		At: atParams{
			AtMobiles: atMobiles,
			IsAtAll:   isAtAll,
		},
	})
}

// SendActionCard send a action card type message.
func (r Robot) SendActionCard(title, text, singleTitle, singleURL, btnOrientation, hideAvatar string) error {
	return r.send(&actionCardMessage{
		MsgType: msgTypeActionCard,
		ActionCard: actionCardParams{
			Title:          title,
			Text:           text,
			SingleTitle:    singleTitle,
			SingleURL:      singleURL,
			BtnOrientation: btnOrientation,
			HideAvatar:     hideAvatar,
		},
	})
}

type dingResponse struct {
	Errcode int
	Errmsg  string
}

func (r Robot) send(msg interface{}) error {
	m, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	resp, err := http.Post(r.Webhook, "application/json", bytes.NewReader(m))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var dr dingResponse
	err = json.Unmarshal(data, &dr)
	if err != nil {
		return err
	}
	if dr.Errcode != 0 {
		return fmt.Errorf("dingrobot send failed: %v", dr.Errmsg)
	}

	return nil
}
