package challenge

import (
	"crawlab/model"
	"github.com/globalsign/mgo/bson"
)

type Login30dService struct {
	UserId bson.ObjectId
}

func (s *Login30dService) Check() (bool, error) {
	days, err := model.GetVisitDays(s.UserId)
	if err != nil {
		return false, err
	}
	return days >= 30, nil
}
