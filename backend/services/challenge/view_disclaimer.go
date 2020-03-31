package challenge

import (
	"crawlab/constants"
	"crawlab/model"
	"github.com/globalsign/mgo/bson"
)

type ViewDisclaimerService struct {
	UserId bson.ObjectId
}

func (s *ViewDisclaimerService) Check() (bool, error) {
	query := bson.M{
		"user_id": s.UserId,
		"type": constants.ActionTypeViewDisclaimer,
	}
	list, err := model.GetActionList(query, 0, 1, "-_id")
	if err != nil {
		return false, err
	}
	return len(list) > 0, nil
}
