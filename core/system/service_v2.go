package system

import (
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/models/service"
	"go.mongodb.org/mongo-driver/bson"
)

type ServiceV2 struct {
}

func (svc *ServiceV2) Init() (err error) {
	// initialize data
	if err := svc.initData(); err != nil {
		return err
	}

	return nil
}

func (svc *ServiceV2) initData() (err error) {
	total, err := service.NewModelServiceV2[models.SettingV2]().Count(bson.M{
		"key": "site_title",
	})
	if err != nil {
		return err
	}
	if total > 0 {
		return nil
	}

	// data to initialize
	settings := []models.SettingV2{
		{
			Key: "site_title",
			Value: bson.M{
				"customize_site_title": false,
				"site_title":           "",
			},
		},
	}
	_, err = service.NewModelServiceV2[models.SettingV2]().InsertMany(settings)
	if err != nil {
		return err
	}
	return nil
}

func NewServiceV2() *ServiceV2 {
	// service
	svc := &ServiceV2{}

	if err := svc.Init(); err != nil {
		panic(err)
	}

	return svc
}

var _serviceV2 *ServiceV2

func GetServiceV2() *ServiceV2 {
	if _serviceV2 == nil {
		_serviceV2 = NewServiceV2()
	}
	return _serviceV2
}
