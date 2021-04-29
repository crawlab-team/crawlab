package apps

type Master struct {
	api *Api
}

func (app *Master) Init() {
	panic("implement me")
}

func (app *Master) Run() {
	panic("implement me")
}

func (app *Master) Stop() {
	panic("implement me")
}

func NewMaster() *Master {
	return &Master{
		api: NewApi(),
	}
}
