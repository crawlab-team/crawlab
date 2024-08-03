package system

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/models/service"
	mongo2 "github.com/crawlab-team/crawlab/db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	col      *mongo2.Col
	modelSvc service.ModelService
}

func (svc *Service) Init() (err error) {
	// initialize data
	if err := svc.initData(); err != nil {
		return err
	}

	return nil
}

func (svc *Service) initData() (err error) {
	total, err := svc.col.Count(bson.M{
		"key": "customize",
	})
	if err != nil {
		return err
	}
	if total > 0 {
		return nil
	}

	// data to initialize
	settings := []models.Setting{
		{
			Id:  primitive.NewObjectID(),
			Key: "customize",
			Value: bson.M{
				"show_custom_title":     false,
				"custom_title":          "",
				"show_custom_logo":      false,
				"custom_logo":           "",
				"hide_platform_version": false,
			},
		},
	}
	var data []interface{}
	for _, s := range settings {
		data = append(data, s)
	}
	_, err = svc.col.InsertMany(data)
	if err != nil {
		return err
	}
	return nil
}

func NewService() *Service {
	// service
	svc := &Service{
		col: mongo2.GetMongoCol(interfaces.ModelColNameSetting),
	}

	// model service
	modelSvc, err := service.GetService()
	if err != nil {
		panic(err)
	}
	svc.modelSvc = modelSvc

	if err := svc.Init(); err != nil {
		panic(err)
	}

	return svc
}

var _service *Service

func GetService() *Service {
	if _service == nil {
		_service = NewService()
	}
	return _service
}
