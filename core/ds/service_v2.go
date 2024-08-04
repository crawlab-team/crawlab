package ds

import (
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab/core/constants"
	constants2 "github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/models/models/v2"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/core/result"
	"github.com/crawlab-team/crawlab/core/utils"
	utils2 "github.com/crawlab-team/crawlab/core/utils"
	"github.com/crawlab-team/crawlab/trace"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
	"time"
)

type ServiceV2 struct {
	// internals
	timeout         time.Duration
	monitorInterval time.Duration
	stopped         bool
}

func (svc *ServiceV2) Init() {
	// result service registry
	reg := result.GetResultServiceRegistry()

	// register result services
	reg.Register(constants.DataSourceTypeMongo, NewDataSourceMongoService)
	reg.Register(constants.DataSourceTypeMysql, NewDataSourceMysqlService)
	reg.Register(constants.DataSourceTypePostgresql, NewDataSourcePostgresqlService)
	reg.Register(constants.DataSourceTypeMssql, NewDataSourceMssqlService)
	reg.Register(constants.DataSourceTypeSqlite, NewDataSourceSqliteService)
	reg.Register(constants.DataSourceTypeCockroachdb, NewDataSourceCockroachdbService)
	reg.Register(constants.DataSourceTypeElasticSearch, NewDataSourceElasticsearchService)
	reg.Register(constants.DataSourceTypeKafka, NewDataSourceKafkaService)
}

func (svc *ServiceV2) Start() {
	// start monitoring
	go svc.Monitor()
}

func (svc *ServiceV2) Wait() {
	utils.DefaultWait()
}

func (svc *ServiceV2) Stop() {
	svc.stopped = true
}

func (svc *ServiceV2) ChangePassword(id primitive.ObjectID, password string, by primitive.ObjectID) (err error) {
	dataSource, err := service.NewModelServiceV2[models.DatabaseV2]().GetById(id)
	if err != nil {
		return err
	}
	dataSource.Password, err = utils.EncryptAES(password)
	if err != nil {
		return err
	}
	dataSource.SetUpdated(by)
	err = service.NewModelServiceV2[models.DatabaseV2]().ReplaceById(id, *dataSource)
	if err != nil {
		return err
	}
	return nil
}

func (svc *ServiceV2) Monitor() {
	for {
		// return if stopped
		if svc.stopped {
			return
		}

		// monitor
		if err := svc.monitor(); err != nil {
			trace.PrintError(err)
		}

		// wait
		time.Sleep(svc.monitorInterval)
	}
}

func (svc *ServiceV2) CheckStatus(id primitive.ObjectID) (err error) {
	ds, err := service.NewModelServiceV2[models.DatabaseV2]().GetById(id)
	if err != nil {
		return err
	}
	return svc.checkStatus(ds, true)
}

func (svc *ServiceV2) SetTimeout(duration time.Duration) {
	svc.timeout = duration
}

func (svc *ServiceV2) SetMonitorInterval(duration time.Duration) {
	svc.monitorInterval = duration
}

func (svc *ServiceV2) monitor() (err error) {
	// start
	tic := time.Now()
	log.Debugf("[DataSourceService] start monitoring")

	// data source list
	dataSources, err := service.NewModelServiceV2[models.DatabaseV2]().GetMany(nil, nil)
	if err != nil {
		return err
	}

	// waiting group
	wg := sync.WaitGroup{}
	wg.Add(len(dataSources))

	// iterate data source list
	for _, ds := range dataSources {
		// async operation
		go func(ds *models.DatabaseV2) {
			// check status and save
			_ = svc.checkStatus(ds, true)

			// release
			wg.Done()
		}(&ds)
	}

	// wait
	wg.Wait()

	// finish
	toc := time.Now()
	log.Debugf("[DataSourceService] finished monitoring. elapsed: %d ms", (toc.Sub(tic)).Milliseconds())

	return nil
}

func (svc *ServiceV2) checkStatus(ds *models.DatabaseV2, save bool) (err error) {
	// check status
	if err := svc._checkStatus(ds); err != nil {
		ds.Status = constants2.DataSourceStatusOffline
		ds.Error = err.Error()
	} else {
		ds.Status = constants2.DataSourceStatusOnline
		ds.Error = ""
	}

	// save
	if save {
		return svc._save(ds)
	}

	return nil
}

func (svc *ServiceV2) _save(ds *models.DatabaseV2) (err error) {
	log.Debugf("[DataSourceService] saving data source: name=%s, type=%s, status=%s, error=%s", ds.Name, ds.Type, ds.Status, ds.Error)
	return service.NewModelServiceV2[models.DatabaseV2]().ReplaceById(ds.Id, *ds)
}

func (svc *ServiceV2) _checkStatus(ds *models.DatabaseV2) (err error) {
	switch ds.Type {
	case constants.DataSourceTypeMongo:
		_, err := utils2.GetMongoClientWithTimeoutV2(ds, svc.timeout)
		if err != nil {
			return err
		}
	case constants.DataSourceTypeMysql:
		s, err := utils2.GetMysqlSessionWithTimeoutV2(ds, svc.timeout)
		if err != nil {
			return err
		}
		if s != nil {
			err := s.Close()
			if err != nil {
				return err
			}
		}
	case constants.DataSourceTypePostgresql:
		s, err := utils2.GetPostgresqlSessionWithTimeoutV2(ds, svc.timeout)
		if err != nil {
			return err
		}
		if s != nil {
			err := s.Close()
			if err != nil {
				return err
			}
		}
	case constants.DataSourceTypeMssql:
		s, err := utils2.GetMssqlSessionWithTimeoutV2(ds, svc.timeout)
		if err != nil {
			return err
		}
		if s != nil {
			err := s.Close()
			if err != nil {
				return err
			}
		}
	case constants.DataSourceTypeSqlite:
		s, err := utils2.GetSqliteSessionWithTimeoutV2(ds, svc.timeout)
		if err != nil {
			return err
		}
		if s != nil {
			err := s.Close()
			if err != nil {
				return err
			}
		}
	case constants.DataSourceTypeCockroachdb:
		s, err := utils2.GetCockroachdbSessionWithTimeoutV2(ds, svc.timeout)
		if err != nil {
			return err
		}
		if s != nil {
			err := s.Close()
			if err != nil {
				return err
			}
		}
	case constants.DataSourceTypeElasticSearch:
		_, err := utils2.GetElasticsearchClientWithTimeoutV2(ds, svc.timeout)
		if err != nil {
			return err
		}
	case constants.DataSourceTypeKafka:
		c, err := utils2.GetKafkaConnectionWithTimeoutV2(ds, svc.timeout)
		if err != nil {
			return err
		}
		if c != nil {
			err := c.Close()
			if err != nil {
				return err
			}
		}
	default:
		log.Warnf("[DataSourceService] invalid data source type: %s", ds.Type)
	}
	return nil
}

func NewDataSourceServiceV2() *ServiceV2 {
	// service
	svc := &ServiceV2{
		monitorInterval: 15 * time.Second,
		timeout:         10 * time.Second,
	}

	// initialize
	svc.Init()

	// start
	svc.Start()

	return svc
}

var _dsSvcV2 *ServiceV2

func GetDataSourceServiceV2() *ServiceV2 {
	if _dsSvcV2 != nil {
		return _dsSvcV2
	}
	_dsSvcV2 = NewDataSourceServiceV2()
	return _dsSvcV2
}
