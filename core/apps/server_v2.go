package apps

import (
	"fmt"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab/core/config"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/node/service"
	"github.com/crawlab-team/crawlab/core/utils"
	"github.com/spf13/viper"
	"net/http"
	_ "net/http/pprof"
)

type ServerV2 struct {
	// settings
	grpcAddress interfaces.Address

	// dependencies
	interfaces.WithConfigPath

	// modules
	nodeSvc interfaces.NodeService
	api     *ApiV2
	dck     *Docker

	// internals
	quit chan int
}

func (app *ServerV2) Init() {
	// log node info
	app.logNodeInfo()

	// pprof
	app.initPprof()
}

func (app *ServerV2) Start() {
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

func (app *ServerV2) Wait() {
	<-app.quit
}

func (app *ServerV2) Stop() {
	app.api.Stop()
	app.quit <- 1
}

func (app *ServerV2) GetApi() ApiApp {
	return app.api
}

func (app *ServerV2) GetNodeService() interfaces.NodeService {
	return app.nodeSvc
}

func (app *ServerV2) logNodeInfo() {
	log.Infof("current node type: %s", utils.GetNodeType())
	if utils.IsDocker() {
		log.Infof("running in docker container")
	}
}

func (app *ServerV2) initPprof() {
	if viper.GetBool("pprof") {
		go func() {
			fmt.Println(http.ListenAndServe("0.0.0.0:6060", nil))
		}()
	}
}

func NewServerV2() (app NodeApp) {
	// server
	svr := &ServerV2{
		WithConfigPath: config.NewConfigPathService(),
		quit:           make(chan int, 1),
	}

	// master modules
	if utils.IsMaster() {
		// api
		svr.api = GetApiV2()

		// docker
		if utils.IsDocker() {
			svr.dck = GetDocker(svr)
		}
	}

	// node service
	var err error
	if utils.IsMaster() {
		svr.nodeSvc, err = service.GetMasterServiceV2()
	} else {
		svr.nodeSvc, err = service.GetWorkerServiceV2()
	}
	if err != nil {
		panic(err)
	}

	return svr
}

var serverV2 NodeApp

func GetServerV2() NodeApp {
	if serverV2 != nil {
		return serverV2
	}
	serverV2 = NewServerV2()
	return serverV2
}
