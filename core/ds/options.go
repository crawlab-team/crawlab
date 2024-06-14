package ds

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
	"time"
)

type DataSourceServiceOption func(svc interfaces.DataSourceService)

func WithMonitorInterval(duration time.Duration) DataSourceServiceOption {
	return func(svc interfaces.DataSourceService) {
		svc.SetMonitorInterval(duration)
	}
}
