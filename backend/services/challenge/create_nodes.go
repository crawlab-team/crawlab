package challenge

import (
	"crawlab/constants"
	"crawlab/model"
	"github.com/globalsign/mgo/bson"
)

type CreateNodesService struct {
	UserId bson.ObjectId
}

func (s *CreateNodesService) Check() (bool, error) {
	query := bson.M{
		"status": constants.StatusOnline,
	}
	list, err := model.GetScheduleList(query)
	if err != nil {
		return false, err
	}
	return len(list) >= 3, nil
}
