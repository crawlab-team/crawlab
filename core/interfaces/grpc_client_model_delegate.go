package interfaces

type GrpcClientModelDelegate interface {
	ModelDelegate
	WithConfigPath
	Close() error
}
