package constants

const (
	// 调度中
	StatusPending string = "pending"
	// 运行中
	StatusRunning string = "running"
	// 已完成
	StatusFinished string = "finished"
	// 错误
	StatusError string = "error"
	// 取消
	StatusCancelled string = "cancelled"
	// 节点重启导致的异常终止
	StatusAbnormal string = "abnormal"
)

const (
	TaskFinish string = "finish"
	TaskCancel string = "cancel"
)

const (
	RunTypeAllNodes      string = "all-nodes"
	RunTypeRandom        string = "random"
	RunTypeSelectedNodes string = "selected-nodes"
)

const (
	TaskTypeSpider string = "spider"
	TaskTypeSystem string = "system"
)
