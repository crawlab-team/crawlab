package interfaces

type NodeService interface {
	Module
	WithConfigPath
	GetConfigService() NodeConfigService
}
