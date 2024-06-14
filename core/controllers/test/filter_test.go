package test

import (
	"encoding/json"
	"fmt"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/entity"
	"github.com/crawlab-team/crawlab/db/mongo"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"testing"
)

func init() {
	viper.Set("mongo.db", "crawlab_test")
}

func TestFilterController_GetColFieldOptions(t *testing.T) {
	T.Setup(t)
	e := T.NewExpect(t)

	// mongo collection
	colName := "test_collection_for_filter"
	field1 := "field1"
	field2 := "field2"
	value1 := "value1"
	col := mongo.GetMongoCol(colName)
	n := 10
	var ids []primitive.ObjectID
	var names []string
	for i := 0; i < n; i++ {
		_id := primitive.NewObjectID()
		ids = append(ids, _id)
		name := fmt.Sprintf("name_%d", i)
		names = append(names, name)
		_, err := col.Insert(bson.M{field1: value1, field2: i % 2, "name": name, "_id": _id})
		require.Nil(t, err)
	}

	// validate filter options field 1
	res := T.WithAuth(e.GET(fmt.Sprintf("/filters/%s/%s/%s", colName, field1, field1))).
		Expect().Status(http.StatusOK).JSON().Object()
	res.Path("$.data").NotNull()
	res.Path("$.data").Array().Length().Equal(1)
	res.Path("$.data").Array().Element(0).Path("$.value").Equal(value1)

	// validate filter options field 2
	res = T.WithAuth(e.GET(fmt.Sprintf("/filters/%s/%s/%s", colName, field2, field2))).
		Expect().Status(http.StatusOK).JSON().Object()
	res.Path("$.data").NotNull()
	res.Path("$.data").Array().Length().Equal(2)

	// validate filter options with query
	conditions := []entity.Condition{{field2, constants.FilterOpEqual, 0}}
	conditionsJson, err := json.Marshal(conditions)
	conditionsJsonStr := string(conditionsJson)
	require.Nil(t, err)
	res = T.WithAuth(e.GET(fmt.Sprintf("/filters/%s/%s/%s", colName, field2, field2))).
		WithQuery(constants.FilterQueryFieldConditions, conditionsJsonStr).
		Expect().Status(http.StatusOK).JSON().Object()
	res.Path("$.data").NotNull()
	res.Path("$.data").Array().Length().Equal(1)

	// validate filter options (basic path)
	res = T.WithAuth(e.GET(fmt.Sprintf("/filters/%s", colName))).
		Expect().Status(http.StatusOK).JSON().Object()
	res.Path("$.data").NotNull()
	res.Path("$.data").Array().Length().Equal(n)
	for i := 0; i < n; i++ {
		res.Path("$.data").Array().Element(i).Object().Value("value").Equal(ids[i])
		res.Path("$.data").Array().Element(i).Object().Value("label").Equal(names[i])
	}
}
