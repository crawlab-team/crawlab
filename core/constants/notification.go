package constants

const (
	NotificationTriggerPatternTask = "^task"
	NotificationTriggerPatternNode = "^node"
)

const (
	NotificationTriggerTaskFinish       = "task_finish"
	NotificationTriggerTaskError        = "task_error"
	NotificationTriggerTaskEmptyResults = "task_empty_results"
	NotificationTriggerNodeStatusChange = "node_status_change"
	NotificationTriggerNodeOnline       = "node_online"
	NotificationTriggerNodeOffline      = "node_offline"
	NotificationTriggerAlert            = "alert"
)

const (
	NotificationTemplateModeRichText = "rich-text"
	NotificationTemplateModeMarkdown = "markdown"
)
