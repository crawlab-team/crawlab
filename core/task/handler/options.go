package handler

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
	"time"
)

type Option func(svc interfaces.TaskHandlerService)

func WithConfigPath(path string) Option {
	return func(svc interfaces.TaskHandlerService) {
		svc.SetConfigPath(path)
	}
}

func WithExitWatchDuration(duration time.Duration) Option {
	return func(svc interfaces.TaskHandlerService) {
		svc.SetExitWatchDuration(duration)
	}
}

func WithReportInterval(interval time.Duration) Option {
	return func(svc interfaces.TaskHandlerService) {
		svc.SetReportInterval(interval)
	}
}

func WithCancelTimeout(timeout time.Duration) Option {
	return func(svc interfaces.TaskHandlerService) {
		svc.SetCancelTimeout(timeout)
	}
}

type RunnerOption func(r interfaces.TaskRunner)

func WithSubscribeTimeout(timeout time.Duration) RunnerOption {
	return func(r interfaces.TaskRunner) {
		r.SetSubscribeTimeout(timeout)
	}
}
