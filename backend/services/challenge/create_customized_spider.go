package challenge

import (
	"crawlab/constants"
	"crawlab/model"
	"github.com/globalsign/mgo/bson"
)

type CreateCustomizedSpiderService struct {
	UserId bson.ObjectId
}

func (s *CreateCustomizedSpiderService) Check() (bool, error) {
	query := bson.M{
		"user_id": s.UserId,
		"type": constants.Customized,
	}
	_, count, err := model.GetSpiderList(query, 0, 1, "-_id")
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
