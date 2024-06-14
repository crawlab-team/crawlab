package client_test

import (
	"context"
	"github.com/crawlab-team/crawlab/core/grpc/server"
	"github.com/crawlab-team/crawlab/core/models/client"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/db/mongo"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
	"time"
)

type TestModel models.TestModel

func setupTestDB() {
	viper.Set("mongo.db", "testdb")
}

func teardownTestDB() {
	db := mongo.GetMongoDb("testdb")
	db.Drop(context.Background())
}

func TestModelServiceV2_GetById(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()
	svr, err := server.NewGrpcServerV2()
	require.Nil(t, err)
	go svr.Start()
	defer svr.Stop()

	m := TestModel{
		Name: "Test Name",
	}
	modelSvc := service.NewModelServiceV2[TestModel]()
	id, err := modelSvc.InsertOne(m)
	require.Nil(t, err)
	m.SetId(id)
	time.Sleep(100 * time.Millisecond)

	c, err := grpc.Dial("localhost:9666", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.Nil(t, err)
	c.Connect()

	clientSvc := client.NewModelServiceV2[TestModel]()
	res, err := clientSvc.GetById(m.Id)
	require.Nil(t, err)
	assert.Equal(t, res.Id, m.Id)
	assert.Equal(t, res.Name, m.Name)
}

func TestModelServiceV2_GetOne(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()
	svr, err := server.NewGrpcServerV2()
	require.Nil(t, err)
	go svr.Start()
	defer svr.Stop()

	m := TestModel{
		Name: "Test Name",
	}
	modelSvc := service.NewModelServiceV2[TestModel]()
	id, err := modelSvc.InsertOne(m)
	require.Nil(t, err)
	m.SetId(id)
	time.Sleep(100 * time.Millisecond)

	c, err := grpc.Dial("localhost:9666", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.Nil(t, err)
	c.Connect()

	clientSvc := client.NewModelServiceV2[TestModel]()
	res, err := clientSvc.GetOne(bson.M{"name": m.Name}, nil)
	require.Nil(t, err)
	assert.Equal(t, res.Id, m.Id)
	assert.Equal(t, res.Name, m.Name)
}

func TestModelServiceV2_GetMany(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()
	svr, err := server.NewGrpcServerV2()
	require.Nil(t, err)
	go svr.Start()
	defer svr.Stop()

	m := TestModel{
		Name: "Test Name",
	}
	modelSvc := service.NewModelServiceV2[TestModel]()
	id, err := modelSvc.InsertOne(m)
	require.Nil(t, err)
	m.SetId(id)
	time.Sleep(100 * time.Millisecond)

	c, err := grpc.Dial("localhost:9666", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.Nil(t, err)
	c.Connect()

	clientSvc := client.NewModelServiceV2[TestModel]()
	res, err := clientSvc.GetMany(bson.M{"name": m.Name}, nil)
	require.Nil(t, err)
	assert.Equal(t, len(res), 1)
	assert.Equal(t, res[0].Id, m.Id)
	assert.Equal(t, res[0].Name, m.Name)
}

func TestModelServiceV2_DeleteById(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()
	svr, err := server.NewGrpcServerV2()
	require.Nil(t, err)
	go svr.Start()
	defer svr.Stop()

	m := TestModel{
		Name: "Test Name",
	}
	modelSvc := service.NewModelServiceV2[TestModel]()
	id, err := modelSvc.InsertOne(m)
	require.Nil(t, err)
	m.SetId(id)
	time.Sleep(100 * time.Millisecond)

	c, err := grpc.Dial("localhost:9666", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.Nil(t, err)
	c.Connect()

	clientSvc := client.NewModelServiceV2[TestModel]()
	err = clientSvc.DeleteById(m.Id)
	require.Nil(t, err)

	res, err := clientSvc.GetById(m.Id)
	assert.NotNil(t, err)
	require.Nil(t, res)
}

func TestModelServiceV2_DeleteOne(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()
	svr, err := server.NewGrpcServerV2()
	require.Nil(t, err)
	go svr.Start()
	defer svr.Stop()

	m := TestModel{
		Name: "Test Name",
	}
	modelSvc := service.NewModelServiceV2[TestModel]()
	id, err := modelSvc.InsertOne(m)
	require.Nil(t, err)
	m.SetId(id)
	time.Sleep(100 * time.Millisecond)

	c, err := grpc.Dial("localhost:9666", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.Nil(t, err)
	c.Connect()

	clientSvc := client.NewModelServiceV2[TestModel]()
	err = clientSvc.DeleteOne(bson.M{"name": m.Name})
	require.Nil(t, err)

	res, err := clientSvc.GetOne(bson.M{"name": m.Name}, nil)
	assert.NotNil(t, err)
	require.Nil(t, res)
}

func TestModelServiceV2_DeleteMany(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()
	svr, err := server.NewGrpcServerV2()
	require.Nil(t, err)
	go svr.Start()
	defer svr.Stop()

	m := TestModel{
		Name: "Test Name",
	}
	modelSvc := service.NewModelServiceV2[TestModel]()
	id, err := modelSvc.InsertOne(m)
	require.Nil(t, err)
	m.SetId(id)
	time.Sleep(100 * time.Millisecond)

	c, err := grpc.Dial("localhost:9666", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.Nil(t, err)
	c.Connect()

	clientSvc := client.NewModelServiceV2[TestModel]()
	err = clientSvc.DeleteMany(bson.M{"name": m.Name})
	require.Nil(t, err)

	res, err := clientSvc.GetMany(bson.M{"name": m.Name}, nil)
	require.Nil(t, err)
	assert.Equal(t, len(res), 0)
}

func TestModelServiceV2_UpdateById(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()
	svr, err := server.NewGrpcServerV2()
	require.Nil(t, err)
	go svr.Start()
	defer svr.Stop()

	m := TestModel{
		Name: "Test Name",
	}
	modelSvc := service.NewModelServiceV2[TestModel]()
	id, err := modelSvc.InsertOne(m)
	require.Nil(t, err)
	m.SetId(id)
	time.Sleep(100 * time.Millisecond)

	c, err := grpc.Dial("localhost:9666", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.Nil(t, err)
	c.Connect()

	clientSvc := client.NewModelServiceV2[TestModel]()
	err = clientSvc.UpdateById(m.Id, bson.M{"$set": bson.M{"name": "New Name"}})
	require.Nil(t, err)

	res, err := clientSvc.GetById(m.Id)
	require.Nil(t, err)
	assert.Equal(t, res.Name, "New Name")
}

func TestModelServiceV2_UpdateOne(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()
	svr, err := server.NewGrpcServerV2()
	require.Nil(t, err)
	go svr.Start()
	defer svr.Stop()

	m := TestModel{
		Name: "Test Name",
	}
	modelSvc := service.NewModelServiceV2[TestModel]()
	id, err := modelSvc.InsertOne(m)
	require.Nil(t, err)
	m.SetId(id)
	time.Sleep(100 * time.Millisecond)

	c, err := grpc.Dial("localhost:9666", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.Nil(t, err)
	c.Connect()

	clientSvc := client.NewModelServiceV2[TestModel]()
	err = clientSvc.UpdateOne(bson.M{"name": m.Name}, bson.M{"$set": bson.M{"name": "New Name"}})
	require.Nil(t, err)

	res, err := clientSvc.GetOne(bson.M{"name": "New Name"}, nil)
	require.Nil(t, err)
	assert.Equal(t, res.Name, "New Name")
}

func TestModelServiceV2_UpdateMany(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()
	svr, err := server.NewGrpcServerV2()
	require.Nil(t, err)
	go svr.Start()
	defer svr.Stop()

	m1 := TestModel{
		Name: "Test Name",
	}
	m2 := TestModel{
		Name: "Test Name",
	}
	modelSvc := service.NewModelServiceV2[TestModel]()
	_, err = modelSvc.InsertOne(m1)
	require.Nil(t, err)
	_, err = modelSvc.InsertOne(m2)
	require.Nil(t, err)
	time.Sleep(100 * time.Millisecond)

	c, err := grpc.Dial("localhost:9666", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.Nil(t, err)
	c.Connect()

	clientSvc := client.NewModelServiceV2[TestModel]()
	err = clientSvc.UpdateMany(bson.M{"name": "Test Name"}, bson.M{"$set": bson.M{"name": "New Name"}})
	require.Nil(t, err)

	res, err := clientSvc.GetMany(bson.M{"name": "New Name"}, nil)
	require.Nil(t, err)
	assert.Equal(t, len(res), 2)
}

func TestModelServiceV2_ReplaceById(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()
	svr, err := server.NewGrpcServerV2()
	require.Nil(t, err)
	go svr.Start()
	defer svr.Stop()

	m := TestModel{
		Name: "Test Name",
	}
	modelSvc := service.NewModelServiceV2[TestModel]()
	id, err := modelSvc.InsertOne(m)
	require.Nil(t, err)
	m.SetId(id)
	time.Sleep(100 * time.Millisecond)

	c, err := grpc.Dial("localhost:9666", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.Nil(t, err)
	c.Connect()

	clientSvc := client.NewModelServiceV2[TestModel]()
	m.Name = "New Name"
	err = clientSvc.ReplaceById(m.Id, m)
	require.Nil(t, err)

	res, err := clientSvc.GetById(m.Id)
	require.Nil(t, err)
	assert.Equal(t, res.Name, "New Name")
}

func TestModelServiceV2_ReplaceOne(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()
	svr, err := server.NewGrpcServerV2()
	require.Nil(t, err)
	go svr.Start()
	defer svr.Stop()

	m := TestModel{
		Name: "Test Name",
	}
	modelSvc := service.NewModelServiceV2[TestModel]()
	id, err := modelSvc.InsertOne(m)
	require.Nil(t, err)
	m.SetId(id)
	time.Sleep(100 * time.Millisecond)

	c, err := grpc.Dial("localhost:9666", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.Nil(t, err)
	c.Connect()

	clientSvc := client.NewModelServiceV2[TestModel]()
	m.Name = "New Name"
	err = clientSvc.ReplaceOne(bson.M{"name": "Test Name"}, m)
	require.Nil(t, err)

	res, err := clientSvc.GetOne(bson.M{"name": "New Name"}, nil)
	require.Nil(t, err)
	assert.Equal(t, res.Name, "New Name")
}

func TestModelServiceV2_InsertOne(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()
	svr, err := server.NewGrpcServerV2()
	require.Nil(t, err)
	go svr.Start()
	defer svr.Stop()

	c, err := grpc.Dial("localhost:9666", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.Nil(t, err)
	c.Connect()

	clientSvc := client.NewModelServiceV2[TestModel]()
	m := TestModel{
		Name: "Test Name",
	}
	id, err := clientSvc.InsertOne(m)
	require.Nil(t, err)

	res, err := clientSvc.GetById(id)
	require.Nil(t, err)
	assert.Equal(t, res.Name, m.Name)
}

func TestModelServiceV2_InsertMany(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()
	svr, err := server.NewGrpcServerV2()
	require.Nil(t, err)
	go svr.Start()
	defer svr.Stop()

	c, err := grpc.Dial("localhost:9666", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.Nil(t, err)
	c.Connect()

	clientSvc := client.NewModelServiceV2[TestModel]()
	models := []TestModel{
		{Name: "Test Name 1"},
		{Name: "Test Name 2"},
	}
	ids, err := clientSvc.InsertMany(models)
	require.Nil(t, err)

	for i, id := range ids {
		res, err := clientSvc.GetById(id)
		require.Nil(t, err)
		assert.Equal(t, res.Name, models[i].Name)
	}
}

func TestModelServiceV2_Count(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()
	svr, err := server.NewGrpcServerV2()
	require.Nil(t, err)
	go svr.Start()
	defer svr.Stop()

	modelSvc := service.NewModelServiceV2[TestModel]()
	for i := 0; i < 5; i++ {
		_, err = modelSvc.InsertOne(TestModel{
			Name: "Test Name",
		})
		require.Nil(t, err)
	}
	time.Sleep(100 * time.Millisecond)

	c, err := grpc.Dial("localhost:9666", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.Nil(t, err)
	c.Connect()

	clientSvc := client.NewModelServiceV2[TestModel]()
	count, err := clientSvc.Count(bson.M{})
	require.Nil(t, err)

	assert.Equal(t, count, 5)
}
