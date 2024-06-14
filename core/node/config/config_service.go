package config

import (
	"encoding/json"
	"github.com/crawlab-team/crawlab/core/config"
	"github.com/crawlab-team/crawlab/core/entity"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/utils"
	"github.com/crawlab-team/crawlab/trace"
	"os"
	"path"
)

type Service struct {
	cfg  *Config
	path string
}

func (svc *Service) Init() (err error) {
	// check config directory path
	configDirPath := path.Dir(svc.path)
	if !utils.Exists(configDirPath) {
		if err := os.MkdirAll(configDirPath, os.FileMode(0766)); err != nil {
			return trace.TraceError(err)
		}
	}

	if !utils.Exists(svc.path) {
		// not exists, set to default config
		// and create a config file for persistence
		svc.cfg = NewConfig(nil)
		data, err := json.Marshal(svc.cfg)
		if err != nil {
			return trace.TraceError(err)
		}
		if err := os.WriteFile(svc.path, data, os.FileMode(0766)); err != nil {
			return trace.TraceError(err)
		}
	} else {
		// exists, read and set to config
		data, err := os.ReadFile(svc.path)
		if err != nil {
			return trace.TraceError(err)
		}
		if err := json.Unmarshal(data, svc.cfg); err != nil {
			return trace.TraceError(err)
		}
	}

	return nil
}

func (svc *Service) Reload() (err error) {
	return svc.Init()
}

func (svc *Service) GetBasicNodeInfo() (res interfaces.Entity) {
	return &entity.NodeInfo{
		Key:        svc.GetNodeKey(),
		Name:       svc.GetNodeName(),
		IsMaster:   svc.IsMaster(),
		AuthKey:    svc.GetAuthKey(),
		MaxRunners: svc.GetMaxRunners(),
	}
}

func (svc *Service) GetNodeKey() (res string) {
	return svc.cfg.Key
}

func (svc *Service) GetNodeName() (res string) {
	return svc.cfg.Name
}

func (svc *Service) IsMaster() (res bool) {
	return svc.cfg.IsMaster
}

func (svc *Service) GetAuthKey() (res string) {
	return svc.cfg.AuthKey
}

func (svc *Service) GetMaxRunners() (res int) {
	return svc.cfg.MaxRunners
}

func (svc *Service) GetConfigPath() (path string) {
	return svc.path
}

func (svc *Service) SetConfigPath(path string) {
	svc.path = path
}

func NewNodeConfigService() (svc2 interfaces.NodeConfigService, err error) {
	// cfg
	cfg := NewConfig(nil)

	// config service
	svc := &Service{
		cfg: cfg,
	}

	// normalize config path
	cfgPath := config.GetConfigPath()
	svc.SetConfigPath(cfgPath)

	// init
	if err := svc.Init(); err != nil {
		return nil, err
	}

	return svc, nil
}

var _service interfaces.NodeConfigService

func GetNodeConfigService() interfaces.NodeConfigService {
	if _service != nil {
		return _service
	}

	var err error
	_service, err = NewNodeConfigService()
	if err != nil {
		panic(err)
	}

	return _service
}
