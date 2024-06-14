package interfaces

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type TaskHandlerService interface {
	TaskBaseService
	// Run task and execute locally
	Run(taskId primitive.ObjectID) (err error)
	// Cancel task locally
	Cancel(taskId primitive.ObjectID) (err error)
	// Fetch tasks and run
	Fetch()
	// ReportStatus periodically report handler status to master
	ReportStatus()
	// Reset reset internals to default
	Reset()
	// IsSyncLocked whether the given task is locked for files sync
	IsSyncLocked(path string) (ok bool)
	// LockSync lock files sync for given task
	LockSync(path string)
	// UnlockSync unlock files sync for given task
	UnlockSync(path string)
	// GetExitWatchDuration get max runners
	GetExitWatchDuration() (duration time.Duration)
	// SetExitWatchDuration set max runners
	SetExitWatchDuration(duration time.Duration)
	// GetFetchInterval get report interval
	GetFetchInterval() (interval time.Duration)
	// SetFetchInterval set report interval
	SetFetchInterval(interval time.Duration)
	// GetReportInterval get report interval
	GetReportInterval() (interval time.Duration)
	// SetReportInterval set report interval
	SetReportInterval(interval time.Duration)
	// GetCancelTimeout get report interval
	GetCancelTimeout() (timeout time.Duration)
	// SetCancelTimeout set report interval
	SetCancelTimeout(timeout time.Duration)
	// GetModelService get model service
	GetModelService() (modelSvc GrpcClientModelService)
	// GetModelSpiderService get model spider service
	GetModelSpiderService() (modelSpiderSvc GrpcClientModelSpiderService)
	// GetModelTaskService get model task service
	GetModelTaskService() (modelTaskSvc GrpcClientModelTaskService)
	// GetModelTaskStatService get model task stat service
	GetModelTaskStatService() (modelTaskStatSvc GrpcClientModelTaskStatService)
	// GetModelEnvironmentService get model environment service
	GetModelEnvironmentService() (modelEnvironmentSvc GrpcClientModelEnvironmentService)
	// GetNodeConfigService get node config service
	GetNodeConfigService() (cfgSvc NodeConfigService)
	// GetCurrentNode get node of the handler
	GetCurrentNode() (n Node, err error)
	// GetTaskById get task by id
	GetTaskById(id primitive.ObjectID) (t Task, err error)
	// GetSpiderById get task by id
	GetSpiderById(id primitive.ObjectID) (t Spider, err error)
}
