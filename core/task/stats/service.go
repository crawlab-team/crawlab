package stats

import (
	"github.com/crawlab-team/crawlab/core/container"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/core/result"
	"github.com/crawlab-team/crawlab/core/task"
	"github.com/crawlab-team/crawlab/core/task/log"
	"github.com/crawlab-team/crawlab/db/mongo"
	"github.com/crawlab-team/crawlab/trace"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
	"time"
)

type Service struct {
	// dependencies
	interfaces.TaskBaseService
	nodeCfgSvc interfaces.NodeConfigService
	modelSvc   service.ModelService

	// internals
	mu             sync.Mutex
	resultServices sync.Map
	rsTtl          time.Duration
	logDriver      log.Driver
}

func (svc *Service) Init() (err error) {
	go svc.cleanup()
	return nil
}

func (svc *Service) InsertData(id primitive.ObjectID, records ...interface{}) (err error) {
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

func (svc *Service) InsertLogs(id primitive.ObjectID, logs ...string) (err error) {
	return svc.logDriver.WriteLines(id.Hex(), logs)
}

func (svc *Service) getResultService(id primitive.ObjectID) (resultSvc interfaces.ResultService, err error) {
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

func (svc *Service) updateTaskStats(id primitive.ObjectID, resultCount int) {
	_ = mongo.GetMongoCol(interfaces.ModelColNameTaskStat).UpdateId(id, bson.M{
		"$inc": bson.M{
			"result_count": resultCount,
		},
	})
}

func (svc *Service) cleanup() {
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

func NewTaskStatsService() (svc2 interfaces.TaskStatsService, err error) {
	// base service
	baseSvc, err := task.NewBaseService()
	if err != nil {
		return nil, trace.TraceError(err)
	}

	// service
	svc := &Service{
		mu:              sync.Mutex{},
		TaskBaseService: baseSvc,
		resultServices:  sync.Map{},
	}

	// dependency injection
	if err := container.GetContainer().Invoke(func(nodeCfgSvc interfaces.NodeConfigService, modelSvc service.ModelService) {
		svc.nodeCfgSvc = nodeCfgSvc
		svc.modelSvc = modelSvc
	}); err != nil {
		return nil, trace.TraceError(err)
	}

	// log driver
	svc.logDriver, err = log.GetLogDriver(log.DriverTypeFile)
	if err != nil {
		return nil, err
	}

	return svc, nil
}

var _service interfaces.TaskStatsService

func GetTaskStatsService() (svr interfaces.TaskStatsService, err error) {
	if _service != nil {
		return _service, nil
	}
	_service, err = NewTaskStatsService()
	if err != nil {
		return nil, err
	}
	return _service, nil
}
