package model

import (
	"crawlab/database"
	"github.com/apex/log"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
	"runtime/debug"
	"time"
)

type User struct {
	Id       bson.ObjectId `json:"_id" bson:"_id"`
	Username string        `json:"username" bson:"username"`
	Password string        `json:"password" bson:"password"`
	Role     string        `json:"role" bson:"role"`

	CreateTs time.Time `json:"create_ts" bson:"create_ts"`
	UpdateTs time.Time `json:"update_ts" bson:"update_ts"`
}

func (user *User) Add() error {
	s, c := database.GetCol("users")
	defer s.Close()

	// 如果存在用户名相同的用户，抛错
	user2, err := GetUserByUsername(user.Username)
	if err != nil {
		if err == mgo.ErrNotFound {
			// pass
		} else {
			log.Errorf(err.Error())
			debug.PrintStack()
			return err
		}
	} else {
		if user2.Username == user.Username {
			return errors.New("username already exists")
		}
	}

	user.Id = bson.NewObjectId()
	user.UpdateTs = time.Now()
	user.CreateTs = time.Now()
	if err := c.Insert(user); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return err
	}

	return nil
}

func GetUser(id bson.ObjectId) (User, error) {
	s, c := database.GetCol("users")
	defer s.Close()
	var user User
	if err := c.Find(bson.M{"_id": id}).One(&user); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return user, err
	}
	return user, nil
}

func GetUserByUsername(username string) (User, error) {
	s, c := database.GetCol("users")
	defer s.Close()

	var user User
	if err := c.Find(bson.M{"username": username}).One(&user); err != nil {
		if err != mgo.ErrNotFound {
			log.Errorf(err.Error())
			debug.PrintStack()
			return user, err
		}
	}

	return user, nil
}

func GetUserList(filter interface{}, skip int, limit int, sortKey string) ([]User, error) {
	s, c := database.GetCol("tasks")
	defer s.Close()

	var users []User
	if err := c.Find(filter).Skip(skip).Limit(limit).Sort(sortKey).All(&users); err != nil {
		debug.PrintStack()
		return users, err
	}
	return users, nil
}

func GetUserListTotal(filter interface{}) (int, error) {
	s, c := database.GetCol("users")
	defer s.Close()

	var result int
	result, err := c.Find(filter).Count()
	if err != nil {
		return result, err
	}
	return result, nil
}
