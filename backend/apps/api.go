package apps

import (
	"context"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab-core/controllers"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/middlewares"
	"github.com/crawlab-team/crawlab-core/routes"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net"
	"net/http"
	"time"
)

type Api struct {
	// dependencies
	interfaces.WithConfigPath

	// internals
	app *gin.Engine
	srv *http.Server
}

func (app *Api) Init() {
	// initialize controllers
	_ = initModule("controllers", controllers.InitControllers)

	// initialize middlewares
	_ = app.initModuleWithApp("middlewares", middlewares.InitMiddlewares)

	// initialize routes
	_ = app.initModuleWithApp("routes", routes.InitRoutes)
}

func (app *Api) Start() {
	host := viper.GetString("server.host")
	port := viper.GetString("server.port")
	address := net.JoinHostPort(host, port)
	app.srv = &http.Server{
		Handler: app.app,
		Addr:    address,
	}
	if err := app.srv.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			log.Error("run server error:" + err.Error())
		} else {
			log.Info("server graceful down")
		}
	}
}

func (app *Api) Wait() {
	DefaultWait()
}

func (app *Api) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.srv.Shutdown(ctx); err != nil {
		log.Error("run server error:" + err.Error())
	}
}

func (app *Api) initModuleWithApp(name string, fn func(app *gin.Engine) error) (err error) {
	return initModule(name, func() error {
		return fn(app.app)
	})
}

func NewApi() *Api {
	api := &Api{
		app: gin.New(),
	}
	// api.Init()
	return api
}
