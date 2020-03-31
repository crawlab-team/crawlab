package challenge

import (
	"crawlab/model"
	"github.com/globalsign/mgo/bson"
)

type CreateUserService struct {
	UserId bson.ObjectId
}

func (s *CreateUserService) Check() (bool, error) {
	query := bson.M{
		"user_id": s.UserId,
	}
	list, err := model.GetUserList(query, 0, 1, "-_id")
	if err != nil {
		return false, err
	}
	return len(list) > 0, nil
}
