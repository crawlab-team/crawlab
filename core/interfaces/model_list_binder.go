package interfaces

type ModelListBinder interface {
	Bind() (l List, err error)
	Process(d interface{}) (l List, err error)
}
