package ds

import (
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/core/utils"
	utils2 "github.com/crawlab-team/crawlab/core/utils"
	"github.com/crawlab-team/crawlab/db/generic"
	"github.com/crawlab-team/crawlab/db/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"time"
)

type MongoService struct {
	// dependencies
	modelSvc service.ModelService

	// internals
	dc  *models.DataCollection // models.DataCollection
	ds  *models.DataSource     // models.DataSource
	c   *mongo2.Client
	db  *mongo2.Database
	col *mongo.Col
	t   time.Time
}

func (svc *MongoService) Insert(records ...interface{}) (err error) {
	_, err = svc.col.InsertMany(records)
	return err
}

func (svc *MongoService) List(query generic.ListQuery, opts *generic.ListOptions) (results []interface{}, err error) {
	var docs []models.Result
	if err := svc.col.Find(utils.GetMongoQuery(query), utils.GetMongoOpts(opts)).All(&docs); err != nil {
		return nil, err
	}
	for i := range docs {
		results = append(results, &docs[i])
	}
	return results, nil
}

func (svc *MongoService) Count(query generic.ListQuery) (n int, err error) {
	return svc.col.Count(utils.GetMongoQuery(query))
}

func NewDataSourceMongoService(colId primitive.ObjectID, dsId primitive.ObjectID) (svc2 interfaces.ResultService, err error) {
	// service
	svc := &MongoService{}

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
		svc.ds.Port = constants.DefaultMongoPort
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

	// mongo client
	svc.c, err = utils2.GetMongoClient(svc.ds)
	if err != nil {
		return nil, err
	}

	// mongo database
	svc.db = mongo.GetMongoDb(svc.ds.Database, mongo.WithDbClient(svc.c))

	// mongo col
	svc.col = mongo.GetMongoColWithDb(svc.dc.Name, svc.db)

	return svc, nil
}

func (svc *MongoService) Index(fields []string) {
	// TODO: implement me
}

func (svc *MongoService) SetTime(t time.Time) {
	svc.t = t
}

func (svc *MongoService) GetTime() (t time.Time) {
	return svc.t
}
