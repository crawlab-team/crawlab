package utils

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
	"sync"
)

var moduleInitializedMap = sync.Map{}

func InitModule(id interfaces.ModuleId, fn func() error) (err error) {
	res, ok := moduleInitializedMap.Load(id)
	if ok {
		initialized, _ := res.(bool)
		if initialized {
			return nil
		}
	}

	if err := fn(); err != nil {
		return err
	}

	moduleInitializedMap.Store(id, true)

	return nil
}

func ForceInitModule(fn func() error) (err error) {
	return fn()
}
