package interfaces

type Color interface {
	Entity
	GetHex() string
	GetName() string
}
