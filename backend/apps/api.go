package apps

import (
	"context"
	"fmt"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab-core/config"
	"github.com/crawlab-team/crawlab-core/middlewares"
	"github.com/crawlab-team/crawlab-core/routes"
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/crawlab-team/crawlab-db/redis"
	"github.com/crawlab-team/go-trace"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Api struct {
	app *gin.Engine
}

func (svc *Api) Init() {
	// initialize config
	_ = svc.initService("config", config.InitConfig)

	// initialize mongo
	_ = svc.initService("mongo", mongo.InitMongo)

	// initialize redis
	_ = svc.initService("redis", redis.InitRedis)

	// initialize middlewares
	_ = svc.initServiceWithApp("middlewares", middlewares.InitMiddlewares)

	// initialize routes
	_ = svc.initServiceWithApp("routes", routes.InitRoutes)
}

func (svc *Api) Run() {
	host := viper.GetString("server.host")
	port := viper.GetString("server.port")
	address := net.JoinHostPort(host, port)
	srv := &http.Server{
		Handler: svc.app,
		Addr:    address,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Error("run server error:" + err.Error())
			} else {
				log.Info("server graceful down")
			}
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx2, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx2); err != nil {
		log.Error("run server error:" + err.Error())
	}
}

func (svc *Api) initServiceWithApp(name string, fn func(app *gin.Engine) error) (err error) {
	return svc.initService(name, func() error {
		return fn(svc.app)
	})
}

func (svc *Api) initService(name string, fn func() error) (err error) {
	if err := fn(); err != nil {
		log.Error(fmt.Sprintf("init %s error: %s", name, err.Error()))
		_ = trace.TraceError(err)
		panic(err)
	}
	log.Info(fmt.Sprintf("initialized %s successfully", name))
	return nil
}

func NewApi() *Api {
	app := gin.New()
	return &Api{
		app: app,
	}
}
