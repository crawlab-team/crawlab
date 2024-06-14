package test

import (
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/db/mongo"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"testing"
	"time"
)

func init() {
	viper.Set("mongo.db", "crawlab_test")
}

func TestExportController_Csv(t *testing.T) {
	T.Setup(t)
	e := T.NewExpect(t)

	// mongo collection
	colName := "test_collection_for_export"
	col := mongo.GetMongoCol(colName)

	// insert test data to mongo collection
	for i := 0; i < 10; i++ {
		_, err := col.Insert(bson.M{
			"field1": i + 1,
			"field2": i + 2,
			"field3": i + 3,
			"field4": i + 4,
		})
		require.Nil(t, err)
	}

	// export from mongo collection
	res := T.WithAuth(e.POST("/export/csv")).
		WithQuery("target", colName).
		Expect().Status(http.StatusOK).JSON().Object()
	res.Path("$.data").NotNull()

	// export id
	exportId := res.Path("$.data").String().Raw()

	// poll export with export id
	for i := 0; i < 10; i++ {
		res = T.WithAuth(e.GET("/export/csv/" + exportId)).Expect().Status(http.StatusOK).JSON().Object()
		status := res.Path("$.data.status").String().Raw()
		if status == constants.TaskStatusFinished {
			break
		}
		time.Sleep(1 * time.Second)
	}

	// download exported csv file
	csvFileBody := T.WithAuth(e.GET("/export/csv/" + exportId + "/download")).Expect().Status(http.StatusOK).Body()
	csvFileBody.NotEmpty()
}
