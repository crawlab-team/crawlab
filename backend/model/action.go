package model

import (
	"crawlab/database"
	"github.com/apex/log"
	"github.com/globalsign/mgo/bson"
	"runtime/debug"
	"time"
)

type Action struct {
	Id     bson.ObjectId `json:"_id" bson:"_id"`
	UserId bson.ObjectId `json:"user_id" bson:"user_id"`
	Type   string        `json:"type" bson:"type"`

	CreateTs time.Time `json:"create_ts" bson:"create_ts"`
	UpdateTs time.Time `json:"update_ts" bson:"update_ts"`
}

func (a *Action) Save() error {
	s, c := database.GetCol("actions")
	defer s.Close()

	a.UpdateTs = time.Now()

	if err := c.UpdateId(a.Id, a); err != nil {
		debug.PrintStack()
		return err
	}
	return nil
}

func (a *Action) Add() error {
	s, c := database.GetCol("actions")
	defer s.Close()

	a.Id = bson.NewObjectId()
	a.UpdateTs = time.Now()
	a.CreateTs = time.Now()
	if err := c.Insert(a); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return err
	}

	return nil
}

func GetAction(id bson.ObjectId) (Action, error) {
	s, c := database.GetCol("actions")
	defer s.Close()
	var user Action
	if err := c.Find(bson.M{"_id": id}).One(&user); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return user, err
	}
	return user, nil
}

func GetActionList(filter interface{}, skip int, limit int, sortKey string) ([]Action, error) {
	s, c := database.GetCol("actions")
	defer s.Close()

	var actions []Action
	if err := c.Find(filter).Skip(skip).Limit(limit).Sort(sortKey).All(&actions); err != nil {
		debug.PrintStack()
		return actions, err
	}
	return actions, nil
}

func GetActionListTotal(filter interface{}) (int, error) {
	s, c := database.GetCol("actions")
	defer s.Close()

	var result int
	result, err := c.Find(filter).Count()
	if err != nil {
		return result, err
	}
	return result, nil
}

func UpdateAction(id bson.ObjectId, item Action) error {
	s, c := database.GetCol("actions")
	defer s.Close()

	var result Action
	if err := c.FindId(id).One(&result); err != nil {
		debug.PrintStack()
		return err
	}

	if err := item.Save(); err != nil {
		return err
	}
	return nil
}

func RemoveAction(id bson.ObjectId) error {
	s, c := database.GetCol("actions")
	defer s.Close()

	var result Action
	if err := c.FindId(id).One(&result); err != nil {
		return err
	}

	if err := c.RemoveId(id); err != nil {
		return err
	}

	return nil
}
