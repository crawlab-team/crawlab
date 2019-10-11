package services

import (
	"crawlab/constants"
	"crawlab/model"
	"crawlab/utils"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo/bson"
	"github.com/spf13/viper"
	"time"
)

func InitUserService() error {
	adminUser := model.User{
		Username: "admin",
		Password: utils.EncryptPassword("admin"),
		Role:     constants.RoleAdmin,
	}
	_ = adminUser.Add()
	return nil
}
func MakeToken(user *model.User) (tokenStr string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.Id,
		"username": user.Username,
		"nbf":      time.Now().Unix(),
	})

	return token.SignedString([]byte(viper.GetString("server.secret")))

}

//func GetToken(username string) (tokenStr string, err error) {
//	user, err := model.GetUserByUsername(username)
//	if err != nil {
//		log.Errorf(err.Error())
//		debug.PrintStack()
//		return
//	}
//
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
//		"id":       user.Id,
//		"username": user.Username,
//		"nbf":      time.Now().Unix(),
//	})
//
//	tokenStr, err = token.SignedString([]byte(viper.GetString("server.secret")))
//	if err != nil {
//		return
//	}
//	return
//}

func SecretFunc() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("server.secret")), nil
	}
}

func CheckToken(tokenStr string) (user model.User, err error) {
	token, err := jwt.Parse(tokenStr, SecretFunc())
	if err != nil {
		return
	}

	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("cannot convert claim to mapclaim")
		return
	}

	//验证token，如果token被修改过则为false
	if !token.Valid {
		err = errors.New("token is invalid")
		return
	}

	id := bson.ObjectIdHex(claim["id"].(string))
	username := claim["username"].(string)
	user, err = model.GetUser(id)
	if err != nil {
		err = errors.New("cannot get user")
		return
	}

	if username != user.Username {
		err = errors.New("username does not match")
		return
	}

	return
}
