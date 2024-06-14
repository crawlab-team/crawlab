package controllers

import (
	"encoding/json"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/container"
	"github.com/crawlab-team/crawlab/core/entity"
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	delegate2 "github.com/crawlab-team/crawlab/core/models/delegate"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/core/utils"
	"github.com/crawlab-team/crawlab/trace"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

var UserController *userController

func getUserActions() []Action {
	userCtx := newUserContext()
	return []Action{
		{
			Method:      http.MethodPost,
			Path:        "/:id/change-password",
			HandlerFunc: userCtx.changePassword,
		},
		{
			Method:      http.MethodGet,
			Path:        "/me",
			HandlerFunc: userCtx.getMe,
		},
		{
			Method:      http.MethodPut,
			Path:        "/me",
			HandlerFunc: userCtx.putMe,
		},
	}
}

type userController struct {
	ListActionControllerDelegate
	d   ListActionControllerDelegate
	ctx *userContext
}

func (ctr *userController) Post(c *gin.Context) {
	var u models.User
	if err := c.ShouldBindJSON(&u); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	if err := ctr.ctx.userSvc.Create(&interfaces.UserCreateOptions{
		Username: u.Username,
		Password: u.Password,
		Email:    u.Email,
		Role:     u.Role,
	}); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccess(c)
}

func (ctr *userController) PostList(c *gin.Context) {
	// users
	var users []models.User
	if err := c.ShouldBindJSON(&users); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	for _, u := range users {
		if err := ctr.ctx.userSvc.Create(&interfaces.UserCreateOptions{
			Username: u.Username,
			Password: u.Password,
			Email:    u.Email,
			Role:     u.Role,
		}); err != nil {
			trace.PrintError(err)
		}
	}

	HandleSuccess(c)
}

func (ctr *userController) PutList(c *gin.Context) {
	// payload
	var payload entity.BatchRequestPayloadWithStringData
	if err := c.ShouldBindJSON(&payload); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// doc to update
	var doc models.User
	if err := json.Unmarshal([]byte(payload.Data), &doc); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// query
	query := bson.M{
		"_id": bson.M{
			"$in": payload.Ids,
		},
	}

	// update users
	if err := ctr.ctx.modelSvc.GetBaseService(interfaces.ModelIdUser).UpdateDoc(query, &doc, payload.Fields); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// update passwords
	if utils.Contains(payload.Fields, "password") {
		for _, id := range payload.Ids {
			if err := ctr.ctx.userSvc.ChangePassword(id, doc.Password); err != nil {
				trace.PrintError(err)
			}
		}
	}

	HandleSuccess(c)
}

type userContext struct {
	modelSvc service.ModelService
	userSvc  interfaces.UserService
}

func (ctx *userContext) changePassword(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	var payload map[string]string
	if err := c.ShouldBindJSON(&payload); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	password, ok := payload["password"]
	if !ok {
		HandleErrorBadRequest(c, errors.ErrorUserMissingRequiredFields)
		return
	}
	if len(password) < 5 {
		HandleErrorBadRequest(c, errors.ErrorUserInvalidPassword)
		return
	}
	if err := ctx.userSvc.ChangePassword(id, password); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccess(c)
}

func (ctx *userContext) getMe(c *gin.Context) {
	u, err := ctx._getMe(c)
	if err != nil {
		HandleErrorUnauthorized(c, errors.ErrorUserUnauthorized)
		return
	}
	HandleSuccessWithData(c, u)
}

func (ctx *userContext) putMe(c *gin.Context) {
	// current user
	u, err := ctx._getMe(c)
	if err != nil {
		HandleErrorUnauthorized(c, errors.ErrorUserUnauthorized)
		return
	}

	// payload
	doc, err := NewJsonBinder(ControllerIdUser).Bind(c)
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	if doc.GetId() != u.GetId() {
		HandleErrorBadRequest(c, errors.ErrorHttpBadRequest)
		return
	}

	// save to db
	if err := delegate2.NewModelDelegate(doc, GetUserFromContext(c)).Save(); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccessWithData(c, doc)
}

func (ctx *userContext) _getMe(c *gin.Context) (u interfaces.User, err error) {
	res, ok := c.Get(constants.UserContextKey)
	if !ok {
		return nil, trace.TraceError(errors.ErrorUserNotExistsInContext)
	}
	u, ok = res.(interfaces.User)
	if !ok {
		return nil, trace.TraceError(errors.ErrorUserInvalidType)
	}
	return u, nil
}

func newUserContext() *userContext {
	// context
	ctx := &userContext{}

	// dependency injection
	if err := container.GetContainer().Invoke(func(
		modelSvc service.ModelService,
		userSvc interfaces.UserService,
	) {
		ctx.modelSvc = modelSvc
		ctx.userSvc = userSvc
	}); err != nil {
		panic(err)
	}

	return ctx
}

func newUserController() *userController {
	modelSvc, err := service.GetService()
	if err != nil {
		panic(err)
	}

	ctr := NewListPostActionControllerDelegate(ControllerIdUser, modelSvc.GetBaseService(interfaces.ModelIdUser), getUserActions())
	d := NewListPostActionControllerDelegate(ControllerIdUser, modelSvc.GetBaseService(interfaces.ModelIdUser), getUserActions())
	ctx := newUserContext()

	return &userController{
		ListActionControllerDelegate: *ctr,
		d:                            *d,
		ctx:                          ctx,
	}
}
