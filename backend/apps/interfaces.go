package apps

import "github.com/crawlab-team/crawlab-core/interfaces"

type App interface {
	Init()
	Start()
	Wait()
	Stop()
}

type MasterApp interface {
	App
	interfaces.WithConfigPath
	SetRunOnMaster(ok bool)
}
