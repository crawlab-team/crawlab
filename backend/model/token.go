package model

import (
	"crawlab/database"
	"github.com/apex/log"
	"github.com/globalsign/mgo/bson"
	"runtime/debug"
	"time"
)

type Token struct {
	Id       bson.ObjectId `json:"_id" bson:"_id"`
	Token    string        `json:"token" bson:"token"`
	UserId   bson.ObjectId `json:"user_id" bson:"user_id"`
	CreateTs time.Time     `json:"create_ts" bson:"create_ts"`
	UpdateTs time.Time     `json:"update_ts" bson:"update_ts"`
}

func (t *Token) Add() error {
	s, c := database.GetCol("tokens")
	defer s.Close()

	if err := c.Insert(t); err != nil {
		log.Errorf("insert token error: " + err.Error())
		debug.PrintStack()
		return err
	}

	return nil
}

func (t *Token) Delete() error {
	s, c := database.GetCol("tokens")
	defer s.Close()

	if err := c.RemoveId(t.Id); err != nil {
		log.Errorf("insert token error: " + err.Error())
		debug.PrintStack()
		return err
	}

	return nil
}

func GetTokenById(id bson.ObjectId) (t Token, err error) {
	s, c := database.GetCol("tokens")
	defer s.Close()

	if err = c.FindId(id).One(&t); err != nil {
		return t, err
	}

	return t, nil
}

func GetTokensByUserId(uid bson.ObjectId) (tokens []Token, err error) {
	s, c := database.GetCol("tokens")
	defer s.Close()

	if err = c.Find(bson.M{"user_id": uid}).All(&tokens); err != nil {
		log.Errorf("find tokens error: " + err.Error())
		debug.PrintStack()
		return tokens, err
	}

	return tokens, nil
}

func DeleteTokenById(id bson.ObjectId) error {
	t, err := GetTokenById(id)
	if err != nil {
		return err
	}

	if err := t.Delete(); err != nil {
		return err
	}

	return nil
}
