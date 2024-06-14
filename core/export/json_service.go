package export

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/ReneKroon/ttlcache"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/entity"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/utils"
	"github.com/crawlab-team/crawlab/db/mongo"
	"github.com/crawlab-team/crawlab/trace"
	"github.com/hashicorp/go-uuid"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"os"
	"path"
	"time"
)

type JsonService struct {
	cache *ttlcache.Cache
}

func (svc *JsonService) GenerateId() (exportId string, err error) {
	exportId, err = uuid.GenerateUUID()
	if err != nil {
		return "", trace.TraceError(err)
	}
	return exportId, nil
}

func (svc *JsonService) Export(exportType, target string, filter interfaces.Filter) (exportId string, err error) {
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

func (svc *JsonService) GetExport(exportId string) (export interfaces.Export, err error) {
	// get export from cache
	res, ok := svc.cache.Get(exportId)
	if !ok {
		return nil, trace.TraceError(errors.New("export not found"))
	}
	export = res.(interfaces.Export)
	return export, nil
}

func (svc *JsonService) export(export *entity.Export) {
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

	// data
	var jsonData []interface{}

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
			break
		}

		// convert raw data to entity
		var data map[string]interface{}
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

		jsonData = append(jsonData, data)
	}

	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		// error
		export.Status = constants.TaskStatusError
		export.EndTs = time.Now()
		log.Errorf("export error (id: %s): %v", export.Id, err)
		trace.PrintError(err)
		svc.cache.Set(export.Id, export)
		return
	}
	jsonString := string(jsonBytes)
	f := utils.OpenFile(export.DownloadPath)
	_, err = f.WriteString(jsonString)
	if err != nil {
		// error
		export.Status = constants.TaskStatusError
		export.EndTs = time.Now()
		log.Errorf("export error (id: %s): %v", export.Id, err)
		trace.PrintError(err)
		svc.cache.Set(export.Id, export)
		return
	}
}

func (svc *JsonService) getExportDir() (dir string, err error) {
	tempDir := os.TempDir()
	exportDir := path.Join(tempDir, "export", "json")
	if !utils.Exists(exportDir) {
		err := os.MkdirAll(exportDir, 0755)
		if err != nil {
			return "", err
		}
	}
	return exportDir, nil
}

func (svc *JsonService) getFileName(exportId string) (fileName string) {
	return exportId + "_" + time.Now().Format("20060102150405") + ".json"
}

// getDownloadPath returns the download path for the export
// format: <tempDir>/export/<exportId>/<exportId>_<timestamp>.csv
func (svc *JsonService) getDownloadPath(exportId string) (downloadPath string) {
	exportDir, err := svc.getExportDir()
	if err != nil {
		return ""
	}
	downloadPath = path.Join(exportDir, svc.getFileName(exportId))
	return downloadPath
}

func NewJsonService() (svc2 interfaces.ExportService) {
	cache := ttlcache.NewCache()
	cache.SetTTL(time.Minute * 5)
	svc := &JsonService{
		cache: cache,
	}
	return svc
}

var _jsonService interfaces.ExportService

func GetJsonService() (svc interfaces.ExportService) {
	if _jsonService == nil {
		_jsonService = NewJsonService()
	}
	return _jsonService
}
