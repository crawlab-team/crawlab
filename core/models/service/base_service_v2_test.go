package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/crawlab-team/crawlab/core/models/models/v2"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/db/mongo"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TestModel struct {
	Id                            primitive.ObjectID `bson:"_id,omitempty" collection:"testmodels"`
	models.BaseModelV2[TestModel] `bson:",inline"`
	Name                          string `bson:"name"`
}

func setupTestDB() {
	viper.Set("mongo.db", "testdb")
}

func teardownTestDB() {
	db := mongo.GetMongoDb("testdb")
	err := db.Drop(context.Background())
	if err != nil {
		return
	}
}

func TestModelServiceV2(t *testing.T) {
	t.Run("GetById", func(t *testing.T) {
		setupTestDB()
		defer teardownTestDB()

		svc := service.NewModelServiceV2[TestModel]()
		testModel := TestModel{Name: "Test Name"}

		id, err := svc.InsertOne(testModel)
		require.Nil(t, err)
		time.Sleep(100 * time.Millisecond)

		result, err := svc.GetById(id)
		require.Nil(t, err)
		assert.Equal(t, testModel.Name, result.Name)
	})

	t.Run("GetOne", func(t *testing.T) {
		setupTestDB()
		defer teardownTestDB()

		svc := service.NewModelServiceV2[TestModel]()
		testModel := TestModel{Name: "Test Name"}

		_, err := svc.InsertOne(testModel)
		require.Nil(t, err)
		time.Sleep(100 * time.Millisecond)

		result, err := svc.GetOne(bson.M{"name": "Test Name"}, nil)
		require.Nil(t, err)
		assert.Equal(t, testModel.Name, result.Name)
	})

	t.Run("GetMany", func(t *testing.T) {
		setupTestDB()
		defer teardownTestDB()

		svc := service.NewModelServiceV2[TestModel]()
		testModels := []TestModel{
			{Name: "Name1"},
			{Name: "Name2"},
		}

		_, err := svc.InsertMany(testModels)
		require.Nil(t, err)
		time.Sleep(100 * time.Millisecond)

		results, err := svc.GetMany(bson.M{}, nil)
		require.Nil(t, err)
		assert.Equal(t, 2, len(results))
	})

	t.Run("InsertOne", func(t *testing.T) {
		setupTestDB()
		defer teardownTestDB()

		svc := service.NewModelServiceV2[TestModel]()
		testModel := TestModel{Name: "Test Name"}

		id, err := svc.InsertOne(testModel)
		require.Nil(t, err)
		assert.NotEqual(t, primitive.NilObjectID, id)
	})

	t.Run("InsertMany", func(t *testing.T) {
		setupTestDB()
		defer teardownTestDB()

		svc := service.NewModelServiceV2[TestModel]()
		testModels := []TestModel{
			{Name: "Name1"},
			{Name: "Name2"},
		}

		ids, err := svc.InsertMany(testModels)
		require.Nil(t, err)
		assert.Equal(t, 2, len(ids))
	})

	t.Run("UpdateById", func(t *testing.T) {
		setupTestDB()
		defer teardownTestDB()

		svc := service.NewModelServiceV2[TestModel]()
		testModel := TestModel{Name: "Old Name"}

		id, err := svc.InsertOne(testModel)
		require.Nil(t, err)
		time.Sleep(100 * time.Millisecond)

		update := bson.M{"$set": bson.M{"name": "New Name"}}
		err = svc.UpdateById(id, update)
		require.Nil(t, err)
		time.Sleep(100 * time.Millisecond)

		result, err := svc.GetById(id)
		require.Nil(t, err)
		assert.Equal(t, "New Name", result.Name)
	})

	t.Run("UpdateOne", func(t *testing.T) {
		setupTestDB()
		defer teardownTestDB()

		svc := service.NewModelServiceV2[TestModel]()
		testModel := TestModel{Name: "Old Name"}

		_, err := svc.InsertOne(testModel)
		require.Nil(t, err)
		time.Sleep(100 * time.Millisecond)

		update := bson.M{"$set": bson.M{"name": "New Name"}}
		err = svc.UpdateOne(bson.M{"name": "Old Name"}, update)
		require.Nil(t, err)
		time.Sleep(100 * time.Millisecond)

		result, err := svc.GetOne(bson.M{"name": "New Name"}, nil)
		require.Nil(t, err)
		assert.Equal(t, "New Name", result.Name)
	})

	t.Run("UpdateMany", func(t *testing.T) {
		setupTestDB()
		defer teardownTestDB()

		svc := service.NewModelServiceV2[TestModel]()
		testModels := []TestModel{
			{Name: "Old Name1"},
			{Name: "Old Name2"},
		}

		_, err := svc.InsertMany(testModels)
		require.Nil(t, err)
		time.Sleep(100 * time.Millisecond)

		update := bson.M{"$set": bson.M{"name": "New Name"}}
		err = svc.UpdateMany(bson.M{"name": bson.M{"$regex": "^Old"}}, update)
		require.Nil(t, err)
		time.Sleep(100 * time.Millisecond)

		results, err := svc.GetMany(bson.M{"name": "New Name"}, nil)
		require.Nil(t, err)
		assert.Equal(t, 2, len(results))
	})

	t.Run("DeleteById", func(t *testing.T) {
		setupTestDB()
		defer teardownTestDB()

		svc := service.NewModelServiceV2[TestModel]()
		testModel := TestModel{Name: "Test Name"}

		id, err := svc.InsertOne(testModel)
		require.Nil(t, err)
		time.Sleep(100 * time.Millisecond)

		err = svc.DeleteById(id)
		require.Nil(t, err)
		time.Sleep(100 * time.Millisecond)

		result, err := svc.GetById(id)
		assert.NotNil(t, err)
		assert.Nil(t, result)
	})

	t.Run("DeleteOne", func(t *testing.T) {
		setupTestDB()
		defer teardownTestDB()

		svc := service.NewModelServiceV2[TestModel]()
		testModel := TestModel{Name: "Test Name"}

		_, err := svc.InsertOne(testModel)
		require.Nil(t, err)
		time.Sleep(100 * time.Millisecond)

		err = svc.DeleteOne(bson.M{"name": "Test Name"})
		require.Nil(t, err)
		time.Sleep(100 * time.Millisecond)

		result, err := svc.GetOne(bson.M{"name": "Test Name"}, nil)
		assert.NotNil(t, err)
		assert.Nil(t, result)
	})

	t.Run("DeleteMany", func(t *testing.T) {
		setupTestDB()
		defer teardownTestDB()

		svc := service.NewModelServiceV2[TestModel]()
		testModels := []TestModel{
			{Name: "Test Name1"},
			{Name: "Test Name2"},
		}

		_, err := svc.InsertMany(testModels)
		require.Nil(t, err)
		time.Sleep(100 * time.Millisecond)

		err = svc.DeleteMany(bson.M{"name": bson.M{"$regex": "^Test Name"}})
		require.Nil(t, err)
		time.Sleep(100 * time.Millisecond)

		results, err := svc.GetMany(bson.M{"name": bson.M{"$regex": "^Test Name"}}, nil)
		require.Nil(t, err)
		assert.Equal(t, 0, len(results))
	})

	t.Run("Count", func(t *testing.T) {
		setupTestDB()
		defer teardownTestDB()

		svc := service.NewModelServiceV2[TestModel]()
		testModels := []TestModel{
			{Name: "Name1"},
			{Name: "Name2"},
		}

		_, err := svc.InsertMany(testModels)
		require.Nil(t, err)
		time.Sleep(100 * time.Millisecond)

		total, err := svc.Count(bson.M{})
		require.Nil(t, err)
		assert.Equal(t, 2, total)
	})
}
