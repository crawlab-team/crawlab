package controllers_test

import (
	"bytes"
	"encoding/json"
	"github.com/crawlab-team/crawlab/core/controllers"
	"github.com/crawlab-team/crawlab/core/middlewares"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateSpider(t *testing.T) {
	SetupTestDB()
	defer CleanupTestDB()

	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.Use(middlewares.AuthorizationMiddlewareV2())
	router.POST("/spiders", controllers.PostSpider)

	payload := models.SpiderV2{
		Name:    "Test Spider",
		ColName: "test_spiders",
	}
	jsonValue, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/spiders", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", TestToken)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response controllers.Response[models.SpiderV2]
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	require.Nil(t, err)
	assert.False(t, response.Data.Id.IsZero())
	assert.Equal(t, payload.Name, response.Data.Name)
	assert.False(t, response.Data.ColId.IsZero())
}

func TestGetSpiderById(t *testing.T) {
	SetupTestDB()
	defer CleanupTestDB()

	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.Use(middlewares.AuthorizationMiddlewareV2())
	router.GET("/spiders/:id", controllers.GetSpiderById)

	model := models.SpiderV2{
		Name:    "Test Spider",
		ColName: "test_spiders",
	}
	id, err := service.NewModelServiceV2[models.SpiderV2]().InsertOne(model)
	require.Nil(t, err)
	ts := models.SpiderStatV2{}
	ts.SetId(id)
	_, err = service.NewModelServiceV2[models.SpiderStatV2]().InsertOne(ts)
	require.Nil(t, err)

	req, _ := http.NewRequest("GET", "/spiders/"+id.Hex(), nil)
	req.Header.Set("Authorization", TestToken)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response controllers.Response[models.SpiderV2]
	err = json.Unmarshal(resp.Body.Bytes(), &response)
	require.Nil(t, err)
	assert.Equal(t, model.Name, response.Data.Name)
}

func TestUpdateSpiderById(t *testing.T) {
	SetupTestDB()
	defer CleanupTestDB()

	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.Use(middlewares.AuthorizationMiddlewareV2())
	router.PUT("/spiders/:id", controllers.PutSpiderById)

	model := models.SpiderV2{
		Name:    "Test Spider",
		ColName: "test_spiders",
	}
	id, err := service.NewModelServiceV2[models.SpiderV2]().InsertOne(model)
	require.Nil(t, err)
	ts := models.SpiderStatV2{}
	ts.SetId(id)
	_, err = service.NewModelServiceV2[models.SpiderStatV2]().InsertOne(ts)
	require.Nil(t, err)

	spiderId := id.Hex()
	payload := models.SpiderV2{
		Name:    "Updated Spider",
		ColName: "test_spider",
	}
	payload.SetId(id)
	jsonValue, _ := json.Marshal(payload)
	req, _ := http.NewRequest("PUT", "/spiders/"+spiderId, bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", TestToken)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response controllers.Response[models.SpiderV2]
	err = json.Unmarshal(resp.Body.Bytes(), &response)
	require.Nil(t, err)
	assert.Equal(t, payload.Name, response.Data.Name)

	svc := service.NewModelServiceV2[models.SpiderV2]()
	resModel, err := svc.GetById(id)
	require.Nil(t, err)
	assert.Equal(t, payload.Name, resModel.Name)
}

func TestDeleteSpiderById(t *testing.T) {
	SetupTestDB()
	defer CleanupTestDB()

	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.Use(middlewares.AuthorizationMiddlewareV2())
	router.DELETE("/spiders/:id", controllers.DeleteSpiderById)

	model := models.SpiderV2{
		Name:    "Test Spider",
		ColName: "test_spiders",
	}
	id, err := service.NewModelServiceV2[models.SpiderV2]().InsertOne(model)
	require.Nil(t, err)
	ts := models.SpiderStatV2{}
	ts.SetId(id)
	_, err = service.NewModelServiceV2[models.SpiderStatV2]().InsertOne(ts)
	require.Nil(t, err)
	task := models.TaskV2{}
	task.SpiderId = id
	taskId, err := service.NewModelServiceV2[models.TaskV2]().InsertOne(task)
	require.Nil(t, err)
	taskStat := models.TaskStatV2{}
	taskStat.SetId(taskId)
	_, err = service.NewModelServiceV2[models.TaskStatV2]().InsertOne(taskStat)
	require.Nil(t, err)

	req, _ := http.NewRequest("DELETE", "/spiders/"+id.Hex(), nil)
	req.Header.Set("Authorization", TestToken)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	_, err = service.NewModelServiceV2[models.SpiderV2]().GetById(id)
	assert.NotNil(t, err)
	_, err = service.NewModelServiceV2[models.SpiderStatV2]().GetById(id)
	assert.NotNil(t, err)
	taskCount, err := service.NewModelServiceV2[models.TaskV2]().Count(bson.M{"spider_id": id})
	require.Nil(t, err)
	assert.Equal(t, 0, taskCount)
	taskStatCount, err := service.NewModelServiceV2[models.TaskStatV2]().Count(bson.M{"_id": taskId})
	require.Nil(t, err)
	assert.Equal(t, 0, taskStatCount)

}

func TestDeleteSpiderList(t *testing.T) {
	SetupTestDB()
	defer CleanupTestDB()

	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.Use(middlewares.AuthorizationMiddlewareV2())
	router.DELETE("/spiders", controllers.DeleteSpiderList)

	modelList := []models.SpiderV2{
		{
			Name:    "Test Name 1",
			ColName: "test_spiders",
		}, {
			Name:    "Test Name 2",
			ColName: "test_spiders",
		},
	}
	var ids []primitive.ObjectID
	var taskIds []primitive.ObjectID
	for _, model := range modelList {
		id, err := service.NewModelServiceV2[models.SpiderV2]().InsertOne(model)
		require.Nil(t, err)
		ts := models.SpiderStatV2{}
		ts.SetId(id)
		_, err = service.NewModelServiceV2[models.SpiderStatV2]().InsertOne(ts)
		require.Nil(t, err)
		task := models.TaskV2{}
		task.SpiderId = id
		taskId, err := service.NewModelServiceV2[models.TaskV2]().InsertOne(task)
		require.Nil(t, err)
		taskStat := models.TaskStatV2{}
		taskStat.SetId(taskId)
		_, err = service.NewModelServiceV2[models.TaskStatV2]().InsertOne(taskStat)
		require.Nil(t, err)
		ids = append(ids, id)
		taskIds = append(taskIds, taskId)
	}

	payload := struct {
		Ids []primitive.ObjectID `json:"ids"`
	}{
		Ids: ids,
	}
	jsonValue, _ := json.Marshal(payload)
	req, _ := http.NewRequest("DELETE", "/spiders", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", TestToken)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	spiderCount, err := service.NewModelServiceV2[models.SpiderV2]().Count(bson.M{"_id": bson.M{"$in": ids}})
	require.Nil(t, err)
	assert.Equal(t, 0, spiderCount)
	spiderStatCount, err := service.NewModelServiceV2[models.SpiderStatV2]().Count(bson.M{"_id": bson.M{"$in": ids}})
	require.Nil(t, err)
	assert.Equal(t, 0, spiderStatCount)
	taskCount, err := service.NewModelServiceV2[models.TaskV2]().Count(bson.M{"_id": bson.M{"$in": taskIds}})
	require.Nil(t, err)
	assert.Equal(t, 0, taskCount)
	taskStatCount, err := service.NewModelServiceV2[models.TaskStatV2]().Count(bson.M{"_id": bson.M{"$in": taskIds}})
	require.Nil(t, err)
	assert.Equal(t, 0, taskStatCount)
}
