package config

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
)

type Option func(svc interfaces.NodeConfigService)

func WithConfigPath(path string) Option {
	return func(svc interfaces.NodeConfigService) {
		svc.SetConfigPath(path)
	}
}
