package stats

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/models/service"
	nodeconfig "github.com/crawlab-team/crawlab/core/node/config"
	"github.com/crawlab-team/crawlab/core/result"
	"github.com/crawlab-team/crawlab/core/task/log"
	"github.com/crawlab-team/crawlab/trace"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
	"time"
)

type ServiceV2 struct {
	// dependencies
	nodeCfgSvc interfaces.NodeConfigService
	modelSvc   service.ModelService

	// internals
	mu             sync.Mutex
	resultServices sync.Map
	rsTtl          time.Duration
	logDriver      log.Driver
}

func (svc *ServiceV2) Init() (err error) {
	go svc.cleanup()
	return nil
}

func (svc *ServiceV2) InsertData(id primitive.ObjectID, records ...interface{}) (err error) {
	resultSvc, err := svc.getResultService(id)
	if err != nil {
		return err
	}
	if err := resultSvc.Insert(records...); err != nil {
		return err
	}
	go svc.updateTaskStats(id, len(records))
	return nil
}

func (svc *ServiceV2) InsertLogs(id primitive.ObjectID, logs ...string) (err error) {
	return svc.logDriver.WriteLines(id.Hex(), logs)
}

func (svc *ServiceV2) getResultService(id primitive.ObjectID) (resultSvc interfaces.ResultService, err error) {
	// atomic operation
	svc.mu.Lock()
	defer svc.mu.Unlock()

	// attempt to get from cache
	res, _ := svc.resultServices.Load(id.Hex())
	if res != nil {
		// hit in cache
		resultSvc, ok := res.(interfaces.ResultService)
		resultSvc.SetTime(time.Now())
		if ok {
			return resultSvc, nil
		}
	}

	// task
	t, err := svc.modelSvc.GetTaskById(id)
	if err != nil {
		return nil, err
	}

	// result service
	resultSvc, err = result.GetResultService(t.SpiderId)
	if err != nil {
		return nil, err
	}

	// store in cache
	svc.resultServices.Store(id.Hex(), resultSvc)

	return resultSvc, nil
}

func (svc *ServiceV2) updateTaskStats(id primitive.ObjectID, resultCount int) {
	err := service.NewModelServiceV2[models.TaskStatV2]().UpdateById(id, bson.M{
		"$inc": bson.M{
			"result_count": resultCount,
		},
	})
	if err != nil {
		trace.PrintError(err)
	}
}

func (svc *ServiceV2) cleanup() {
	for {
		// atomic operation
		svc.mu.Lock()

		svc.resultServices.Range(func(key, value interface{}) bool {
			rs := value.(interfaces.ResultService)
			if time.Now().After(rs.GetTime().Add(svc.rsTtl)) {
				svc.resultServices.Delete(key)
			}
			return true
		})

		svc.mu.Unlock()

		time.Sleep(10 * time.Minute)
	}
}

func NewTaskStatsServiceV2() (svc2 *ServiceV2, err error) {
	// service
	svc := &ServiceV2{
		mu:             sync.Mutex{},
		resultServices: sync.Map{},
	}

	svc.nodeCfgSvc = nodeconfig.GetNodeConfigService()

	// log driver
	svc.logDriver, err = log.GetLogDriver(log.DriverTypeFile)
	if err != nil {
		return nil, err
	}

	return svc, nil
}

var _serviceV2 *ServiceV2

func GetTaskStatsServiceV2() (svr *ServiceV2, err error) {
	if _serviceV2 != nil {
		return _serviceV2, nil
	}
	_serviceV2, err = NewTaskStatsServiceV2()
	if err != nil {
		return nil, err
	}
	return _serviceV2, nil
}
