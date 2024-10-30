package interfaces

import "time"

type NodeWorkerService interface {
	NodeService
	Register()
	ReportStatus()
	SetHeartbeatInterval(duration time.Duration)
}
