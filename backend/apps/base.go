package apps

import (
	"fmt"
	"github.com/apex/log"
	"github.com/crawlab-team/go-trace"
)

type App interface {
	Init()
	Run()
}

type BaseApp struct {
}

func (app *BaseApp) Init() {
	panic("implement me")
}

func (app *BaseApp) Run() {
	panic("implement me")
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
