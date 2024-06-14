package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/crawlab-team/crawlab/core/controllers"
	"github.com/crawlab-team/crawlab/core/middlewares"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/core/user"
	"github.com/spf13/viper"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/crawlab-team/crawlab/db/mongo"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {

}

// TestModel is a simple struct to be used as a model in tests
type TestModel models.TestModel

var TestToken string

// SetupTestDB sets up the test database
func SetupTestDB() {
	viper.Set("mongo.db", "testdb")
	modelSvc := service.NewModelServiceV2[models.UserV2]()
	u := models.UserV2{
		Username: "admin",
	}
	id, err := modelSvc.InsertOne(u)
	if err != nil {
		panic(err)
	}
	u.SetId(id)

	userSvc, err := user.GetUserServiceV2()
	if err != nil {
		panic(err)
	}
	token, err := userSvc.MakeToken(&u)
	if err != nil {
		panic(err)
	}
	TestToken = token
}

// SetupRouter sets up the gin router for testing
func SetupRouter() *gin.Engine {
	router := gin.Default()
	return router
}

// CleanupTestDB cleans up the test database
func CleanupTestDB() {
	mongo.GetMongoDb("testdb").Drop(context.Background())
}

func TestBaseControllerV2_GetById(t *testing.T) {
	SetupTestDB()
	defer CleanupTestDB()

	// Insert a test document
	id, err := service.NewModelServiceV2[TestModel]().InsertOne(TestModel{Name: "test"})
	assert.NoError(t, err)

	// Initialize the controller
	ctr := controllers.NewControllerV2[TestModel]()

	// Set up the router
	router := SetupRouter()
	router.Use(middlewares.AuthorizationMiddlewareV2())
	router.GET("/testmodels/:id", ctr.GetById)

	// Create a test request
	req, _ := http.NewRequest("GET", "/testmodels/"+id.Hex(), nil)
	req.Header.Set("Authorization", TestToken)
	w := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(w, req)

	// Check the response
	assert.Equal(t, http.StatusOK, w.Code)

	var response controllers.Response[TestModel]
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "test", response.Data.Name)
}

func TestBaseControllerV2_Post(t *testing.T) {
	SetupTestDB()
	defer CleanupTestDB()

	// Initialize the controller
	ctr := controllers.NewControllerV2[TestModel]()

	// Set up the router
	router := SetupRouter()
	router.Use(middlewares.AuthorizationMiddlewareV2())
	router.POST("/testmodels", ctr.Post)

	// Create a test request
	testModel := TestModel{Name: "test"}
	jsonValue, _ := json.Marshal(testModel)
	req, _ := http.NewRequest("POST", "/testmodels", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", TestToken)
	w := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(w, req)

	// Check the response
	assert.Equal(t, http.StatusOK, w.Code)

	var response controllers.Response[TestModel]
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "test", response.Data.Name)

	// Check if the document was inserted into the database
	result, err := service.NewModelServiceV2[TestModel]().GetById(response.Data.Id)
	assert.NoError(t, err)
	assert.Equal(t, "test", result.Name)
}

func TestBaseControllerV2_DeleteById(t *testing.T) {
	SetupTestDB()
	defer CleanupTestDB()

	// Insert a test document
	id, err := service.NewModelServiceV2[TestModel]().InsertOne(TestModel{Name: "test"})
	assert.NoError(t, err)

	// Initialize the controller
	ctr := controllers.NewControllerV2[TestModel]()

	// Set up the router
	router := SetupRouter()
	router.Use(middlewares.AuthorizationMiddlewareV2())
	router.DELETE("/testmodels/:id", ctr.DeleteById)

	// Create a test request
	req, _ := http.NewRequest("DELETE", "/testmodels/"+id.Hex(), nil)
	req.Header.Set("Authorization", TestToken)
	w := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(w, req)

	// Check the response
	assert.Equal(t, http.StatusOK, w.Code)

	// Check if the document was deleted from the database
	_, err = service.NewModelServiceV2[TestModel]().GetById(id)
	assert.Error(t, err)
}
