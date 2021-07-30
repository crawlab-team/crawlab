package apps

import (
	"github.com/crawlab-team/crawlab-core/config"
	"github.com/crawlab-team/crawlab-core/controllers"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/node/service"
	"go.uber.org/dig"
)

type Master struct {
	// settings
	runOnMaster bool
	grpcAddress interfaces.Address

	// dependencies
	interfaces.WithConfigPath
	masterSvc interfaces.NodeMasterService

	// modules
	api *Api

	// internals
	quit chan int
}

func (app *Master) SetGrpcAddress(address interfaces.Address) {
	app.grpcAddress = address
}

func (app *Master) SetRunOnMaster(ok bool) {
	app.runOnMaster = ok
}

func (app *Master) Init() {
	// initialize controllers
	if err := controllers.InitControllers(); err != nil {
		panic(err)
	}
}

func (app *Master) Start() {
	go app.api.Start()
	go app.masterSvc.Start()
}

func (app *Master) Wait() {
	<-app.quit
}

func (app *Master) Stop() {
	app.api.Stop()
	app.quit <- 1
}

func NewMaster(opts ...MasterOption) (app MasterApp) {
	// master
	m := &Master{
		WithConfigPath: config.NewConfigPathService(),
		api:            NewApi(),
		quit:           make(chan int, 1),
	}

	// apply options
	for _, opt := range opts {
		opt(m)
	}

	// service options
	var svcOpts []service.Option
	if m.grpcAddress != nil {
		svcOpts = append(svcOpts, service.WithAddress(m.grpcAddress))
	}

	// dependency injection
	c := dig.New()
	if err := c.Provide(service.ProvideMasterService(m.GetConfigPath(), svcOpts...)); err != nil {
		panic(err)
	}
	if err := c.Invoke(func(masterSvc interfaces.NodeMasterService) {
		m.masterSvc = masterSvc
	}); err != nil {
		panic(err)
	}

	return m
}
