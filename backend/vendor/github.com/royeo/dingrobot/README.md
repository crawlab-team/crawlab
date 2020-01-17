# Dingrobot

钉钉自定义机器人 Golang API.

支持的消息类型：
- 文本类型
- link 类型
- markdown 类型
- 整体跳转 ActionCard 类型

## Installation

Install:

```sh
go get -u github.com/royeo/dingrobot
```

Import:

```go
import "github.com/royeo/dingrobot"
```

## Quick start

发送文本类型的消息：

```go
func main() {
	// You should replace the webhook here with your own.
	webhook := "https://oapi.dingtalk.com/robot/send?access_token=xxx"
	robot := dingrobot.NewRobot(webhook)

	content := "我就是我,  @1825718XXXX 是不一样的烟火"
	atMobiles := []string{"1825718XXXX"}
	isAtAll := false

	err := robot.SendText(content, atMobiles, isAtAll)
	if err != nil {
		log.Fatal(err)
	}
}
```

发送 link 类型的消息：

```go
func main() {
	// You should replace the webhook here with your own.
	webhook := "https://oapi.dingtalk.com/robot/send?access_token=xxx"
	robot := dingrobot.NewRobot(webhook)

	title := "自定义机器人协议"
	text := "群机器人是钉钉群的高级扩展功能。群机器人可以将第三方服务的信息聚合到群聊中，实现自动化的信息同步。例如：通过聚合GitHub，GitLab等源码管理服务，实现源码更新同步；通过聚合Trello，JIRA等项目协调服务，实现项目信息同步。不仅如此，群机器人支持Webhook协议的自定义接入，支持更多可能性，例如：你可将运维报警提醒通过自定义机器人聚合到钉钉群。"
	messageUrl := "https://open-doc.dingtalk.com/docs/doc.htm?spm=a219a.7629140.0.0.Rqyvqo&treeId=257&articleId=105735&docType=1"
	picUrl := ""

	err := robot.SendLink(title, text, messageUrl, picUrl)
	if err != nil {
		log.Fatal(err)
	}
}
```

发送 markdown 类型的消息：

```go
func main() {
	// You should replace the webhook here with your own.
	webhook := "https://oapi.dingtalk.com/robot/send?access_token=xxx"
	robot := dingrobot.NewRobot(webhook)

	title := "杭州天气"
	text := "#### 杭州天气  \n > 9度，@1825718XXXX 西北风1级，空气良89，相对温度73%\n\n > ![screenshot](http://i01.lw.aliimg.com/media/lALPBbCc1ZhJGIvNAkzNBLA_1200_588.png)\n  > ###### 10点20分发布 [天气](http://www.thinkpage.cn/) "
	atMobiles := []string{"1825718XXXX"}
	isAtAll := false

	err := robot.SendMarkdown(title, text, atMobiles, isAtAll)
	if err != nil {
		log.Fatal(err)
	}
}
```

发送整体跳转 ActionCard 类型的消息：

```go
func main() {
	// You should replace the webhook here with your own.
	webhook := "https://oapi.dingtalk.com/robot/send?access_token=xxx"
	robot := dingrobot.NewRobot(webhook)

	title := "乔布斯 20 年前想打造一间苹果咖啡厅，而它正是 Apple Store 的前身"
	text := "![screenshot](@lADOpwk3K80C0M0FoA) \n #### 乔布斯 20 年前想打造的苹果咖啡厅 \n\n Apple Store 的设计正从原来满满的科技感走向生活化，而其生活化的走向其实可以追溯到 20 年前苹果一个建立咖啡馆的计划"
	singleTitle := "阅读全文"
	singleURL := "https://www.dingtalk.com/"
	btnOrientation := "0"
	hideAvatar := "0"

	err := robot.SendActionCard(title, text, singleTitle, singleURL, btnOrientation, hideAvatar)
	if err != nil {
		log.Fatal(err)
	}
}
```

## License

MIT Copyright (c) 2018 Royeo
