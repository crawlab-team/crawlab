package result

import (
	"fmt"
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/trace"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
)

func NewResultService(registryKey string, s *models.Spider) (svc2 interfaces.ResultService, err error) {
	// result service function
	var fn interfaces.ResultServiceRegistryFn

	if registryKey == "" {
		// default
		fn = NewResultServiceMongo
	} else {
		// from registry
		reg := GetResultServiceRegistry()
		fn = reg.Get(registryKey)
		if fn == nil {
			return nil, errors.NewResultError(fmt.Sprintf("%s is not implemented", registryKey))
		}
	}

	// generate result service
	svc, err := fn(s.ColId, s.DataSourceId)
	if err != nil {
		return nil, trace.TraceError(err)
	}

	return svc, nil
}

var store = sync.Map{}

func GetResultService(spiderId primitive.ObjectID, opts ...Option) (svc2 interfaces.ResultService, err error) {
	// model service
	modelSvc, err := service.GetService()
	if err != nil {
		return nil, trace.TraceError(err)
	}

	// spider
	s, err := modelSvc.GetSpiderById(spiderId)
	if err != nil {
		return nil, trace.TraceError(err)
	}

	// apply options
	_opts := &Options{}
	for _, opt := range opts {
		opt(_opts)
	}

	// store key
	storeKey := s.ColId.Hex() + ":" + s.DataSourceId.Hex()

	// attempt to load result service from store
	res, _ := store.Load(storeKey)
	if res != nil {
		svc, ok := res.(interfaces.ResultService)
		if ok {
			return svc, nil
		}
	}

	// registry key
	var registryKey string
	ds, _ := modelSvc.GetDataSourceById(s.DataSourceId)
	if ds != nil {
		registryKey = ds.Type
	}

	// create a new result service if not exists
	svc, err := NewResultService(registryKey, s)
	if err != nil {
		return nil, err
	}

	// save into store
	store.Store(storeKey, svc)

	return svc, nil
}
