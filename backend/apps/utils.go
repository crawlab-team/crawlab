package apps

import (
	"fmt"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab-core/utils"
	"github.com/crawlab-team/go-trace"
)

func Start(app App) {
	start(app)
}

func start(app App) {
	app.Init()
	go app.Start()
	app.Wait()
	app.Stop()
}

func DefaultWait() {
	utils.DefaultWait()
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

func initApp(name string, app App) {
	_ = initModule(name, func() error {
		app.Init()
		return nil
	})
}
