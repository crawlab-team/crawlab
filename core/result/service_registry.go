package result

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
	"sync"
)

type ServiceRegistry struct {
	// internals
	services sync.Map
}

func (r *ServiceRegistry) Register(key string, fn interfaces.ResultServiceRegistryFn) {
	r.services.Store(key, fn)
}

func (r *ServiceRegistry) Unregister(key string) {
	r.services.Delete(key)
}

func (r *ServiceRegistry) Get(key string) (fn interfaces.ResultServiceRegistryFn) {
	res, ok := r.services.Load(key)
	if ok {
		fn, ok = res.(interfaces.ResultServiceRegistryFn)
		if !ok {
			return nil
		}
		return fn
	}
	return nil
}

func NewResultServiceRegistry() (r interfaces.ResultServiceRegistry) {
	r = &ServiceRegistry{
		services: sync.Map{},
	}
	return r
}

var _svc interfaces.ResultServiceRegistry

func GetResultServiceRegistry() (r interfaces.ResultServiceRegistry) {
	if _svc != nil {
		return _svc
	}
	_svc = NewResultServiceRegistry()
	return _svc
}
