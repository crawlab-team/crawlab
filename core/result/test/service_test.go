package test

import (
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestResultService_GetList(t *testing.T) {
	var err error
	T.Setup(t)

	n := 1000
	var docs []interface{}
	for i := 0; i < n; i++ {
		d := &models.Result{
			"i": i,
		}
		docs = append(docs, d)
	}
	_, err = T.TestCol.InsertMany(docs)
	require.Nil(t, err)

	// get all
	results, err := T.resultSvc.List(nil, nil)
	require.Nil(t, err)
	require.Equal(t, n, len(results))

	//query := bson.M{
	//	"i": bson.M{
	//		"$lt": n / 2,
	//	},
	//}
	//results, err = T.resultSvc.List(query, nil)
	//require.Nil(t, err)
	//require.Equal(t, n/2, len(results))
}

func TestResultService_Count(t *testing.T) {
	var err error
	T.Setup(t)

	n := 1000
	var docs []interface{}
	for i := 0; i < n; i++ {
		d := &models.Result{
			"i": i,
		}
		docs = append(docs, d)
	}
	_, err = T.TestCol.InsertMany(docs)
	require.Nil(t, err)

	// get all
	total, err := T.resultSvc.Count(nil)
	require.Nil(t, err)
	require.Equal(t, n, total)

	//query := bson.M{
	//	"i": bson.M{
	//		"$lt": n / 2,
	//	},
	//}
	//total, err = T.resultSvc.Count(query)
	//require.Nil(t, err)
	//require.Equal(t, n/2, total)
}
