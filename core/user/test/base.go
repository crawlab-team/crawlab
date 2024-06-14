package test

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/core/user"
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
	modelSvc service.ModelService
	userSvc  interfaces.UserService

	// test data
	TestUsername    string
	TestPassword    string
	TestNewPassword string
}

func (t *Test) Setup(t2 *testing.T) {
	var err error
	t.userSvc, err = user.NewUserService()
	if err != nil {
		panic(err)
	}
	t2.Cleanup(t.Cleanup)
}

func (t *Test) Cleanup() {
	_ = t.modelSvc.GetBaseService(interfaces.ModelIdUser).DeleteList(nil)
}

func NewTest() (t *Test, err error) {
	// test
	t = &Test{
		TestUsername:    "test_username",
		TestPassword:    "test_password",
		TestNewPassword: "test_new_password",
	}

	// dependency injection
	c := dig.New()
	if err := c.Provide(service.GetService); err != nil {
		return nil, err
	}
	if err := c.Invoke(func(modelSvc service.ModelService) {
		t.modelSvc = modelSvc
	}); err != nil {
		return nil, err
	}

	return t, nil
}
