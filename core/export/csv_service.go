package export

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/ReneKroon/ttlcache"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/entity"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/utils"
	"github.com/crawlab-team/crawlab/db/mongo"
	"github.com/crawlab-team/crawlab/trace"
	"github.com/hashicorp/go-uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"os"
	"path"
	"sort"
	"strconv"
	"time"
)

type CsvService struct {
	cache *ttlcache.Cache
}

func (svc *CsvService) GenerateId() (exportId string, err error) {
	exportId, err = uuid.GenerateUUID()
	if err != nil {
		return "", trace.TraceError(err)
	}
	return exportId, nil
}

func (svc *CsvService) Export(exportType, target string, filter interfaces.Filter) (exportId string, err error) {
	// generate export id
	exportId, err = svc.GenerateId()
	if err != nil {
		return "", err
	}

	// export
	export := &entity.Export{
		Id:           exportId,
		Type:         exportType,
		Target:       target,
		Filter:       filter,
		Status:       constants.TaskStatusRunning,
		StartTs:      time.Now(),
		FileName:     svc.getFileName(exportId),
		DownloadPath: svc.getDownloadPath(exportId),
		Limit:        100,
	}

	// save to cache
	svc.cache.Set(exportId, export)

	// execute export
	go svc.export(export)

	return exportId, nil
}

func (svc *CsvService) GetExport(exportId string) (export interfaces.Export, err error) {
	// get export from cache
	res, ok := svc.cache.Get(exportId)
	if !ok {
		return nil, trace.TraceError(errors.New("export not found"))
	}
	export = res.(interfaces.Export)
	return export, nil
}

func (svc *CsvService) export(export *entity.Export) {
	// check empty
	if export.Target == "" {
		err := errors.New("empty target")
		export.Status = constants.TaskStatusError
		export.EndTs = time.Now()
		log.Errorf("export error (id: %s): %v", export.Id, err)
		trace.PrintError(err)
		svc.cache.Set(export.Id, export)
		return
	}

	// mongo collection
	col := mongo.GetMongoCol(export.Target)

	// mongo query
	query, err := utils.FilterToQuery(export.Filter)
	if err != nil {
		export.Status = constants.TaskStatusError
		export.EndTs = time.Now()
		log.Errorf("export error (id: %s): %v", export.Id, err)
		trace.PrintError(err)
		svc.cache.Set(export.Id, export)
		return
	}

	// mongo cursor
	cur := col.Find(query, nil).GetCursor()

	// csv writer
	csvWriter, csvFile, err := svc.getCsvWriter(export)
	defer func() {
		csvWriter.Flush()
		_ = csvFile.Close()
	}()
	if err != nil {
		export.Status = constants.TaskStatusError
		export.EndTs = time.Now()
		log.Errorf("export error (id: %s): %v", export.Id, err)
		trace.PrintError(err)
		svc.cache.Set(export.Id, export)
		return
	}

	// write bom
	bom := []byte{0xEF, 0xBB, 0xBF}
	_, err = csvFile.Write(bom)
	if err != nil {
		trace.PrintError(err)
		return
	}

	// write csv header row
	columns, err := svc.getColumns(query, export)
	err = csvWriter.Write(columns)
	if err != nil {
		export.Status = constants.TaskStatusError
		export.EndTs = time.Now()
		log.Errorf("export error (id: %s): %v", export.Id, err)
		trace.PrintError(err)
		svc.cache.Set(export.Id, export)
		return
	}
	csvWriter.Flush()

	// iterate cursor
	i := 0
	for {
		// increment counter
		i++

		// check error
		err := cur.Err()
		if err != nil {
			if err != mongo2.ErrNoDocuments {
				// error
				export.Status = constants.TaskStatusError
				export.EndTs = time.Now()
				log.Errorf("export error (id: %s): %v", export.Id, err)
				trace.PrintError(err)
			} else {
				// no more data
				export.Status = constants.TaskStatusFinished
				export.EndTs = time.Now()
				log.Infof("export finished (id: %s)", export.Id)
			}
			svc.cache.Set(export.Id, export)
			return
		}

		// has data
		if !cur.Next(context.Background()) {
			// no more data
			export.Status = constants.TaskStatusFinished
			export.EndTs = time.Now()
			log.Infof("export finished (id: %s)", export.Id)
			svc.cache.Set(export.Id, export)
			return
		}

		// convert raw data to entity
		var data bson.M
		err = cur.Decode(&data)
		if err != nil {
			// error
			export.Status = constants.TaskStatusError
			export.EndTs = time.Now()
			log.Errorf("export error (id: %s): %v", export.Id, err)
			trace.PrintError(err)
			svc.cache.Set(export.Id, export)
			return
		}

		// write csv row cells
		cells := svc.getRowCells(columns, data)
		err = csvWriter.Write(cells)
		if err != nil {
			// error
			export.Status = constants.TaskStatusError
			export.EndTs = time.Now()
			log.Errorf("export error (id: %s): %v", export.Id, err)
			trace.PrintError(err)
			svc.cache.Set(export.Id, export)
			return
		}

		// flush if limit reached
		if i >= export.Limit {
			csvWriter.Flush()
			i = 0
		}
	}
}

func (svc *CsvService) getExportDir() (dir string, err error) {
	tempDir := os.TempDir()
	exportDir := path.Join(tempDir, "export", "csv")
	if !utils.Exists(exportDir) {
		err := os.MkdirAll(exportDir, 0755)
		if err != nil {
			return "", err
		}
	}
	return exportDir, nil
}

func (svc *CsvService) getFileName(exportId string) (fileName string) {
	return exportId + "_" + time.Now().Format("20060102150405") + ".csv"
}

// getDownloadPath returns the download path for the export
// format: <tempDir>/export/<exportId>/<exportId>_<timestamp>.csv
func (svc *CsvService) getDownloadPath(exportId string) (downloadPath string) {
	exportDir, err := svc.getExportDir()
	if err != nil {
		return ""
	}
	downloadPath = path.Join(exportDir, svc.getFileName(exportId))
	return downloadPath
}

func (svc *CsvService) getCsvWriter(export *entity.Export) (csvWriter *csv.Writer, csvFile *os.File, err error) {
	// open file
	csvFile, err = os.Create(export.DownloadPath)
	if err != nil {
		return nil, nil, trace.TraceError(err)
	}

	// create csv writer
	csvWriter = csv.NewWriter(csvFile)

	return csvWriter, csvFile, nil
}

func (svc *CsvService) getColumns(query bson.M, export interfaces.Export) (columns []string, err error) {
	// get mongo collection
	col := mongo.GetMongoCol(export.GetTarget())

	// get 10 records
	var data []bson.M
	if err := col.Find(query, &mongo.FindOptions{Limit: 10}).All(&data); err != nil {
		return nil, trace.TraceError(err)
	}

	// columns set
	columnsSet := make(map[string]bool)
	for _, d := range data {
		for k := range d {
			columnsSet[k] = true
		}
	}

	// columns
	columns = make([]string, 0, len(columnsSet))
	for k := range columnsSet {
		// skip task key
		if k == constants.TaskKey {
			continue
		}

		// skip _id
		if k == "_id" {
			continue
		}

		// append to columns
		columns = append(columns, k)
	}

	// order columns
	sort.Strings(columns)

	return columns, nil
}

func (svc *CsvService) getRowCells(columns []string, data bson.M) (cells []string) {
	for _, c := range columns {
		v, ok := data[c]
		if !ok {
			cells = append(cells, "")
			continue
		}
		switch v.(type) {
		case string:
			cells = append(cells, v.(string))
		case time.Time:
			cells = append(cells, v.(time.Time).Format("2006-01-02 15:04:05"))
		case int:
			cells = append(cells, strconv.Itoa(v.(int)))
		case int32:
			cells = append(cells, strconv.Itoa(int(v.(int32))))
		case int64:
			cells = append(cells, strconv.FormatInt(v.(int64), 10))
		case float32:
			cells = append(cells, strconv.FormatFloat(float64(v.(float32)), 'f', -1, 32))
		case float64:
			cells = append(cells, strconv.FormatFloat(v.(float64), 'f', -1, 64))
		case bool:
			cells = append(cells, strconv.FormatBool(v.(bool)))
		case primitive.ObjectID:
			cells = append(cells, v.(primitive.ObjectID).Hex())
		case primitive.DateTime:
			cells = append(cells, v.(primitive.DateTime).Time().Format("2006-01-02 15:04:05"))
		default:
			cells = append(cells, fmt.Sprintf("%v", v))
		}
	}
	return cells
}

func NewCsvService() (svc2 interfaces.ExportService) {
	cache := ttlcache.NewCache()
	cache.SetTTL(time.Minute * 5)
	svc := &CsvService{
		cache: cache,
	}
	return svc
}

var _csvService interfaces.ExportService

func GetCsvService() (svc interfaces.ExportService) {
	if _csvService == nil {
		_csvService = NewCsvService()
	}
	return _csvService
}
