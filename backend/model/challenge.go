package model

import (
	"crawlab/database"
	"github.com/apex/log"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"runtime/debug"
	"time"
)

type Challenge struct {
	Id            bson.ObjectId `json:"_id" bson:"_id"`
	Name          string        `json:"name" bson:"name"`
	TitleCn       string        `json:"title_cn" bson:"title_cn"`
	TitleEn       string        `json:"title_en" bson:"title_en"`
	DescriptionCn string        `json:"description_cn" bson:"description_cn"`
	DescriptionEn string        `json:"description_en" bson:"description_en"`
	Difficulty    int           `json:"difficulty" bson:"difficulty"`
	Path          string        `json:"path" bson:"path"`

	// 前端展示
	Achieved bool `json:"achieved" bson:"achieved"`

	CreateTs time.Time `json:"create_ts" bson:"create_ts"`
	UpdateTs time.Time `json:"update_ts" bson:"update_ts"`
}

func (ch *Challenge) Save() error {
	s, c := database.GetCol("challenges")
	defer s.Close()

	ch.UpdateTs = time.Now()

	if err := c.UpdateId(ch.Id, c); err != nil {
		debug.PrintStack()
		return err
	}
	return nil
}

func (ch *Challenge) Add() error {
	s, c := database.GetCol("challenges")
	defer s.Close()

	ch.Id = bson.NewObjectId()
	ch.UpdateTs = time.Now()
	ch.CreateTs = time.Now()
	if err := c.Insert(ch); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return err
	}

	return nil
}

func GetChallenge(id bson.ObjectId) (Challenge, error) {
	s, c := database.GetCol("challenges")
	defer s.Close()

	var ch Challenge
	if err := c.Find(bson.M{"_id": id}).One(&ch); err != nil {
		if err != mgo.ErrNotFound {
			log.Errorf(err.Error())
			debug.PrintStack()
			return ch, err
		}
	}

	return ch, nil
}

func GetChallengeByName(name string) (Challenge, error) {
	s, c := database.GetCol("challenges")
	defer s.Close()

	var ch Challenge
	if err := c.Find(bson.M{"name": name}).One(&ch); err != nil {
		if err != mgo.ErrNotFound {
			log.Errorf(err.Error())
			debug.PrintStack()
			return ch, err
		}
	}

	return ch, nil
}

func GetChallengeList(filter interface{}, skip int, limit int, sortKey string) ([]Challenge, error) {
	s, c := database.GetCol("challenges")
	defer s.Close()

	var challenges []Challenge
	if err := c.Find(filter).Skip(skip).Limit(limit).Sort(sortKey).All(&challenges); err != nil {
		debug.PrintStack()
		return challenges, err
	}

	//for _, ch := range challenges {
	//}

	return challenges, nil
}

func GetChallengeListTotal(filter interface{}) (int, error) {
	s, c := database.GetCol("challenges")
	defer s.Close()

	var result int
	result, err := c.Find(filter).Count()
	if err != nil {
		return result, err
	}
	return result, nil
}

type ChallengeAchievement struct {
	Id          bson.ObjectId `json:"_id" bson:"_id"`
	ChallengeId bson.ObjectId `json:"challenge_id" bson:"challenge_id"`
	UserId      bson.ObjectId `json:"user_id" bson:"user_id"`

	CreateTs time.Time `json:"create_ts" bson:"create_ts"`
	UpdateTs time.Time `json:"update_ts" bson:"update_ts"`
}

func (ca *ChallengeAchievement) Save() error {
	s, c := database.GetCol("challenges_achievements")
	defer s.Close()

	ca.UpdateTs = time.Now()

	if err := c.UpdateId(ca.Id, c); err != nil {
		debug.PrintStack()
		return err
	}
	return nil
}

func (ca *ChallengeAchievement) Add() error {
	s, c := database.GetCol("challenges_achievements")
	defer s.Close()

	ca.Id = bson.NewObjectId()
	ca.UpdateTs = time.Now()
	ca.CreateTs = time.Now()
	if err := c.Insert(ca); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return err
	}

	return nil
}
