package interfaces

type GrpcModelDelegateMessage interface {
	GetModelId() ModelId
	GetMethod() ModelDelegateMethod
	GetData() []byte
	ToBytes() (data []byte)
}
