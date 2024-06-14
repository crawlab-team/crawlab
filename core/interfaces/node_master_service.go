package interfaces

import (
	"time"
)

type NodeMasterService interface {
	NodeService
	Monitor()
	SetMonitorInterval(duration time.Duration)
	Register() error
	StopOnError()
	GetServer() GrpcServer
}
