package interfaces

type ModuleId int

type Module interface {
	Init() error
	Start()
	Wait()
	Stop()
}
