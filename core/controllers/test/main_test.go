package test

import (
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/controllers"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/core/user"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func TestMain(m *testing.M) {
	// init user
	modelSvc, err := service.GetService()
	if err != nil {
		panic(err)
	}
	_, err = modelSvc.GetUserByUsername(constants.DefaultAdminUsername, nil)
	if err != nil {
		if err.Error() != mongo.ErrNoDocuments.Error() {
			panic(err)
		}
		userSvc, err := user.GetUserService()
		if err != nil {
			panic(err)
		}
		if err := userSvc.Create(&interfaces.UserCreateOptions{
			Username: constants.DefaultAdminUsername,
			Password: constants.DefaultAdminPassword,
			Role:     constants.RoleAdmin,
		}); err != nil {
			panic(err)
		}
	}

	if err := controllers.InitControllers(); err != nil {
		panic(err)
	}

	m.Run()

	T.Cleanup()
}
