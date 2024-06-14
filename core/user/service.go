package user

import (
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/container"
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/delegate"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/core/utils"
	mongo2 "github.com/crawlab-team/crawlab/db/mongo"
	"github.com/crawlab-team/crawlab/trace"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Service struct {
	// settings variables
	jwtSecret        string
	jwtSigningMethod jwt.SigningMethod

	// dependencies
	modelSvc service.ModelService
}

func (svc *Service) Init() (err error) {
	_, err = svc.modelSvc.GetUserByUsername(constants.DefaultAdminUsername, nil)
	if err == nil {
		return nil
	}
	if err.Error() != mongo.ErrNoDocuments.Error() {
		return err
	}
	return svc.Create(&interfaces.UserCreateOptions{
		Username: constants.DefaultAdminUsername,
		Password: constants.DefaultAdminPassword,
		Role:     constants.RoleAdmin,
	})
}

func (svc *Service) SetJwtSecret(secret string) {
	svc.jwtSecret = secret
}

func (svc *Service) SetJwtSigningMethod(method jwt.SigningMethod) {
	svc.jwtSigningMethod = method
}

func (svc *Service) Create(opts *interfaces.UserCreateOptions, args ...interface{}) (err error) {
	actor := utils.GetUserFromArgs(args...)

	// validate options
	if opts.Username == "" || opts.Password == "" {
		return trace.TraceError(errors.ErrorUserMissingRequiredFields)
	}
	if len(opts.Password) < 5 {
		return trace.TraceError(errors.ErrorUserInvalidPassword)
	}

	// normalize options
	if opts.Role == "" {
		opts.Role = constants.RoleNormal
	}

	// check if user exists
	if u, err := svc.modelSvc.GetUserByUsername(opts.Username, nil); err == nil && u != nil && !u.Id.IsZero() {
		return trace.TraceError(errors.ErrorUserAlreadyExists)
	}

	// transaction
	return mongo2.RunTransaction(func(ctx mongo.SessionContext) error {
		// add user
		u := &models.User{
			Username: opts.Username,
			Role:     opts.Role,
			Email:    opts.Email,
		}
		if err := delegate.NewModelDelegate(u, actor).Add(); err != nil {
			return err
		}

		// add password
		p := &models.Password{
			Id:       u.Id,
			Password: utils.EncryptMd5(opts.Password),
		}
		if err := delegate.NewModelDelegate(p, actor).Add(); err != nil {
			return err
		}

		return nil
	})
}

func (svc *Service) Login(opts *interfaces.UserLoginOptions) (token string, u interfaces.User, err error) {
	u, err = svc.modelSvc.GetUserByUsername(opts.Username, nil)
	if err != nil {
		return "", nil, err
	}
	p, err := svc.modelSvc.GetPasswordById(u.GetId())
	if err != nil {
		return "", nil, err
	}
	if p.Password != utils.EncryptMd5(opts.Password) {
		return "", nil, errors.ErrorUserMismatch
	}
	token, err = svc.makeToken(u)
	if err != nil {
		return "", nil, err
	}
	return token, u, nil
}

func (svc *Service) CheckToken(tokenStr string) (u interfaces.User, err error) {
	return svc.checkToken(tokenStr)
}

func (svc *Service) ChangePassword(id primitive.ObjectID, password string, args ...interface{}) (err error) {
	actor := utils.GetUserFromArgs(args...)

	p, err := svc.modelSvc.GetPasswordById(id)
	if err != nil {
		return err
	}
	p.Password = utils.EncryptMd5(password)
	if err := delegate.NewModelDelegate(p, actor).Save(); err != nil {
		return err
	}
	return nil
}

func (svc *Service) MakeToken(user interfaces.User) (tokenStr string, err error) {
	return svc.makeToken(user)
}

func (svc *Service) GetCurrentUser(c *gin.Context) (user interfaces.User, err error) {
	// token string
	tokenStr := c.GetHeader("Authorization")

	// user
	u, err := userSvc.CheckToken(tokenStr)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (svc *Service) makeToken(user interfaces.User) (tokenStr string, err error) {
	token := jwt.NewWithClaims(svc.jwtSigningMethod, jwt.MapClaims{
		"id":       user.GetId(),
		"username": user.GetUsername(),
		"nbf":      time.Now().Unix(),
	})
	return token.SignedString([]byte(svc.jwtSecret))
}

func (svc *Service) checkToken(tokenStr string) (user interfaces.User, err error) {
	token, err := jwt.Parse(tokenStr, svc.getSecretFunc())
	if err != nil {
		return
	}

	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.ErrorUserInvalidType
		return
	}

	if !token.Valid {
		err = errors.ErrorUserInvalidToken
		return
	}

	id, err := primitive.ObjectIDFromHex(claim["id"].(string))
	if err != nil {
		return user, err
	}
	username := claim["username"].(string)
	user, err = svc.modelSvc.GetUserById(id)
	if err != nil {
		err = errors.ErrorUserNotExists
		return
	}

	if username != user.GetUsername() {
		err = errors.ErrorUserMismatch
		return
	}

	return
}

func (svc *Service) getSecretFunc() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(svc.jwtSecret), nil
	}
}

func NewUserService() (svc2 interfaces.UserService, err error) {
	// service
	svc := &Service{
		jwtSecret:        "crawlab",
		jwtSigningMethod: jwt.SigningMethodHS256,
	}

	// dependency injection
	if err := container.GetContainer().Invoke(func(modelSvc service.ModelService) {
		svc.modelSvc = modelSvc
	}); err != nil {
		return nil, trace.TraceError(err)
	}

	// initialize
	if err := svc.Init(); err != nil {
		return nil, trace.TraceError(err)
	}

	return svc, nil
}

var userSvc interfaces.UserService

func GetUserService() (svc interfaces.UserService, err error) {
	if userSvc != nil {
		return userSvc, nil
	}
	svc, err = NewUserService()
	if err != nil {
		return nil, err
	}
	userSvc = svc
	return svc, nil
}
