package interfaces

type GrpcClientModelService interface {
	WithConfigPath
	NewBaseServiceDelegate(id ModelId) (GrpcClientModelBaseService, error)
}
