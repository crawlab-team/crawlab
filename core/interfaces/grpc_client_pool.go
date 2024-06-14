package interfaces

type GrpcClientPool interface {
	WithConfigPath
	Init() error
	NewClient() error
	GetClient() (GrpcClient, error)
	SetSize(int)
}
