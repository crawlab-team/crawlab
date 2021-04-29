package apps

import (
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab-core/config"
	"github.com/crawlab-team/crawlab-core/grpc"
)

type Handler struct {
}

func (app *Handler) Init() {
	// config
	_ = initModule("config", config.InitConfig)

	// grpc
	_ = initModule("grpc", grpc.InitGrpcServices)
}

func (app *Handler) Start() {
}

func (app *Handler) Wait() {
	DefaultWait()
}

func (app *Handler) Stop() {
	log.Info("handler has stopped")
}

func NewHandler() *Handler {
	return &Handler{}
}
