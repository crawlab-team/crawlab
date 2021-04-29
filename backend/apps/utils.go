package apps

import (
	"fmt"
	"github.com/apex/log"
	"github.com/crawlab-team/go-trace"
	"os"
	"os/signal"
	"syscall"
)

func Start(app App) {
	start(app)
}

func start(app App) {
	app.Init()
	go app.Run()
	waitForStop()
	app.Stop()
}

func waitForStop() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

func initModule(name string, fn func() error) (err error) {
	if err := fn(); err != nil {
		log.Error(fmt.Sprintf("init %s error: %s", name, err.Error()))
		_ = trace.TraceError(err)
		panic(err)
	}
	log.Info(fmt.Sprintf("initialized %s successfully", name))
	return nil
}
