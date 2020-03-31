package challenge

import (
	"crawlab/constants"
	"crawlab/model"
	"github.com/globalsign/mgo/bson"
)

type RunRandomService struct {
	UserId bson.ObjectId
}

func (s *RunRandomService) Check() (bool, error) {
	query := bson.M{
		"user_id":     s.UserId,
		"run_type":    constants.RunTypeRandom,
		"status":      constants.StatusFinished,
		"schedule_id": bson.ObjectIdHex(constants.ObjectIdNull),
	}
	list, err := model.GetTaskList(query, 0, 1, "-_id")
	if err != nil {
		return false, err
	}
	return len(list) > 0, nil
}
