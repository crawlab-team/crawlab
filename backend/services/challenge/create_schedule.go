package challenge

import (
	"crawlab/model"
	"github.com/globalsign/mgo/bson"
)

type CreateScheduleService struct {
	UserId bson.ObjectId
}

func (s *CreateScheduleService) Check() (bool, error) {
	query := bson.M{
		"user_id": s.UserId,
	}
	list, err := model.GetScheduleList(query)
	if err != nil {
		return false, err
	}
	return len(list) > 0, nil
}
