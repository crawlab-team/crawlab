package entity

type NodeMessage struct {
	// 通信类别
	Type string `json:"type"`

	// 任务相关
	TaskId string `json:"task_id"` // 任务ID

	// 节点相关
	NodeId string `json:"node_id"` // 节点ID

	// 日志相关
	LogPath string `json:"log_path"` // 日志路径
	Log     string `json:"log"`      // 日志

	// 系统信息
	SysInfo SystemInfo `json:"sys_info"`

	// 爬虫相关
	SpiderId string `json:"spider_id"` //爬虫ID

	// 语言相关
	Lang Lang `json:"lang"`

	// 错误相关
	Error string `json:"error"`
}
