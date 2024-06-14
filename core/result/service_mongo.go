package result

import (
	"github.com/crawlab-team/crawlab/trace"
	"time"

	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/core/utils"
	"github.com/crawlab-team/crawlab/db/generic"
	"github.com/crawlab-team/crawlab/db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ServiceMongo struct {
	// dependencies
	modelSvc    service.ModelService
	modelColSvc interfaces.ModelBaseService

	// internals
	colId primitive.ObjectID     // _id of models.DataCollection
	dc    *models.DataCollection // models.DataCollection
	t     time.Time
}

func (svc *ServiceMongo) List(query generic.ListQuery, opts *generic.ListOptions) (results []interface{}, err error) {
	_query := svc.getQuery(query)
	_opts := svc.getOpts(opts)
	return svc.getList(_query, _opts)
}

func (svc *ServiceMongo) Count(query generic.ListQuery) (n int, err error) {
	_query := svc.getQuery(query)
	return svc.modelColSvc.Count(_query)
}

func (svc *ServiceMongo) Insert(docs ...interface{}) (err error) {
	if svc.dc.Dedup.Enabled && len(svc.dc.Dedup.Keys) > 0 {
		for _, doc := range docs {
			hash, err := utils.GetResultHash(doc, svc.dc.Dedup.Keys)
			if err != nil {
				return err
			}
			doc.(interfaces.Result).SetValue(constants.HashKey, hash)
			query := bson.M{constants.HashKey: hash}
			switch svc.dc.Dedup.Type {
			case constants.DedupTypeOverwrite:
				err = mongo.GetMongoCol(svc.dc.Name).ReplaceWithOptions(query, doc, &options.ReplaceOptions{Upsert: &[]bool{true}[0]})
				if err != nil {
					return trace.TraceError(err)
				}
			default:
				var o bson.M
				err := mongo.GetMongoCol(svc.dc.Name).Find(query, &mongo.FindOptions{Limit: 1}).One(&o)
				if err == nil {
					// exists, ignore
					continue
				}
				if err != mongo2.ErrNoDocuments {
					// error
					return trace.TraceError(err)
				}
				// not exists, insert
				_, err = mongo.GetMongoCol(svc.dc.Name).Insert(doc)
				if err != nil {
					return trace.TraceError(err)
				}
			}
		}
	} else {
		_, err = mongo.GetMongoCol(svc.dc.Name).InsertMany(docs)
		if err != nil {
			return trace.TraceError(err)
		}
	}
	return nil
}

func (svc *ServiceMongo) Index(fields []string) {
	for _, field := range fields {
		_ = mongo.GetMongoCol(svc.dc.Name).CreateIndex(mongo2.IndexModel{Keys: bson.M{field: 1}})
	}
}

func (svc *ServiceMongo) SetTime(t time.Time) {
	svc.t = t
}

func (svc *ServiceMongo) GetTime() (t time.Time) {
	return svc.t
}

func (svc *ServiceMongo) getList(query bson.M, opts *mongo.FindOptions) (results []interface{}, err error) {
	list, err := svc.modelColSvc.GetList(query, opts)
	if err != nil {
		return nil, err
	}
	for _, d := range list.GetModels() {
		r, ok := d.(interfaces.Result)
		if ok {
			results = append(results, r)
		}
	}
	return results, nil
}

func (svc *ServiceMongo) getQuery(query generic.ListQuery) (res bson.M) {
	return utils.GetMongoQuery(query)
}

func (svc *ServiceMongo) getOpts(opts *generic.ListOptions) (res *mongo.FindOptions) {
	return utils.GetMongoOpts(opts)
}

func NewResultServiceMongo(colId primitive.ObjectID, _ primitive.ObjectID) (svc2 interfaces.ResultService, err error) {
	// service
	svc := &ServiceMongo{
		colId: colId,
		t:     time.Now(),
	}

	// dependency injection
	svc.modelSvc, err = service.GetService()
	if err != nil {
		return nil, err
	}

	// data collection
	svc.dc, _ = svc.modelSvc.GetDataCollectionById(colId)
	go func() {
		for {
			time.Sleep(1 * time.Second)
			svc.dc, _ = svc.modelSvc.GetDataCollectionById(colId)
		}
	}()

	// data collection model service
	svc.modelColSvc = service.GetBaseServiceByColName(interfaces.ModelIdResult, svc.dc.Name)

	return svc, nil
}
