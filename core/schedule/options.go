package schedule

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
	"time"
)

type Option func(svc interfaces.ScheduleService)

func WithConfigPath(path string) Option {
	return func(svc interfaces.ScheduleService) {
		svc.SetConfigPath(path)
	}
}

func WithLocation(loc *time.Location) Option {
	return func(svc interfaces.ScheduleService) {
		svc.SetLocation(loc)
	}
}

func WithDelayIfStillRunning() Option {
	return func(svc interfaces.ScheduleService) {
		svc.SetDelay(true)
	}
}

func WithSkipIfStillRunning() Option {
	return func(svc interfaces.ScheduleService) {
		svc.SetSkip(true)
	}
}

func WithUpdateInterval(interval time.Duration) Option {
	return func(svc interfaces.ScheduleService) {
	}
}
