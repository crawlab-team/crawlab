package utils

import (
	"crawlab/constants"
	"github.com/globalsign/mgo/bson"
	"strconv"
	"time"
)

func IsObjectIdNull(id bson.ObjectId) bool {
	return id.Hex() == constants.ObjectIdNull
}

func InterfaceToString(value interface{}) string {
	switch value.(type) {
	case bson.ObjectId:
		return value.(bson.ObjectId).Hex()
	case string:
		return value.(string)
	case int:
		return strconv.Itoa(value.(int))
	case time.Time:
		return value.(time.Time).String()
	}
	return ""
}
