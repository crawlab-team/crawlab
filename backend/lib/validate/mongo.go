package validate

import (
	"github.com/globalsign/mgo/bson"
	"github.com/go-playground/validator/v10"
)

func MongoID(sl validator.FieldLevel) bool {
	return bson.IsObjectIdHex(sl.Field().String())
}
