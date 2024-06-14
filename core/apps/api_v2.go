package apps

import (
	"context"
	"errors"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab/core/controllers"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/middlewares"
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

type ApiV2 struct {
	// dependencies
	interfaces.WithConfigPath

	// internals
	app   *gin.Engine
	ln    net.Listener
	srv   *http.Server
	ready bool
}

func (app *ApiV2) Init() {
	// initialize middlewares
	_ = app.initModuleWithApp("middlewares", middlewares.InitMiddlewares)

	// initialize routes
	_ = app.initModuleWithApp("routes", controllers.InitRoutes)
}

func (app *ApiV2) Start() {
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

func (app *ApiV2) Wait() {
	DefaultWait()
}

func (app *ApiV2) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.srv.Shutdown(ctx); err != nil {
		log.Error("run server error:" + err.Error())
	}
}

func (app *ApiV2) GetGinEngine() *gin.Engine {
	return app.app
}

func (app *ApiV2) GetHttpServer() *http.Server {
	return app.srv
}

func (app *ApiV2) Ready() (ok bool) {
	return app.ready
}

func (app *ApiV2) initModuleWithApp(name string, fn func(app *gin.Engine) error) (err error) {
	return initModule(name, func() error {
		return fn(app.app)
	})
}

func NewApiV2() *ApiV2 {
	api := &ApiV2{
		app: gin.New(),
	}
	api.Init()
	return api
}

var apiV2 *ApiV2

func GetApiV2() *ApiV2 {
	if apiV2 != nil {
		return apiV2
	}
	apiV2 = NewApiV2()
	return apiV2
}
