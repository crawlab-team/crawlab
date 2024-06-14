package test

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/delegate"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/core/schedule"
	"go.uber.org/dig"
	"testing"
)

func init() {
	var err error
	T, err = NewTest()
	if err != nil {
		panic(err)
	}
}

var T *Test

type Test struct {
	// dependencies
	modelSvc    service.ModelService
	scheduleSvc interfaces.ScheduleService

	// test data
	TestSchedule interfaces.Schedule
	TestSpider   interfaces.Spider
	ScriptName   string
	Script       string
}

func (t *Test) Setup(t2 *testing.T) {
	t.scheduleSvc.Start()
	t2.Cleanup(t.Cleanup)
}

func (t *Test) Cleanup() {
	t.scheduleSvc.Stop()
	_ = t.modelSvc.GetBaseService(interfaces.ModelIdTask).Delete(nil)
}

func NewTest() (t *Test, err error) {
	// test
	t = &Test{
		TestSpider: &models.Spider{
			Name: "test_spider",
			Cmd:  "go run main.go",
		},
		ScriptName: "main.go",
		Script: `package main
import "fmt"
func main() {
  fmt.Println("it works")
}`,
	}

	// dependency injection
	c := dig.New()
	if err := c.Provide(service.GetService); err != nil {
		return nil, err
	}
	if err := c.Provide(schedule.NewScheduleService); err != nil {
		return nil, err
	}
	if err := c.Invoke(func(modelSvc service.ModelService, scheduleSvc interfaces.ScheduleService) {
		t.modelSvc = modelSvc
		t.scheduleSvc = scheduleSvc
	}); err != nil {
		return nil, err
	}

	// add spider to db
	if err := delegate.NewModelDelegate(t.TestSpider).Add(); err != nil {
		return nil, err
	}

	// test schedule
	t.TestSchedule = &models.Schedule{
		Name:     "test_schedule",
		SpiderId: t.TestSpider.GetId(),
		Cron:     "* * * * *",
	}
	if err := delegate.NewModelDelegate(t.TestSchedule).Add(); err != nil {
		return nil, err
	}

	return t, nil
}
