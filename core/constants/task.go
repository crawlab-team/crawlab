package constants

const (
	TaskStatusPending   = "pending"
	TaskStatusRunning   = "running"
	TaskStatusFinished  = "finished"
	TaskStatusError     = "error"
	TaskStatusCancelled = "cancelled"
	TaskStatusAbnormal  = "abnormal"
)

const (
	RunTypeAllNodes      = "all-nodes"
	RunTypeRandom        = "random"
	RunTypeSelectedNodes = "selected-nodes"
)

type TaskSignal int

const (
	TaskSignalFinish TaskSignal = iota
	TaskSignalCancel
	TaskSignalError
	TaskSignalLost
)

const (
	TaskKey   = "_tid"
	SpiderKey = "_sid"
)
