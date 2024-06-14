package apps

import (
	"context"
	"errors"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab/core/controllers"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/middlewares"
	"github.com/crawlab-team/crawlab/core/routes"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net"
	"net/http"
	"time"
)

func init() {
	// set gin mode
	if viper.GetString("gin.mode") == "" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(viper.GetString("gin.mode"))
	}
}

type Api struct {
	// dependencies
	interfaces.WithConfigPath

	// internals
	app   *gin.Engine
	ln    net.Listener
	srv   *http.Server
	ready bool
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
	// address
	host := viper.GetString("server.host")
	port := viper.GetString("server.port")
	address := net.JoinHostPort(host, port)

	// http server
	app.srv = &http.Server{
		Handler: app.app,
		Addr:    address,
	}

	// listen
	var err error
	app.ln, err = net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	app.ready = true

	// serve
	if err := http.Serve(app.ln, app.app); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
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

func (app *Api) GetGinEngine() *gin.Engine {
	return app.app
}

func (app *Api) GetHttpServer() *http.Server {
	return app.srv
}

func (app *Api) Ready() (ok bool) {
	return app.ready
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
	api.Init()
	return api
}

var api *Api

func GetApi() *Api {
	if api != nil {
		return api
	}
	api = NewApi()
	return api
}
