package apps

import (
	"github.com/apex/log"
)

type Scheduler struct {
}

func (app *Scheduler) Init() {
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
