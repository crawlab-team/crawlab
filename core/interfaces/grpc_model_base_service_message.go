package interfaces

type GrpcModelBaseServiceMessage interface {
	GetModelId() ModelId
	GetData() []byte
	ToBytes() (data []byte)
}
