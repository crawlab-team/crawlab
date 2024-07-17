package admin

import (
	log2 "github.com/apex/log"
	config2 "github.com/crawlab-team/crawlab/core/config"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	models2 "github.com/crawlab-team/crawlab/core/models/models/v2"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/core/node/config"
	"github.com/crawlab-team/crawlab/core/task/scheduler"
	"github.com/crawlab-team/crawlab/trace"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
	"time"
)

type ServiceV2 struct {
	// dependencies
	nodeCfgSvc   interfaces.NodeConfigService
	schedulerSvc *scheduler.ServiceV2
	cron         *cron.Cron
	syncLock     bool

	// settings
	cfgPath string
}

func (svc *ServiceV2) Schedule(id primitive.ObjectID, opts *interfaces.SpiderRunOptions) (taskIds []primitive.ObjectID, err error) {
	// spider
	s, err := service.NewModelServiceV2[models2.SpiderV2]().GetById(id)
	if err != nil {
		return nil, err
	}

	// assign tasks
	return svc.scheduleTasks(s, opts)
}

func (svc *ServiceV2) scheduleTasks(s *models2.SpiderV2, opts *interfaces.SpiderRunOptions) (taskIds []primitive.ObjectID, err error) {
	// main task
	mainTask := &models2.TaskV2{
		SpiderId:   s.Id,
		Mode:       opts.Mode,
		NodeIds:    opts.NodeIds,
		Cmd:        opts.Cmd,
		Param:      opts.Param,
		ScheduleId: opts.ScheduleId,
		Priority:   opts.Priority,
		UserId:     opts.UserId,
		CreateTs:   time.Now(),
	}
	mainTask.SetId(primitive.NewObjectID())

	// normalize
	if mainTask.Mode == "" {
		mainTask.Mode = s.Mode
	}
	if mainTask.NodeIds == nil {
		mainTask.NodeIds = s.NodeIds
	}
	if mainTask.Cmd == "" {
		mainTask.Cmd = s.Cmd
	}
	if mainTask.Param == "" {
		mainTask.Param = s.Param
	}
	if mainTask.Priority == 0 {
		mainTask.Priority = s.Priority
	}

	if svc.isMultiTask(opts) {
		// multi tasks
		nodeIds, err := svc.getNodeIds(opts)
		if err != nil {
			return nil, err
		}
		for _, nodeId := range nodeIds {
			t := &models2.TaskV2{
				SpiderId:   s.Id,
				Mode:       opts.Mode,
				Cmd:        opts.Cmd,
				Param:      opts.Param,
				NodeId:     nodeId,
				ScheduleId: opts.ScheduleId,
				Priority:   opts.Priority,
				UserId:     opts.UserId,
				CreateTs:   time.Now(),
			}
			t.SetId(primitive.NewObjectID())
			t2, err := svc.schedulerSvc.Enqueue(t, opts.UserId)
			if err != nil {
				return nil, err
			}
			taskIds = append(taskIds, t2.Id)
		}
	} else {
		// single task
		nodeIds, err := svc.getNodeIds(opts)
		if err != nil {
			return nil, err
		}
		if len(nodeIds) > 0 {
			mainTask.NodeId = nodeIds[0]
		}
		t2, err := svc.schedulerSvc.Enqueue(mainTask, opts.UserId)
		if err != nil {
			return nil, err
		}
		taskIds = append(taskIds, t2.Id)
	}

	return taskIds, nil
}

func (svc *ServiceV2) getNodeIds(opts *interfaces.SpiderRunOptions) (nodeIds []primitive.ObjectID, err error) {
	if opts.Mode == constants.RunTypeAllNodes {
		query := bson.M{
			"active":  true,
			"enabled": true,
			"status":  constants.NodeStatusOnline,
		}
		nodes, err := service.NewModelServiceV2[models2.NodeV2]().GetMany(query, nil)
		if err != nil {
			return nil, err
		}
		for _, node := range nodes {
			nodeIds = append(nodeIds, node.Id)
		}
	} else if opts.Mode == constants.RunTypeSelectedNodes {
		nodeIds = opts.NodeIds
	}
	return nodeIds, nil
}

func (svc *ServiceV2) isMultiTask(opts *interfaces.SpiderRunOptions) (res bool) {
	if opts.Mode == constants.RunTypeAllNodes {
		query := bson.M{
			"active":  true,
			"enabled": true,
			"status":  constants.NodeStatusOnline,
		}
		nodes, err := service.NewModelServiceV2[models2.NodeV2]().GetMany(query, nil)
		if err != nil {
			trace.PrintError(err)
			return false
		}
		return len(nodes) > 1
	} else if opts.Mode == constants.RunTypeRandom {
		return false
	} else if opts.Mode == constants.RunTypeSelectedNodes {
		return len(opts.NodeIds) > 1
	} else {
		return false
	}
}

func newSpiderAdminServiceV2() (svc2 *ServiceV2, err error) {
	svc := &ServiceV2{
		nodeCfgSvc: config.GetNodeConfigService(),
		cfgPath:    config2.GetConfigPath(),
	}
	svc.schedulerSvc, err = scheduler.GetTaskSchedulerServiceV2()
	if err != nil {
		return nil, err
	}

	// cron
	svc.cron = cron.New()

	// validate node type
	if !svc.nodeCfgSvc.IsMaster() {
		return nil, trace.TraceError(errors.ErrorSpiderForbidden)
	}

	return svc, nil
}

var svcV2 *ServiceV2
var svcV2Once = new(sync.Once)

func GetSpiderAdminServiceV2() (svc2 *ServiceV2, err error) {
	if svcV2 != nil {
		return svcV2, nil
	}
	svcV2Once.Do(func() {
		svcV2, err = newSpiderAdminServiceV2()
		if err != nil {
			log2.Errorf("[GetSpiderAdminServiceV2] error: %v", err)
		}
	})
	if err != nil {
		return nil, err
	}
	return svcV2, nil
}
