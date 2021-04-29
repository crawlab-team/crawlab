package apps

import (
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab-core/grpc"
)

type Handler struct {
}

func (app *Handler) Init() {
	_ = initModule("grpc", grpc.InitGrpcServices)
}

func (app *Handler) Run() {
	log.Info("handler has started")
}

func (app *Handler) Stop() {
	log.Info("handler has stopped")
}

func NewHandler() *Handler {
	return &Handler{}
}
