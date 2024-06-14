package client

import (
	config2 "github.com/crawlab-team/crawlab/core/config"
	"github.com/crawlab-team/crawlab/core/container"
	"github.com/crawlab-team/crawlab/core/interfaces"
)

type ServiceDelegate struct {
	// settings
	cfgPath string

	// internals
	c interfaces.GrpcClient
}

func (d *ServiceDelegate) GetConfigPath() string {
	return d.cfgPath
}

func (d *ServiceDelegate) SetConfigPath(path string) {
	d.cfgPath = path
}

func (d *ServiceDelegate) NewBaseServiceDelegate(id interfaces.ModelId) (svc interfaces.GrpcClientModelBaseService, err error) {
	var opts []ModelBaseServiceDelegateOption
	opts = append(opts, WithBaseServiceModelId(id))
	if d.cfgPath != "" {
		opts = append(opts, WithBaseServiceConfigPath(d.cfgPath))
	}
	return NewBaseServiceDelegate(opts...)
}

func NewServiceDelegate() (svc2 interfaces.GrpcClientModelService, err error) {
	// service delegate
	svc := &ServiceDelegate{
		cfgPath: config2.GetConfigPath(),
	}

	// dependency injection
	if err := container.GetContainer().Invoke(func(client interfaces.GrpcClient) {
		svc.c = client
	}); err != nil {
		return nil, err
	}

	return svc, nil
}
