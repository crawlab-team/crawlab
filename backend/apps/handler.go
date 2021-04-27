package apps

import (
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab-core/services"
)

type Handler struct {
	BaseApp
}

func (app *Handler) init() {
	_ = app.initModule("task-service", services.InitTaskService)
}

func (app *Handler) run() {
	log.Info("handler has started")
}

func NewHandler() *Handler {
	return &Handler{}
}
