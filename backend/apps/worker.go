package apps

import (
	"github.com/crawlab-team/crawlab-core/config"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/node/service"
	"go.uber.org/dig"
)

type Worker struct {
	// settings
	grpcAddress interfaces.Address

	// dependencies
	interfaces.WithConfigPath
	workerSvc interfaces.NodeWorkerService

	// internals
	quit chan int
}

func (app *Worker) SetGrpcAddress(address interfaces.Address) {
	app.grpcAddress = address
}

func (app *Worker) Init() {
}

func (app *Worker) Start() {
	go app.workerSvc.Start()
}

func (app *Worker) Wait() {
	<-app.quit
}

func (app *Worker) Stop() {
	app.quit <- 1
}

func NewWorker(opts ...WorkerOption) (app *Worker) {
	// worker
	w := &Worker{
		WithConfigPath: config.NewConfigPathService(),
		quit:           make(chan int, 1),
	}

	// apply options
	for _, opt := range opts {
		opt(w)
	}

	// service options
	var svcOpts []service.Option
	if w.grpcAddress != nil {
		svcOpts = append(svcOpts, service.WithAddress(w.grpcAddress))
	}

	// dependency injection
	c := dig.New()
	if err := c.Provide(service.ProvideWorkerService(w.GetConfigPath(), svcOpts...)); err != nil {
		panic(err)
	}
	if err := c.Invoke(func(workerSvc interfaces.NodeWorkerService) {
		w.workerSvc = workerSvc
	}); err != nil {
		panic(err)
	}

	return w
}
