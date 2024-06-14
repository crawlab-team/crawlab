package export

import (
	"encoding/csv"
	"fmt"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/db/mongo"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestCsvService_Export(t *testing.T) {
	// test data rows
	var rows []interface{}
	for i := 0; i < 10; i++ {
		data := bson.M{
			"no":              i,
			"string_field":    "test",
			"int_field":       1,
			"float_field":     1.1,
			"bool_field":      true,
			"time_field":      time.Now(),
			"object_id_field": primitive.NewObjectID(),
		}
		rows = append(rows, data)
	}

	// test mongo collection name
	collectionName := "test_collection"

	// test mongo collection
	collection := mongo.GetMongoCol(collectionName)

	// delete records of test mongo collection after test
	t.Cleanup(func() {
		_ = collection.Delete(bson.M{})
	})

	// save test data rows to mongo collection
	_, err := collection.InsertMany(rows)
	require.Nil(t, err)

	// export service
	csvSvc := NewCsvService()

	// export
	exportId, err := csvSvc.Export(collectionName, collectionName, nil)
	require.Nil(t, err)

	// get export
	export, err := csvSvc.GetExport(exportId)
	require.Nil(t, err)
	require.NotNil(t, export)
	require.Equal(t, exportId, export.GetId())
	require.NotNil(t, export.GetDownloadPath())

	// wait for export to finish with timeout of 5 seconds
	timeout := time.After(5 * time.Second)
	finished := false
	for {
		if finished {
			break
		}
		select {
		case <-timeout:
			t.Fatal("export timeout")
		default:
			if export.GetStatus() == constants.TaskStatusFinished {
				finished = true
				continue
			}
			time.Sleep(100 * time.Millisecond)
		}
	}

	// export file path
	exportFilePath := export.GetDownloadPath()
	require.FileExists(t, exportFilePath)

	// csv file
	csvFile, err := os.Open(exportFilePath)
	require.Nil(t, err)
	defer csvFile.Close()

	// csv file reader
	csvFileReader := csv.NewReader(csvFile)

	// csv file rows
	csvFileRows, err := csvFileReader.ReadAll()
	require.Nil(t, err)
	require.Equal(t, len(rows), len(csvFileRows)-1)

	// csv file columns
	csvFileColumns := csvFileRows[0]

	// iterate csv file records and compare with test data rows
	for i, row := range rows {
		// csv file record
		csvFileRecord := csvFileRows[i+1]

		// iterate csv file columns and compare with test data row
		for j, column := range csvFileColumns {
			// csv file column value
			csvFileColumnValue := csvFileRecord[j]

			// compare csv file column value with test data row
			switch column {
			case "no":
				// convert int to string
				stringValue := fmt.Sprintf("%d", row.(bson.M)["no"].(int))
				require.Equal(t, stringValue, csvFileColumnValue)
			case "string_field":
				require.Equal(t, row.(bson.M)["string_field"], csvFileColumnValue)
			case "int_field":
				// convert int to string
				stringValue := fmt.Sprintf("%d", row.(bson.M)["int_field"])
				require.Equal(t, stringValue, csvFileColumnValue)
			case "float_field":
				// convert string to float
				floatValue, err := strconv.ParseFloat(csvFileColumnValue, 64)
				require.Nil(t, err)
				require.Equal(t, row.(bson.M)["float_field"].(float64), floatValue)
			case "bool_field":
				// convert bool to string
				stringValue := fmt.Sprintf("%t", row.(bson.M)["bool_field"])
				require.Equal(t, stringValue, csvFileColumnValue)
			case "time_field":
				// convert time to string
				stringValue := row.(bson.M)["time_field"].(time.Time).Format("2006-01-02 15:04:05")
				require.Equal(t, stringValue, csvFileColumnValue)
			}
		}
	}
}
