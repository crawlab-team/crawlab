package interfaces

type WithModelId interface {
	GetModelId() (id ModelId)
	SetModelId(id ModelId)
}
