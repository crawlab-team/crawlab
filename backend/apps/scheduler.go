package apps

import (
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab-core/config"
	"github.com/crawlab-team/crawlab-core/grpc"
)

type Scheduler struct {
}

func (app *Scheduler) Init() {
	// config
	_ = initModule("config", config.InitConfig)

	// grpc
	_ = initModule("grpc", grpc.InitGrpcServices)
}

func (app *Scheduler) Start() {
	log.Info("scheduler has started")
}

func (app *Scheduler) Wait() {
	DefaultWait()
}

func (app *Scheduler) Stop() {
	log.Info("scheduler has stopped")
}

func NewScheduler() *Scheduler {
	return &Scheduler{}
}
