package apps

import (
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab-core/grpc"
)

type Worker struct {
	handler *Handler
}

func (app *Worker) Init() {
	_ = initModule("grpc", grpc.InitGrpcServices)
}

func (app *Worker) Run() {
	log.Info("worker has started")
}

func (app *Worker) Stop() {
	log.Info("worker has stopped")
}

func NewWorker() *Worker {
	return &Worker{
		handler: NewHandler(),
	}
}
