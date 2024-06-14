package controllers_test

import (
	"github.com/crawlab-team/crawlab/core/controllers"
	"github.com/crawlab-team/crawlab/core/middlewares"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPostUserChangePassword_Success(t *testing.T) {
	SetupTestDB()
	defer CleanupTestDB()

	modelSvc := service.NewModelServiceV2[models.UserV2]()
	u := models.UserV2{}
	id, err := modelSvc.InsertOne(u)
	assert.Nil(t, err)
	u.SetId(id)

	router := gin.Default()
	router.Use(middlewares.AuthorizationMiddlewareV2())
	router.POST("/users/:id/change-password", controllers.PostUserChangePassword)

	password := "newPassword"
	reqBody := strings.NewReader(`{"password":"` + password + `"}`)
	req, _ := http.NewRequest(http.MethodPost, "/users/"+id.Hex()+"/change-password", reqBody)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", TestToken)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetUserMe_Success(t *testing.T) {
	SetupTestDB()
	defer CleanupTestDB()

	modelSvc := service.NewModelServiceV2[models.UserV2]()
	u := models.UserV2{}
	id, err := modelSvc.InsertOne(u)
	assert.Nil(t, err)
	u.SetId(id)

	router := gin.Default()
	router.Use(middlewares.AuthorizationMiddlewareV2())
	router.GET("/users/me", controllers.GetUserMe)

	req, _ := http.NewRequest(http.MethodGet, "/users/me", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", TestToken)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPutUserById_Success(t *testing.T) {
	SetupTestDB()
	defer CleanupTestDB()

	modelSvc := service.NewModelServiceV2[models.UserV2]()
	u := models.UserV2{}
	id, err := modelSvc.InsertOne(u)
	assert.Nil(t, err)
	u.SetId(id)

	router := gin.Default()
	router.Use(middlewares.AuthorizationMiddlewareV2())
	router.PUT("/users/me", controllers.PutUserById)

	reqBody := strings.NewReader(`{"id":"` + id.Hex() + `","username":"newUsername","email":"newEmail@test.com"}`)
	req, _ := http.NewRequest(http.MethodPut, "/users/me", reqBody)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", TestToken)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
