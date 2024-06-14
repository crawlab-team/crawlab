package apps

import (
	"fmt"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab/core/config"
	"github.com/crawlab-team/crawlab/core/controllers"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/node/service"
	"github.com/crawlab-team/crawlab/core/utils"
	"github.com/spf13/viper"
	"net/http"
	_ "net/http/pprof"
)

func init() {
	injectModules()
}

type Server struct {
	// settings
	grpcAddress interfaces.Address

	// dependencies
	interfaces.WithConfigPath
	nodeSvc interfaces.NodeService

	// modules
	api *Api
	dck *Docker

	// internals
	quit chan int
}

func (app *Server) SetGrpcAddress(address interfaces.Address) {
	app.grpcAddress = address
}

func (app *Server) GetApi() (api ApiApp) {
	return app.api
}

func (app *Server) GetNodeService() (svc interfaces.NodeService) {
	return app.nodeSvc
}

func (app *Server) Init() {
	// log node info
	app.logNodeInfo()

	if utils.IsMaster() {

		// initialize controllers
		if err := controllers.InitControllers(); err != nil {
			panic(err)
		}
	}

	// pprof
	app.initPprof()
}

func (app *Server) Start() {
	if utils.IsMaster() {
		// start docker app
		if utils.IsDocker() {
			go app.dck.Start()
		}

		// start api
		go app.api.Start()
	}

	// start node service
	go app.nodeSvc.Start()
}

func (app *Server) Wait() {
	<-app.quit
}

func (app *Server) Stop() {
	app.api.Stop()
	app.quit <- 1
}

func (app *Server) logNodeInfo() {
	log.Infof("current node type: %s", utils.GetNodeType())
	if utils.IsDocker() {
		log.Infof("running in docker container")
	}
}

func (app *Server) initPprof() {
	if viper.GetBool("pprof") {
		go func() {
			fmt.Println(http.ListenAndServe("0.0.0.0:6060", nil))
		}()
	}
}

func NewServer() (app NodeApp) {
	// server
	svr := &Server{
		WithConfigPath: config.NewConfigPathService(),
		quit:           make(chan int, 1),
	}

	// service options
	var svcOpts []service.Option
	if svr.grpcAddress != nil {
		svcOpts = append(svcOpts, service.WithAddress(svr.grpcAddress))
	}

	// master modules
	if utils.IsMaster() {
		// api
		svr.api = GetApi()

		// docker
		if utils.IsDocker() {
			svr.dck = GetDocker(svr)
		}
	}

	// node service
	var err error
	if utils.IsMaster() {
		svr.nodeSvc, err = service.NewMasterService(svcOpts...)
	} else {
		svr.nodeSvc, err = service.NewWorkerService(svcOpts...)
	}
	if err != nil {
		panic(err)
	}

	return svr
}

var server NodeApp

func GetServer() NodeApp {
	if server != nil {
		return server
	}
	server = NewServer()
	return server
}
