package interfaces

type Tag interface {
	Model
	GetName() string
	GetColor() string
	SetCol(string)
}
