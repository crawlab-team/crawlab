package challenge

import (
	"crawlab/constants"
	"crawlab/model"
	"encoding/json"
	"github.com/apex/log"
	"github.com/globalsign/mgo/bson"
	"io/ioutil"
	"path"
	"runtime/debug"
)

type Service interface {
	Check() (bool, error)
}

func GetService(name string, uid bson.ObjectId) Service {
	switch name {
	case constants.ChallengeLogin7d:
		return &Login7dService{UserId: uid}
	case constants.ChallengeCreateCustomizedSpider:
		return &CreateCustomizedSpiderService{UserId: uid}
	case constants.ChallengeRunRandom:
		return &RunRandomService{UserId: uid}
	}
	return nil
}

func AddChallengeAchievement(name string, uid bson.ObjectId) error {
	ch, err := model.GetChallengeByName(name)
	if err != nil {
		return err
	}
	ca := model.ChallengeAchievement{
		ChallengeId: ch.Id,
		UserId:      uid,
	}
	if err := ca.Add(); err != nil {
		return err
	}
	return nil
}

func CheckChallengeAndUpdate(ch model.Challenge, uid bson.ObjectId) error {
	svc := GetService(ch.Name, uid)
	achieved, err := svc.Check()
	if err != nil {
		return err
	}
	if achieved && !ch.Achieved {
		if err := AddChallengeAchievement(ch.Name, uid); err != nil {
			return err
		}
	}
	return nil
}

func CheckChallengeAndUpdateAll(uid bson.ObjectId) error {
	challenges, err := model.GetChallengeListWithAchieved(nil, 0, constants.Infinite, "-_id", uid)
	if err != nil {
		return err
	}
	for _, ch := range challenges {
		if err := CheckChallengeAndUpdate(ch, uid); err != nil {
			continue
		}
	}
	return nil
}

func InitChallengeService() error {
	// 读取文件
	contentBytes, err := ioutil.ReadFile(path.Join("data", "challenge_data.json"))
	if err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return err
	}

	// 反序列化
	var challenges []model.Challenge
	if err := json.Unmarshal(contentBytes, &challenges); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return err
	}

	for _, ch := range challenges {
		chDb, err := model.GetChallengeByName(ch.Name)
		if err != nil {
			continue
		}
		if chDb.Name == "" {
			if err := ch.Add(); err != nil {
				log.Errorf(err.Error())
				debug.PrintStack()
				continue
			}
		} else {
			if err := ch.Save(); err != nil {
				log.Errorf(err.Error())
				debug.PrintStack()
				continue
			}
		}
	}

	return nil
}
