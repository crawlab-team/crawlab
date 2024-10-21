package service_test

import (
	"context"
	"github.com/apex/log"
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
		log.Errorf("dropping test db error: %v", err)
		return
	}
	log.Infof("dropped test db")
}

func TestModelServiceV2(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()

	t.Run("GetById", func(t *testing.T) {
		svc := service.NewModelServiceV2[TestModel]()
		testModel := TestModel{Name: "GetById Test"}

		id, err := svc.InsertOne(testModel)
		require.Nil(t, err)
		time.Sleep(100 * time.Millisecond)

		result, err := svc.GetById(id)
		require.Nil(t, err)
		assert.Equal(t, testModel.Name, result.Name)
	})

	t.Run("GetOne", func(t *testing.T) {
		svc := service.NewModelServiceV2[TestModel]()
		testModel := TestModel{Name: "GetOne Test"}

		_, err := svc.InsertOne(testModel)
		require.Nil(t, err)
		time.Sleep(100 * time.Millisecond)

		result, err := svc.GetOne(bson.M{"name": "GetOne Test"}, nil)
		require.Nil(t, err)
		assert.Equal(t, testModel.Name, result.Name)
	})

	t.Run("GetMany", func(t *testing.T) {
		svc := service.NewModelServiceV2[TestModel]()
		testModels := []TestModel{
			{Name: "GetMany Test 1"},
			{Name: "GetMany Test 2"},
		}

		_, err := svc.InsertMany(testModels)
		require.Nil(t, err)
		time.Sleep(100 * time.Millisecond)

		results, err := svc.GetMany(bson.M{"name": bson.M{"$regex": "^GetMany Test"}}, nil)
		require.Nil(t, err)
		assert.Equal(t, 2, len(results))
	})

	t.Run("InsertOne", func(t *testing.T) {
		svc := service.NewModelServiceV2[TestModel]()
		testModel := TestModel{Name: "InsertOne Test"}

		id, err := svc.InsertOne(testModel)
		require.Nil(t, err)
		assert.NotEqual(t, primitive.NilObjectID, id)
	})

	t.Run("InsertMany", func(t *testing.T) {
		svc := service.NewModelServiceV2[TestModel]()
		testModels := []TestModel{
			{Name: "InsertMany Test 1"},
			{Name: "InsertMany Test 2"},
		}

		ids, err := svc.InsertMany(testModels)
		require.Nil(t, err)
		assert.Equal(t, 2, len(ids))
	})

	t.Run("UpdateById", func(t *testing.T) {
		svc := service.NewModelServiceV2[TestModel]()
		testModel := TestModel{Name: "UpdateById Test"}

		id, err := svc.InsertOne(testModel)
		require.Nil(t, err)
		time.Sleep(100 * time.Millisecond)

		update := bson.M{"$set": bson.M{"name": "UpdateById Test New Name"}}
		err = svc.UpdateById(id, update)
		require.Nil(t, err)
		time.Sleep(100 * time.Millisecond)

		result, err := svc.GetById(id)
		require.Nil(t, err)
		assert.Equal(t, "UpdateById Test New Name", result.Name)
	})

	t.Run("UpdateOne", func(t *testing.T) {
		svc := service.NewModelServiceV2[TestModel]()
		testModel := TestModel{Name: "UpdateOne Test"}

		_, err := svc.InsertOne(testModel)
		require.Nil(t, err)
		time.Sleep(100 * time.Millisecond)

		update := bson.M{"$set": bson.M{"name": "UpdateOne Test New Name"}}
		err = svc.UpdateOne(bson.M{"name": "UpdateOne Test"}, update)
		require.Nil(t, err)
		time.Sleep(100 * time.Millisecond)

		result, err := svc.GetOne(bson.M{"name": "UpdateOne Test New Name"}, nil)
		require.Nil(t, err)
		assert.Equal(t, "UpdateOne Test New Name", result.Name)
	})

	t.Run("UpdateMany", func(t *testing.T) {
		svc := service.NewModelServiceV2[TestModel]()
		testModels := []TestModel{
			{Name: "UpdateMany Test 1"},
			{Name: "UpdateMany Test 2"},
		}

		_, err := svc.InsertMany(testModels)
		require.Nil(t, err)
		time.Sleep(100 * time.Millisecond)

		update := bson.M{"$set": bson.M{"name": "UpdateMany Test New Name"}}
		err = svc.UpdateMany(bson.M{"name": bson.M{"$regex": "^UpdateMany Test"}}, update)
		require.Nil(t, err)
		time.Sleep(100 * time.Millisecond)

		results, err := svc.GetMany(bson.M{"name": "UpdateMany Test New Name"}, nil)
		require.Nil(t, err)
		assert.Equal(t, 2, len(results))
	})

	t.Run("DeleteById", func(t *testing.T) {
		svc := service.NewModelServiceV2[TestModel]()
		testModel := TestModel{Name: "DeleteById Test"}

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
		svc := service.NewModelServiceV2[TestModel]()
		testModel := TestModel{Name: "DeleteOne Test"}

		_, err := svc.InsertOne(testModel)
		require.Nil(t, err)
		time.Sleep(100 * time.Millisecond)

		err = svc.DeleteOne(bson.M{"name": "DeleteOne Test"})
		require.Nil(t, err)
		time.Sleep(100 * time.Millisecond)

		result, err := svc.GetOne(bson.M{"name": "DeleteOne Test"}, nil)
		assert.NotNil(t, err)
		assert.Nil(t, result)
	})

	t.Run("DeleteMany", func(t *testing.T) {
		svc := service.NewModelServiceV2[TestModel]()
		testModels := []TestModel{
			{Name: "DeleteMany Test 1"},
			{Name: "DeleteMany Test 2"},
		}

		_, err := svc.InsertMany(testModels)
		require.Nil(t, err)
		time.Sleep(100 * time.Millisecond)

		err = svc.DeleteMany(bson.M{"name": bson.M{"$regex": "^DeleteMany Test"}})
		require.Nil(t, err)
		time.Sleep(100 * time.Millisecond)

		results, err := svc.GetMany(bson.M{"name": bson.M{"$regex": "^DeleteMany Test"}}, nil)
		require.Nil(t, err)
		assert.Equal(t, 0, len(results))
	})

	t.Run("Count", func(t *testing.T) {
		svc := service.NewModelServiceV2[TestModel]()
		testModels := []TestModel{
			{Name: "Count Test 1"},
			{Name: "Count Test 2"},
		}

		_, err := svc.InsertMany(testModels)
		require.Nil(t, err)
		time.Sleep(100 * time.Millisecond)

		total, err := svc.Count(bson.M{"name": bson.M{"$regex": "^Count Test"}})
		require.Nil(t, err)
		assert.Equal(t, 2, total)
	})
}
