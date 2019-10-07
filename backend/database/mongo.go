package database

import (
	"github.com/globalsign/mgo"
	"github.com/spf13/viper"
	"net"
	"time"
)

var Session *mgo.Session

func GetSession() *mgo.Session {
	return Session.Copy()
}

func GetDb() (*mgo.Session, *mgo.Database) {
	s := GetSession()
	return s, s.DB(viper.GetString("mongo.db"))
}

func GetCol(collectionName string) (*mgo.Session, *mgo.Collection) {
	s := GetSession()
	db := s.DB(viper.GetString("mongo.db"))
	col := db.C(collectionName)
	return s, col
}

func GetGridFs(prefix string) (*mgo.Session, *mgo.GridFS) {
	s, db := GetDb()
	gf := db.GridFS(prefix)
	return s, gf
}

func InitMongo() error {
	var mongoHost = viper.GetString("mongo.host")
	var mongoPort = viper.GetString("mongo.port")
	var mongoDb = viper.GetString("mongo.db")
	var mongoUsername = viper.GetString("mongo.username")
	var mongoPassword = viper.GetString("mongo.password")
	var mongoAuth = viper.GetString("mongo.authSource")

	if Session == nil {
		var dialInfo mgo.DialInfo
		addr := net.JoinHostPort(mongoHost, mongoPort)
		timeout := time.Second * 10
		dialInfo = mgo.DialInfo{
			Addrs:         []string{addr},
			Timeout:       timeout,
			Database:      mongoDb,
			PoolLimit:     100,
			PoolTimeout:   timeout,
			ReadTimeout:   timeout,
			WriteTimeout:  timeout,
			AppName:       "crawlab",
			FailFast:      true,
			MinPoolSize:   10,
			MaxIdleTimeMS: 1000 * 30,
		}
		if mongoUsername != "" {
			dialInfo.Username = mongoUsername
			dialInfo.Password = mongoPassword
			dialInfo.Source = mongoAuth
		}
		sess, err := mgo.DialWithInfo(&dialInfo)
		if err != nil {
			return err
		}
		Session = sess
	}
	return nil
}
