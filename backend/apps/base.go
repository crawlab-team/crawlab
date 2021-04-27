package apps

import (
	"fmt"
	"github.com/apex/log"
	"github.com/crawlab-team/go-trace"
	"os"
	"os/signal"
	"syscall"
)

type App interface {
	init()
	run()
	stop()
}

type BaseApp struct {
}

func (app *BaseApp) init() {
	panic("implement me")
}

func (app *BaseApp) run() {
	panic("implement me")
}

func (app *BaseApp) stop() {
	panic("implement me")
}

func (app *BaseApp) start() {
	app.init()
	go app.run()
	app.waitForStop()
	app.stop()
}

func (app *BaseApp) Start() {
	app.start()
}

func (app *BaseApp) waitForStop() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

func (app *BaseApp) initModule(name string, fn func() error) (err error) {
	if err := fn(); err != nil {
		log.Error(fmt.Sprintf("init %s error: %s", name, err.Error()))
		_ = trace.TraceError(err)
		panic(err)
	}
	log.Info(fmt.Sprintf("initialized %s successfully", name))
	return nil
}
