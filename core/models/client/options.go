package client

import "github.com/crawlab-team/crawlab/core/interfaces"

type ModelDelegateOption func(delegate interfaces.GrpcClientModelDelegate)

func WithDelegateConfigPath(path string) ModelDelegateOption {
	return func(d interfaces.GrpcClientModelDelegate) {
		d.SetConfigPath(path)
	}
}

type ModelServiceDelegateOption func(delegate interfaces.GrpcClientModelService)

func WithServiceConfigPath(path string) ModelServiceDelegateOption {
	return func(d interfaces.GrpcClientModelService) {
		d.SetConfigPath(path)
	}
}

type ModelBaseServiceDelegateOption func(delegate interfaces.GrpcClientModelBaseService)

func WithBaseServiceModelId(id interfaces.ModelId) ModelBaseServiceDelegateOption {
	return func(d interfaces.GrpcClientModelBaseService) {
		d.SetModelId(id)
	}
}

func WithBaseServiceConfigPath(path string) ModelBaseServiceDelegateOption {
	return func(d interfaces.GrpcClientModelBaseService) {
		d.SetConfigPath(path)
	}
}
