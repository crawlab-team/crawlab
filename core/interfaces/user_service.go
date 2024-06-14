package interfaces

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService interface {
	Init() (err error)
	SetJwtSecret(secret string)
	SetJwtSigningMethod(method jwt.SigningMethod)
	Create(opts *UserCreateOptions, args ...interface{}) (err error)
	Login(opts *UserLoginOptions) (token string, u User, err error)
	CheckToken(token string) (u User, err error)
	ChangePassword(id primitive.ObjectID, password string, args ...interface{}) (err error)
	MakeToken(user User) (tokenStr string, err error)
	GetCurrentUser(c *gin.Context) (u User, err error)
}
