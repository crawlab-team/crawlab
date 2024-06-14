package test

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/delegate"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/core/result"
	"github.com/crawlab-team/crawlab/db/mongo"
	"go.uber.org/dig"
	"testing"
)

func init() {
	T = NewTest()
}

var T *Test

type Test struct {
	// dependencies
	modelSvc  service.ModelService
	resultSvc interfaces.ResultService

	// test data
	TestColName string
	TestCol     *mongo.Col
	TestDc      *models.DataCollection
}

func (t *Test) Setup(t2 *testing.T) {
	t2.Cleanup(t.Cleanup)
}

func (t *Test) Cleanup() {
	_ = t.modelSvc.DropAll()
}

func NewTest() *Test {
	var err error

	// test
	t := &Test{
		TestColName: "test_results",
	}

	// dependency injection
	c := dig.New()
	if err := c.Provide(service.NewService); err != nil {
		panic(err)
	}
	if err := c.Invoke(func(
		modelSvc service.ModelService,
	) {
		t.modelSvc = modelSvc
	}); err != nil {
		panic(err)
	}

	// data collection
	t.TestDc = &models.DataCollection{
		Name: t.TestColName,
	}
	if err := delegate.NewModelDelegate(t.TestDc).Add(); err != nil {
		panic(err)
	}
	t.TestCol = mongo.GetMongoCol(t.TestColName)

	// result service
	t.resultSvc, err = result.GetResultService(t.TestDc.GetId())
	if err != nil {
		panic(err)
	}

	return t
}
