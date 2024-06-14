package interfaces

type ColorService interface {
	Injectable
	GetByName(name string) (res Color, err error)
	GetRandom() (res Color, err error)
}
