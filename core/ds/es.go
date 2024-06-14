package ds

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/crawlab-team/crawlab/core/constants"
	constants2 "github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/entity"
	entity2 "github.com/crawlab-team/crawlab/core/entity"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/core/utils"
	"github.com/crawlab-team/crawlab/db/generic"
	"github.com/crawlab-team/crawlab/trace"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"sync"
	"time"
)

type ElasticsearchService struct {
	// dependencies
	modelSvc service.ModelService

	// internals
	dc *models.DataCollection // models.DataCollection
	ds *models.DataSource     // models.DataSource
	c  *elasticsearch.Client  // elasticsearch.Client
	t  time.Time
}

func (svc *ElasticsearchService) Insert(records ...interface{}) (err error) {
	// wait group
	var wg sync.WaitGroup
	wg.Add(len(records))

	// iterate records
	for _, r := range records {
		// async operation
		go func(r interface{}) {
			switch r.(type) {
			case entity.Result:
				// convert type to entity.Result
				d := r.(entity.Result)

				// get document id
				id := d.GetValue("id")
				var docId string
				switch id.(type) {
				case string:
					docId = id.(string)
				}
				if docId == "" {
					docId = uuid.New().String() // generate new uuid if id is empty
				}

				// collection
				d[constants2.DataCollectionKey] = svc.dc.Name

				// index request
				req := esapi.IndexRequest{
					Index:      svc.getIndexName(),
					DocumentID: docId,
					Body:       strings.NewReader(d.String()),
				}

				// perform request
				res, err := req.Do(context.Background(), svc.c)
				if err != nil {
					trace.PrintError(err)
					wg.Done()
					return
				}
				defer res.Body.Close()
				if res.IsError() {
					trace.PrintError(errors.New(fmt.Sprintf("[ElasticsearchService] [%s] error inserting record: %v", res.Status(), r)))
				}

				// release
				wg.Done()
			default:
				wg.Done()
				return
			}
		}(r)
	}

	// wait
	wg.Wait()

	return nil
}

func (svc *ElasticsearchService) List(query generic.ListQuery, opts *generic.ListOptions) (results []interface{}, err error) {
	data, err := svc.getListResponse(query, opts, false)
	if err != nil {
		return nil, err
	}
	for _, hit := range data.Hits.Hits {
		results = append(results, hit.Source)
	}
	return results, nil
}

func (svc *ElasticsearchService) Count(query generic.ListQuery) (n int, err error) {
	data, err := svc.getListResponse(query, nil, true)
	if err != nil {
		return n, err
	}
	return int(data.Hits.Total.Value), nil
}

func (svc *ElasticsearchService) getListResponse(query generic.ListQuery, opts *generic.ListOptions, trackTotalHits bool) (data *entity2.ElasticsearchResponseData, err error) {
	if opts == nil {
		opts = &generic.ListOptions{}
	}
	query = append(query, generic.ListQueryCondition{
		Key:   constants2.DataCollectionKey,
		Op:    constants2.FilterOpEqual,
		Value: svc.dc.Name,
	})
	res, err := svc.c.Search(
		svc.c.Search.WithContext(context.Background()),
		svc.c.Search.WithIndex(svc.getIndexName()),
		svc.c.Search.WithBody(utils.GetElasticsearchQueryWithOptions(query, opts)),
		svc.c.Search.WithTrackTotalHits(trackTotalHits),
	)
	if err != nil {
		return nil, trace.TraceError(err)
	}
	defer res.Body.Close()
	if res.IsError() {
		err = errors.New(fmt.Sprintf("[ElasticsearchService] [%s] error listing records: response=%s, query=%v opts=%v", res.Status(), res.String(), query, opts))
		trace.PrintError(err)
		return nil, err
	}
	data = &entity2.ElasticsearchResponseData{}
	if err := json.NewDecoder(res.Body).Decode(data); err != nil {
		return nil, trace.TraceError(err)
	}
	return data, nil
}

func (svc *ElasticsearchService) getIndexName() (index string) {
	if svc.ds.Database == "" {
		return svc.dc.Name
	} else {
		return svc.ds.Name
	}
}

func NewDataSourceElasticsearchService(colId primitive.ObjectID, dsId primitive.ObjectID) (svc2 interfaces.ResultService, err error) {
	// service
	svc := &ElasticsearchService{}

	// dependency injection
	svc.modelSvc, err = service.GetService()
	if err != nil {
		return nil, err
	}

	// data source
	if dsId.IsZero() {
		svc.ds = &models.DataSource{}
	} else {
		svc.ds, err = svc.modelSvc.GetDataSourceById(dsId)
		if err != nil {
			return nil, err
		}
	}

	// data source defaults
	if svc.ds.Host == "" {
		svc.ds.Host = constants.DefaultHost
	}
	if svc.ds.Port == "" {
		svc.ds.Port = constants.DefaultElasticsearchPort
	}

	// data source password
	pwd, err := svc.modelSvc.GetPasswordById(svc.ds.Id)
	if err == nil {
		svc.ds.Password, err = utils.DecryptAES(pwd.Password)
		if err != nil {
			return nil, err
		}
	}

	// data collection
	svc.dc, err = svc.modelSvc.GetDataCollectionById(colId)
	if err != nil {
		return nil, err
	}

	// client
	svc.c, err = utils.GetElasticsearchClient(svc.ds)
	if err != nil {
		return nil, err
	}

	return svc, nil
}

func (svc *ElasticsearchService) Index(fields []string) {
	// TODO: implement me
}

func (svc *ElasticsearchService) SetTime(t time.Time) {
	svc.t = t
}

func (svc *ElasticsearchService) GetTime() (t time.Time) {
	return svc.t
}
