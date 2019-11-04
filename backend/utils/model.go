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
	switch realValue := value.(type) {
	case bson.ObjectId:
		return realValue.Hex()
	case string:
		return realValue
	case int:
		return strconv.Itoa(realValue)
	case time.Time:
		return realValue.String()
	default:
		return ""
	}
}
