package apps

type Worker struct {
	handler *Handler
	quit    chan int
}

func (app *Worker) Init() {
	initApp("handler", app.handler) // handler
}

func (app *Worker) Start() {
	go app.handler.Start()
}

func (app *Worker) Wait() {
	<-app.quit
}

func (app *Worker) Stop() {
	app.handler.Stop()

	app.quit <- 1
}

func NewWorker() *Worker {
	return &Worker{
		handler: NewHandler(),
		quit:    make(chan int, 1),
	}
}
