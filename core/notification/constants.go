package notification

const (
	TypeMail   = "mail"
	TypeMobile = "mobile"
)

const (
	ChannelMailProviderGmail          = "gmail"
	ChannelMailProviderOutlook        = "outlook"
	ChannelMailProviderYahoo          = "yahoo"
	ChannelMailProviderHotmail        = "hotmail"
	ChannelMailProviderAol            = "aol"
	ChannelMailProviderZoho           = "zoho"
	ChannelMailProviderYandex         = "yandex"
	ChannelMailProviderICloud         = "icloud"
	ChannelMailProviderQQMailProvider = "qq"
	ChannelMailProvider163            = "163"
	ChannelMailProvider126            = "126"
	ChannelMailProviderSina           = "sina"
	ChannelMailProviderSohu           = "sohu"
	ChannelMailProvider21CN           = "21cn"
	ChannelMailProviderTencent        = "tencent"
	ChannelMailProviderHuawei         = "huawei"
	ChannelMailProviderAliyun         = "aliyun"

	ChannelIMProviderSlack             = "slack"              // https://api.slack.com/messaging/webhooks
	ChannelIMProviderMSTeams           = "ms_teams"           // https://learn.microsoft.com/en-us/microsoftteams/platform/webhooks-and-connectors/how-to/add-incoming-webhook?tabs=newteams%2Cjavascript
	ChannelIMProviderTelegram          = "telegram"           // https://core.telegram.org/bots/api
	ChannelIMProviderDiscord           = "discord"            // https://support.discord.com/hc/en-us/articles/228383668-Intro-to-Webhooks
	ChannelIMProviderWhatsappBusiness  = "whatsapp_business"  // https://developers.facebook.com/docs/whatsapp/cloud-api/guides/send-messages
	ChannelIMProviderFacebookMessenger = "facebook_messenger" // https://developers.facebook.com/docs/messenger-platform/send-messages
	ChannelIMProviderWechatWork        = "wechat_work"        // https://developer.work.weixin.qq.com/document/path/91770
	ChannelIMProviderDingtalk          = "dingtalk"           // https://open.dingtalk.com/document/orgapp/custom-robot-access
	ChannelIMProviderLark              = "lark"               // https://www.larksuite.com/hc/en-US/articles/099698615114-use-webhook-triggers
)
