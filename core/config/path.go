package config

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
)

type PathService struct {
	cfgPath string
}

func (svc *PathService) GetConfigPath() (path string) {
	return svc.cfgPath
}

func (svc *PathService) SetConfigPath(path string) {
	svc.cfgPath = path
}

func NewConfigPathService() (svc interfaces.WithConfigPath) {
	svc = &PathService{}
	svc.SetConfigPath(GetConfigPath())
	return svc
}
