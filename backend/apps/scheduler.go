package apps

import (
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab-core/grpc"
)

type Scheduler struct {
}

func (app *Scheduler) Init() {
	_ = initModule("grpc", grpc.InitGrpcServices)
}

func (app *Scheduler) Run() {
	log.Info("scheduler has started")
}

func (app *Scheduler) Stop() {
	log.Info("scheduler has stopped")
}

func NewScheduler() *Scheduler {
	return &Scheduler{}
}
