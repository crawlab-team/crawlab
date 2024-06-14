package interfaces

type NodeService interface {
	Module
	WithConfigPath
	WithAddress
	GetConfigService() NodeConfigService
}
