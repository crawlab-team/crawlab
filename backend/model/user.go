package model

import (
	"crawlab/database"
	"crawlab/utils"
	"github.com/apex/log"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"runtime/debug"
	"time"
)

type User struct {
	Id       bson.ObjectId `json:"_id" bson:"_id"`
	Username string        `json:"username" bson:"username"`
	Password string        `json:"password" bson:"password"`
	Role     string        `json:"role" bson:"role"`
	Salt     string        `json:"-"`
	Enable   bool          `json:"enable"`

	RePasswordTs time.Time `json:"-" bson:"re_password_ts"`
	LastLoginTs  time.Time `json:"-" bson:"last_login_ts"`
	CreateTs     time.Time `json:"create_ts" bson:"create_ts"`
	UpdateTs     time.Time `json:"update_ts" bson:"update_ts"`
}

func (user *User) Save() error {
	s, c := database.GetCol("users")
	defer s.Close()

	user.UpdateTs = time.Now()

	if err := c.UpdateId(user.Id, user); err != nil {
		debug.PrintStack()
		return err
	}
	return nil
}
func (user *User) ShouldUpgradeSecurityLevel() bool {
	return len(user.Salt) == 0
}
func (user *User) ValidatePassword(password string) bool {
	if user.ShouldUpgradeSecurityLevel() {
		if utils.EncryptPasswordV1(password) != user.Password {
			return false
		}
		user.Salt = utils.RandomString(10)
		user.Enable = true
		user.Password = utils.EncryptPasswordV2(password, user.Salt)
	} else if utils.EncryptPasswordV2(password, user.Salt) != user.Password {

		return false
	}
	return true
}
func (user *User) LoginWithPassword(password string) bool {

	if !user.ValidatePassword(password) {
		return false
	}
	user.LastLoginTs = time.Now()
	_ = user.Save()
	return true
}
func (user *User) Add() error {
	s, c := database.GetCol("users")
	defer s.Close()

	user.Id = bson.NewObjectId()
	if user.RePasswordTs == (time.Time{}) {
		user.RePasswordTs = time.Now().AddDate(0, 3, 0)
	}
	user.UpdateTs = time.Now()
	user.CreateTs = time.Now()
	return c.Insert(user)
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
	s, c := database.GetCol("users")
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

func UpdateUser(id bson.ObjectId, item User) error {
	s, c := database.GetCol("users")
	defer s.Close()

	var result User
	if err := c.FindId(id).One(&result); err != nil {
		debug.PrintStack()
		return err
	}

	if item.Password == "" {
		item.Password = result.Password
	} else {
		item.Password = utils.EncryptPassword(item.Password)
	}

	if err := item.Save(); err != nil {
		return err
	}
	return nil
}
func ChangePassword(user *User, unencryptedPassword string) (err error) {
	user.Salt = utils.RandomString(10)
	user.Password = utils.EncryptPasswordV2(unencryptedPassword, user.Salt)
	user.RePasswordTs = time.Now().AddDate(0, 3, 0)
	user.UpdateTs = time.Now()
	return user.Save()
}
func RemoveUser(id bson.ObjectId) error {
	s, c := database.GetCol("users")
	defer s.Close()

	var result User
	if err := c.FindId(id).One(&result); err != nil {
		return err
	}

	if err := c.RemoveId(id); err != nil {
		return err
	}

	return nil
}
func CreateUserIndex() error {
	s, c := database.GetCol("users")
	defer s.Close()
	usernameUniqueIndex := mgo.Index{
		Key:              []string{"username"},
		Unique:           true,
		DropDups:         false,
		Background:       true,
		Sparse:           false,
		PartialFilter:    nil,
		ExpireAfter:      0,
		Name:             "",
		Min:              0,
		Max:              0,
		Minf:             0,
		Maxf:             0,
		BucketSize:       0,
		Bits:             0,
		DefaultLanguage:  "",
		LanguageOverride: "",
		Weights:          nil,
		Collation:        nil,
	}

	return c.EnsureIndex(usernameUniqueIndex)
}
