package apps

type App interface {
	Init()
	Start()
	Wait()
	Stop()
}
