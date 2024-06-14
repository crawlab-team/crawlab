package interfaces

import "time"

type NodeWorkerService interface {
	NodeService
	Register()
	Recv()
	ReportStatus()
	SetHeartbeatInterval(duration time.Duration)
}
