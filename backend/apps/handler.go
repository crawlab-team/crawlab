package apps

import (
	"github.com/crawlab-team/crawlab-core/services"
)

type Handler struct {
	BaseApp
}

func (app *Handler) Init() {
	_ = app.initModule("task-service", services.InitTaskService)
}

func (app *Handler) Run() {
	panic("implement me")
}

func NewHandler() *Handler {
	return &Handler{}
}
