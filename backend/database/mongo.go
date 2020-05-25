package database

import (
	"crawlab/constants"
	"github.com/apex/log"
	"github.com/cenkalti/backoff/v4"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/spf13/viper"
	"net"
	"reflect"
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

func FillNullObjectId(doc interface{}) {
	t := reflect.TypeOf(doc)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return
	}
	v := reflect.ValueOf(doc)
	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)
		fv := v.Elem().Field(i)
		val := fv.Interface()
		switch val.(type) {
		case bson.ObjectId:
			if !val.(bson.ObjectId).Valid() {
				v.FieldByName(ft.Name).Set(reflect.ValueOf(bson.ObjectIdHex(constants.ObjectIdNull)))
			}
		}
	}
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
		bp := backoff.NewExponentialBackOff()
		var err error

		err = backoff.Retry(func() error {
			Session, err = mgo.DialWithInfo(&dialInfo)
			if err != nil {
				log.WithError(err).Warnf("waiting for connect mongo database, after %f seconds try again.", bp.NextBackOff().Seconds())
			}
			return err
		}, bp)
	}
	//Add Unique index for 'key'
	keyIndex := mgo.Index{
		Key:    []string{"key"},
		Unique: true,
	}
	s, c := GetCol("nodes")
	defer s.Close()
	c.EnsureIndex(keyIndex)

	return nil
}
