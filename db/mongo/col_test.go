package mongo

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"strconv"
	"testing"
)

type ColTestObject struct {
	dbName  string
	colName string
	col     *Col
}

type TestDocument struct {
	Key   string   `bson:"key"`
	Value int      `bson:"value"`
	Tags  []string `bson:"tags"`
}

type TestAggregateResult struct {
	Id    string `bson:"_id"`
	Count int    `bson:"count"`
	Value int    `bson:"value"`
}

func setupColTest() (to *ColTestObject, err error) {
	dbName := "test_db"
	colName := "test_col"
	viper.Set("mongo.db", dbName)
	col := GetMongoCol(colName)
	if err := col.db.Drop(col.ctx); err != nil {
		return nil, err
	}
	return &ColTestObject{
		dbName:  dbName,
		colName: colName,
		col:     col,
	}, nil
}

func cleanupColTest(to *ColTestObject) {
	_ = to.col.db.Drop(to.col.ctx)
}

func TestGetMongoCol(t *testing.T) {
	colName := "test_col"

	col := GetMongoCol(colName)
	require.Equal(t, colName, col.c.Name())
}

func TestGetMongoColWithDb(t *testing.T) {
	dbName := "test_db"
	colName := "test_col"

	col := GetMongoColWithDb(colName, dbName)
	require.Equal(t, colName, col.c.Name())
	require.Equal(t, dbName, col.db.Name())
}

func TestCol_Insert(t *testing.T) {
	to, err := setupColTest()
	require.Nil(t, err)

	id, err := to.col.Insert(bson.M{"key": "value"})
	require.Nil(t, err)
	require.IsType(t, primitive.ObjectID{}, id)

	var doc map[string]string
	err = to.col.FindId(id).One(&doc)
	require.Nil(t, err)
	require.Equal(t, doc["key"], "value")

	cleanupColTest(to)
}

func TestCol_InsertMany(t *testing.T) {
	to, err := setupColTest()
	require.Nil(t, err)

	n := 10
	var docs []interface{}
	for i := 0; i < n; i++ {
		docs = append(docs, bson.M{"key": fmt.Sprintf("value-%d", i)})
	}
	ids, err := to.col.InsertMany(docs)
	require.Nil(t, err)
	require.Equal(t, n, len(ids))

	var resDocs []map[string]string
	err = to.col.Find(nil, &FindOptions{Sort: bson.D{{"_id", 1}}}).All(&resDocs)
	require.Nil(t, err)
	require.Equal(t, n, len(resDocs))
	for i, doc := range resDocs {
		require.Equal(t, fmt.Sprintf("value-%d", i), doc["key"])
	}

	cleanupColTest(to)
}

func TestCol_UpdateId(t *testing.T) {
	to, err := setupColTest()
	require.Nil(t, err)

	id, err := to.col.Insert(bson.M{"key": "old-value"})
	require.Nil(t, err)

	err = to.col.UpdateId(id, bson.M{
		"$set": bson.M{
			"key": "new-value",
		},
	})
	require.Nil(t, err)

	var doc map[string]string
	err = to.col.FindId(id).One(&doc)
	require.Nil(t, err)
	require.Equal(t, "new-value", doc["key"])

	cleanupColTest(to)
}

func TestCol_Update(t *testing.T) {
	to, err := setupColTest()
	require.Nil(t, err)

	n := 10
	var docs []interface{}
	for i := 0; i < n; i++ {
		docs = append(docs, bson.M{"key": fmt.Sprintf("old-value-%d", i)})
	}

	err = to.col.Update(nil, bson.M{
		"$set": bson.M{
			"key": "new-value",
		},
	})
	require.Nil(t, err)

	var resDocs []map[string]string
	err = to.col.Find(nil, &FindOptions{Sort: bson.D{{"_id", 1}}}).All(&resDocs)
	require.Nil(t, err)
	for _, doc := range resDocs {
		require.Equal(t, "new-value", doc["key"])
	}

	cleanupColTest(to)
}

func TestCol_ReplaceId(t *testing.T) {
	to, err := setupColTest()
	require.Nil(t, err)

	id, err := to.col.Insert(bson.M{"key": "old-value"})
	require.Nil(t, err)

	var doc map[string]interface{}
	err = to.col.FindId(id).One(&doc)
	require.Nil(t, err)
	doc["key"] = "new-value"

	err = to.col.ReplaceId(id, doc)
	require.Nil(t, err)

	err = to.col.FindId(id).One(doc)
	require.Nil(t, err)
	require.Equal(t, "new-value", doc["key"])

	cleanupColTest(to)
}

func TestCol_Replace(t *testing.T) {
	to, err := setupColTest()
	require.Nil(t, err)

	id, err := to.col.Insert(bson.M{"key": "old-value"})
	require.Nil(t, err)

	var doc map[string]interface{}
	err = to.col.FindId(id).One(&doc)
	require.Nil(t, err)
	doc["key"] = "new-value"

	err = to.col.Replace(bson.M{"key": "old-value"}, doc)
	require.Nil(t, err)

	err = to.col.FindId(id).One(&doc)
	require.Nil(t, err)
	require.Equal(t, "new-value", doc["key"])

	cleanupColTest(to)
}

func TestCol_DeleteId(t *testing.T) {
	to, err := setupColTest()
	require.Nil(t, err)

	id, err := to.col.Insert(bson.M{"key": "value"})
	require.Nil(t, err)

	err = to.col.DeleteId(id)
	require.Nil(t, err)

	total, err := to.col.Count(nil)
	require.Nil(t, err)
	require.Equal(t, 0, total)

	cleanupColTest(to)
}

func TestCol_Delete(t *testing.T) {
	to, err := setupColTest()
	require.Nil(t, err)

	n := 10
	var docs []interface{}
	for i := 0; i < n; i++ {
		docs = append(docs, bson.M{"key": fmt.Sprintf("value-%d", i)})
	}
	ids, err := to.col.InsertMany(docs)
	require.Nil(t, err)
	require.Equal(t, n, len(ids))

	err = to.col.Delete(bson.M{"key": "value-0"})
	require.Nil(t, err)

	total, err := to.col.Count(nil)
	require.Nil(t, err)
	require.Equal(t, n-1, total)

	err = to.col.Delete(nil)
	require.Nil(t, err)

	total, err = to.col.Count(nil)
	require.Nil(t, err)
	require.Equal(t, 0, total)

	cleanupColTest(to)
}

func TestCol_FindId(t *testing.T) {
	to, err := setupColTest()
	require.Nil(t, err)

	id, err := to.col.Insert(bson.M{"key": "value"})
	require.Nil(t, err)

	var doc map[string]string
	err = to.col.FindId(id).One(&doc)
	require.Nil(t, err)
	require.Equal(t, "value", doc["key"])

	cleanupColTest(to)
}

func TestCol_Find(t *testing.T) {
	to, err := setupColTest()
	require.Nil(t, err)

	n := 10
	var docs []interface{}
	for i := 0; i < n; i++ {
		docs = append(docs, TestDocument{
			Key:  fmt.Sprintf("value-%d", i),
			Tags: []string{"test tag"},
		})
	}
	ids, err := to.col.InsertMany(docs)
	require.Nil(t, err)
	require.Equal(t, n, len(ids))

	err = to.col.Find(nil, nil).All(&docs)
	require.Nil(t, err)
	require.Equal(t, n, len(docs))

	err = to.col.Find(bson.M{"key": bson.M{"$gte": fmt.Sprintf("value-%d", 5)}}, nil).All(&docs)
	require.Nil(t, err)
	require.Equal(t, n-5, len(docs))

	err = to.col.Find(nil, &FindOptions{
		Skip: 5,
	}).All(&docs)
	require.Nil(t, err)
	require.Equal(t, n-5, len(docs))

	err = to.col.Find(nil, &FindOptions{
		Limit: 5,
	}).All(&docs)
	require.Nil(t, err)
	require.Equal(t, 5, len(docs))

	var resDocs []TestDocument
	err = to.col.Find(nil, &FindOptions{
		Sort: bson.D{{"key", 1}},
	}).All(&resDocs)
	require.Nil(t, err)
	require.Greater(t, len(resDocs), 0)
	require.Equal(t, "value-0", resDocs[0].Key)

	err = to.col.Find(nil, &FindOptions{
		Sort: bson.D{{"key", -1}},
	}).All(&resDocs)
	require.Nil(t, err)
	require.Greater(t, len(resDocs), 0)
	require.Equal(t, fmt.Sprintf("value-%d", n-1), resDocs[0].Key)

	var resDocs2 []TestDocument
	err = to.col.Find(bson.M{"tags": bson.M{"$in": []string{"test tag"}}}, nil).All(&resDocs2)
	require.Nil(t, err)
	require.Greater(t, len(resDocs2), 0)

	cleanupColTest(to)
}

func TestCol_CreateIndex(t *testing.T) {
	to, err := setupColTest()
	require.Nil(t, err)

	err = to.col.CreateIndex(mongo.IndexModel{
		Keys: bson.D{{"key", 1}},
	})
	require.Nil(t, err)

	indexes, err := to.col.ListIndexes()
	require.Nil(t, err)
	require.Equal(t, 2, len(indexes))

	cleanupColTest(to)
}

func TestCol_Aggregate(t *testing.T) {
	to, err := setupColTest()
	require.Nil(t, err)

	n := 10
	v := 2
	var docs []interface{}
	for i := 0; i < n; i++ {
		docs = append(docs, TestDocument{
			Key:   fmt.Sprintf("%d", i%2),
			Value: v,
		})
	}
	ids, err := to.col.InsertMany(docs)
	require.Nil(t, err)
	require.Equal(t, n, len(ids))

	pipeline := mongo.Pipeline{
		{
			{
				"$group",
				bson.D{
					{"_id", "$key"},
					{
						"count",
						bson.D{{"$sum", 1}},
					},
					{
						"value",
						bson.D{{"$sum", "$value"}},
					},
				},
			},
		},
		{
			{
				"$sort",
				bson.D{{"_id", 1}},
			},
		},
	}
	var results []TestAggregateResult
	err = to.col.Aggregate(pipeline, nil).All(&results)
	require.Nil(t, err)
	require.Equal(t, 2, len(results))

	for i, r := range results {
		require.Equal(t, strconv.Itoa(i), r.Id)
		require.Equal(t, n/2, r.Count)
		require.Equal(t, n*v/2, r.Value)
	}
}

func TestCol_CreateIndexes(t *testing.T) {
	to, err := setupColTest()
	require.Nil(t, err)

	err = to.col.CreateIndexes([]mongo.IndexModel{
		{
			Keys: bson.D{{"key", 1}},
		},
		{
			Keys: bson.D{{"empty-key", 1}},
		},
	})
	require.Nil(t, err)

	indexes, err := to.col.ListIndexes()
	require.Nil(t, err)
	require.Equal(t, 3, len(indexes))

	cleanupColTest(to)
}

func TestCol_DeleteIndex(t *testing.T) {
	to, err := setupColTest()
	require.Nil(t, err)

	err = to.col.CreateIndex(mongo.IndexModel{
		Keys: bson.D{{"key", 1}},
	})
	require.Nil(t, err)

	indexes, err := to.col.ListIndexes()
	require.Nil(t, err)
	require.Equal(t, 2, len(indexes))
	for _, index := range indexes {
		name, ok := index["name"].(string)
		require.True(t, ok)

		if name != "_id_" {
			err = to.col.DeleteIndex(name)
			require.Nil(t, err)
			break
		}
	}

	indexes, err = to.col.ListIndexes()
	require.Nil(t, err)
	require.Equal(t, 1, len(indexes))

	cleanupColTest(to)
}

func TestCol_DeleteIndexes(t *testing.T) {
	to, err := setupColTest()
	require.Nil(t, err)

	err = to.col.CreateIndexes([]mongo.IndexModel{
		{
			Keys: bson.D{{"key", 1}},
		},
		{
			Keys: bson.D{{"empty-key", 1}},
		},
	})
	require.Nil(t, err)

	err = to.col.DeleteAllIndexes()
	require.Nil(t, err)

	indexes, err := to.col.ListIndexes()
	require.Nil(t, err)
	require.Equal(t, 1, len(indexes))

	cleanupColTest(to)
}
