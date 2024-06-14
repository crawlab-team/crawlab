package apps

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/gin-gonic/gin"
	"net/http"
)

type App interface {
	Init()
	Start()
	Wait()
	Stop()
}

type ApiApp interface {
	App
	GetGinEngine() (engine *gin.Engine)
	GetHttpServer() (svr *http.Server)
	Ready() (ok bool)
}

type NodeApp interface {
	App
	interfaces.WithConfigPath
}

type ServerApp interface {
	NodeApp
	GetApi() (api ApiApp)
	GetNodeService() (masterSvc interfaces.NodeService)
}

type DockerApp interface {
	App
	GetParent() (parent NodeApp)
	SetParent(parent NodeApp)
	Ready() (ok bool)
}
