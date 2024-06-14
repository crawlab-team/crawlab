package service

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/db/mongo"
)

type Option func(ModelService)

type BaseServiceOption func(svc interfaces.ModelBaseService)

func WithBaseServiceModelId(id interfaces.ModelId) BaseServiceOption {
	return func(svc interfaces.ModelBaseService) {
		svc.SetModelId(id)
	}
}

func WithBaseServiceCol(col *mongo.Col) BaseServiceOption {
	return func(svc interfaces.ModelBaseService) {
		_svc, ok := svc.(*BaseService)
		if ok {
			_svc.SetCol(col)
		}
	}
}
