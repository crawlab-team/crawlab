package model

import (
	"crawlab/database"
	"errors"
	"github.com/apex/log"
	"github.com/globalsign/mgo/bson"
	"runtime/debug"
)

/**
全局变量
*/

type Variable struct {
	Id     bson.ObjectId `json:"_id" bson:"_id"`
	Key    string        `json:"key" bson:"key"`
	Value  string        `json:"value" bson:"value"`
	Remark string        `json:"remark" bson:"remark"`
}

func (model *Variable) Save() error {
	s, c := database.GetCol("variable")
	defer s.Close()

	if err := c.UpdateId(model.Id, model); err != nil {
		log.Errorf("update variable error: %s", err.Error())
		return err
	}
	return nil
}

func (model *Variable) Add() error {
	s, c := database.GetCol("variable")
	defer s.Close()

	// key 去重
	_, err := GetByKey(model.Key)
	if err == nil {
		return errors.New("key already exists")
	}

	model.Id = bson.NewObjectId()
	if err := c.Insert(model); err != nil {
		log.Errorf("add variable error: %s", err.Error())
		debug.PrintStack()
		return err
	}
	return nil
}

func (model *Variable) Delete() error {
	s, c := database.GetCol("variable")
	defer s.Close()

	if err := c.RemoveId(model.Id); err != nil {
		log.Errorf("remove variable error: %s", err.Error())
		debug.PrintStack()
		return err
	}
	return nil
}

func GetByKey(key string) (Variable, error) {
	s, c := database.GetCol("variable")
	defer s.Close()

	var model Variable
	if err := c.Find(bson.M{"key": key}).One(&model); err != nil {
		log.Errorf("variable found error: %s, key: %s", err.Error(), key)
		return model, err
	}
	return model, nil
}

func GetVariable(id bson.ObjectId) (Variable, error) {
	s, c := database.GetCol("variable")
	defer s.Close()

	var model Variable
	if err := c.FindId(id).One(&model); err != nil {
		log.Errorf("variable found error: %s", err.Error())
		return model, err
	}
	return model, nil
}

func GetVariableList() []Variable {
	s, c := database.GetCol("variable")
	defer s.Close()

	var list []Variable
	if err := c.Find(nil).All(&list); err != nil {

	}
	return list
}
