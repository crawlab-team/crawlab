package services

import (
	"crawlab/constants"
	"crawlab/model"
	"crawlab/utils"
	"github.com/apex/log"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"runtime/debug"
	"time"
)

func InitUserService() error {
	adminUser := model.User{
		Username: "admin",
		Password: utils.EncryptPassword("admin"),
		Role:     constants.RoleAdmin,
	}
	if err := adminUser.Add(); err != nil {
		// pass
	}
	return nil
}

func GetToken(username string) (tokenStr string, err error) {
	user, err := model.GetUserByUsername(username)
	if err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.Id,
		"username": user.Username,
		"nbf":      time.Now().Unix(),
	})

	tokenStr, err = token.SignedString([]byte(viper.GetString("server.secret")))
	if err != nil {
		return
	}
	return
}
