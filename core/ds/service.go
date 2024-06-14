package ds

import (
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab/core/constants"
	constants2 "github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/container"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/delegate"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/core/result"
	"github.com/crawlab-team/crawlab/core/utils"
	utils2 "github.com/crawlab-team/crawlab/core/utils"
	"github.com/crawlab-team/crawlab/trace"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
	"time"
)

type Service struct {
	// dependencies
	modelSvc service.ModelService

	// internals
	timeout         time.Duration
	monitorInterval time.Duration
	stopped         bool
}

func (svc *Service) Init() {
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

func (svc *Service) Start() {
	// start monitoring
	go svc.Monitor()
}

func (svc *Service) Wait() {
	utils.DefaultWait()
}

func (svc *Service) Stop() {
	svc.stopped = true
}

func (svc *Service) ChangePassword(id primitive.ObjectID, password string) (err error) {
	p, err := svc.modelSvc.GetPasswordById(id)
	if err == nil {
		// exists, save
		encryptedPassword, err := utils.EncryptAES(password)
		if err != nil {
			return err
		}
		p.Password = encryptedPassword
		if err := delegate.NewModelDelegate(p).Save(); err != nil {
			return err
		}
		return nil
	} else if err.Error() == mongo.ErrNoDocuments.Error() {
		// not exists, add
		encryptedPassword, err := utils.EncryptAES(password)
		if err != nil {
			return err
		}
		p = &models.Password{
			Id:       id,
			Password: encryptedPassword,
		}
		if err := delegate.NewModelDelegate(p).Add(); err != nil {
			return err
		}
		return nil
	} else {
		// error
		return err
	}
}

func (svc *Service) Monitor() {
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

func (svc *Service) CheckStatus(id primitive.ObjectID) (err error) {
	ds, err := svc.modelSvc.GetDataSourceById(id)
	if err != nil {
		return err
	}
	return svc.checkStatus(ds, true)
}

func (svc *Service) SetTimeout(duration time.Duration) {
	svc.timeout = duration
}

func (svc *Service) SetMonitorInterval(duration time.Duration) {
	svc.monitorInterval = duration
}

func (svc *Service) monitor() (err error) {
	// start
	tic := time.Now()
	log.Debugf("[DataSourceService] start monitoring")

	// data source list
	dsList, err := svc.modelSvc.GetDataSourceList(nil, nil)
	if err != nil {
		return err
	}

	// waiting group
	wg := sync.WaitGroup{}
	wg.Add(len(dsList))

	// iterate data source list
	for _, ds := range dsList {
		// async operation
		go func(ds models.DataSource) {
			// check status and save
			_ = svc.checkStatus(&ds, true)

			// release
			wg.Done()
		}(ds)
	}

	// wait
	wg.Wait()

	// finish
	toc := time.Now()
	log.Debugf("[DataSourceService] finished monitoring. elapsed: %d ms", (toc.Sub(tic)).Milliseconds())

	return nil
}

func (svc *Service) checkStatus(ds *models.DataSource, save bool) (err error) {
	// password
	if ds.Password == "" {
		pwd, err := svc.modelSvc.GetPasswordById(ds.Id)
		if err == nil {
			ds.Password, err = utils.DecryptAES(pwd.Password)
			if err != nil {
				return err
			}
		} else if err.Error() != mongo.ErrNoDocuments.Error() {
			return trace.TraceError(err)
		}
	}

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

func (svc *Service) _save(ds *models.DataSource) (err error) {
	log.Debugf("[DataSourceService] saving data source: name=%s, type=%s, status=%s, error=%s", ds.Name, ds.Type, ds.Status, ds.Error)
	return delegate.NewModelDelegate(ds).Save()
}

func (svc *Service) _checkStatus(ds *models.DataSource) (err error) {
	switch ds.Type {
	case constants.DataSourceTypeMongo:
		_, err := utils2.GetMongoClientWithTimeout(ds, svc.timeout)
		if err != nil {
			return err
		}
	case constants.DataSourceTypeMysql:
		s, err := utils2.GetMysqlSessionWithTimeout(ds, svc.timeout)
		if err != nil {
			return err
		}
		if s != nil {
			s.Close()
		}
	case constants.DataSourceTypePostgresql:
		s, err := utils2.GetPostgresqlSessionWithTimeout(ds, svc.timeout)
		if err != nil {
			return err
		}
		if s != nil {
			s.Close()
		}
	case constants.DataSourceTypeMssql:
		s, err := utils2.GetMssqlSessionWithTimeout(ds, svc.timeout)
		if err != nil {
			return err
		}
		if s != nil {
			s.Close()
		}
	case constants.DataSourceTypeSqlite:
		s, err := utils2.GetSqliteSessionWithTimeout(ds, svc.timeout)
		if err != nil {
			return err
		}
		if s != nil {
			s.Close()
		}
	case constants.DataSourceTypeCockroachdb:
		s, err := utils2.GetCockroachdbSessionWithTimeout(ds, svc.timeout)
		if err != nil {
			return err
		}
		if s != nil {
			s.Close()
		}
	case constants.DataSourceTypeElasticSearch:
		_, err := utils2.GetElasticsearchClientWithTimeout(ds, svc.timeout)
		if err != nil {
			return err
		}
	case constants.DataSourceTypeKafka:
		c, err := utils2.GetKafkaConnectionWithTimeout(ds, svc.timeout)
		if err != nil {
			return err
		}
		if c != nil {
			c.Close()
		}
	default:
		log.Warnf("[DataSourceService] invalid data source type: %s", ds.Type)
	}
	return nil
}

func NewDataSourceService(opts ...DataSourceServiceOption) (svc2 interfaces.DataSourceService, err error) {
	// service
	svc := &Service{
		monitorInterval: 15 * time.Second,
		timeout:         10 * time.Second,
	}

	// apply options
	for _, opt := range opts {
		opt(svc)
	}

	// dependency injection
	if err := container.GetContainer().Invoke(func(modelSvc service.ModelService) {
		svc.modelSvc = modelSvc
	}); err != nil {
		return nil, trace.TraceError(err)
	}

	// initialize
	svc.Init()

	// start
	svc.Start()

	return svc, nil
}

var _dsSvc interfaces.DataSourceService

func GetDataSourceService() (svc interfaces.DataSourceService, err error) {
	if _dsSvc != nil {
		return _dsSvc, nil
	}
	svc, err = NewDataSourceService()
	if err != nil {
		return nil, err
	}
	_dsSvc = svc
	return svc, nil
}
