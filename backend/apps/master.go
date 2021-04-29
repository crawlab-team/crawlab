package apps

type Master struct {
	api       *Api
	scheduler *Scheduler
	quit      chan int
}

func (app *Master) Init() {
	// api
	initApp("api", app.api)

	// scheduler
	initApp("scheduler", app.scheduler)
}

func (app *Master) Start() {
	go app.api.Start()
	go app.scheduler.Start()
}

func (app *Master) Wait() {
	<-app.quit
}

func (app *Master) Stop() {
	app.api.Stop()
	app.scheduler.Stop()

	app.quit <- 1
}

func NewMaster() *Master {
	return &Master{
		api:       NewApi(),
		scheduler: NewScheduler(),
		quit:      make(chan int, 1),
	}
}
