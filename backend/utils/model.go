package utils

import (
	"crawlab/constants"
	"encoding/json"
	"github.com/globalsign/mgo/bson"
	"strings"
)

func IsObjectIdNull(id bson.ObjectId) bool {
	return id.Hex() == constants.ObjectIdNull
}

func InterfaceToString(value interface{}) string {
	bytes, err := json.Marshal(value)
	if err != nil {
		return ""
	}
	str := string(bytes)
	if strings.HasPrefix(str, "\"") && strings.HasSuffix(str, "\"") {
		str = str[1 : len(str)-1]
	}
	return str
}
