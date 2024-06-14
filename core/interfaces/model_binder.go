package interfaces

type ModelBinder interface {
	Bind() (res Model, err error)
	Process(d Model) (res Model, err error)
}
