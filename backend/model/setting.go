package model

import (
	"crawlab/database"
	"github.com/globalsign/mgo/bson"
	"time"
)

type Setting struct {
	Keyword  string
	Document bson.Raw
}

func GetRawSetting(keyword string, pointer interface{}) error {
	s, col := database.GetCol("settings")
	defer s.Close()
	var setting Setting
	err := col.Find(bson.M{"keyword": keyword}).One(&setting)
	if err != nil {
		return err
	}
	return setting.Document.Unmarshal(pointer)
}

type DocumentMeta struct {
	DocumentVersion  int
	DocStructVersion int
	UpdateTime       time.Time
	CreateTime       time.Time
	DeleteTime       time.Time
}

//demo
type SecuritySetting struct {
	EnableRegister   bool
	EnableInvitation bool
	DocumentMeta     `bson:"inline" json:"inline"`
}

func GetSecuritySetting() (SecuritySetting, error) {
	var app SecuritySetting
	err := GetRawSetting("security", &app)
	return app, err

}
