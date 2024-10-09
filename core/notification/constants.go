package notification

const (
	TypeMail = "mail"
	TypeIM   = "im"
)

const (
	ChannelMailProviderGmail   = "gmail"
	ChannelMailProviderOutlook = "outlook"
	ChannelMailProviderYahoo   = "yahoo"
	ChannelMailProviderICloud  = "icloud"
	ChannelMailProviderAol     = "aol"
	ChannelMailProviderZoho    = "zoho"
	ChannelMailProviderQQ      = "qq"
	ChannelMailProvider163     = "163"
	ChannelMailProviderExmail  = "exmail"

	ChannelIMProviderSlack      = "slack"       // https://api.slack.com/messaging/webhooks
	ChannelIMProviderTelegram   = "telegram"    // https://core.telegram.org/bots/api
	ChannelIMProviderDiscord    = "discord"     // https://support.discord.com/hc/en-us/articles/228383668-Intro-to-Webhooks
	ChannelIMProviderMSTeams    = "ms_teams"    // https://learn.microsoft.com/en-us/microsoftteams/platform/webhooks-and-connectors/how-to/add-incoming-webhook?tabs=newteams%2Cjavascript
	ChannelIMProviderWechatWork = "wechat_work" // https://developer.work.weixin.qq.com/document/path/91770
	ChannelIMProviderDingtalk   = "dingtalk"    // https://open.dingtalk.com/document/orgapp/custom-robot-access
	ChannelIMProviderLark       = "lark"        // https://www.larksuite.com/hc/en-US/articles/099698615114-use-webhook-triggers
)

const (
	StatusSending = "sending"
	StatusSuccess = "success"
	StatusError   = "error"
)
