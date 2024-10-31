package admin

import (
	"errors"
	log2 "github.com/apex/log"
	config2 "github.com/crawlab-team/crawlab/core/config"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/interfaces"
	models2 "github.com/crawlab-team/crawlab/core/models/models/v2"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/core/node/config"
	"github.com/crawlab-team/crawlab/core/task/scheduler"
	"github.com/crawlab-team/crawlab/trace"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
)

type Service struct {
	// dependencies
	nodeCfgSvc   interfaces.NodeConfigService
	schedulerSvc *scheduler.Service
	syncLock     bool

	// settings
	cfgPath string
}

func (svc *Service) Schedule(id primitive.ObjectID, opts *interfaces.SpiderRunOptions) (taskIds []primitive.ObjectID, err error) {
	// spider
	s, err := service.NewModelServiceV2[models2.SpiderV2]().GetById(id)
	if err != nil {
		return nil, err
	}

	// assign tasks
	return svc.scheduleTasks(s, opts)
}

func (svc *Service) scheduleTasks(s *models2.SpiderV2, opts *interfaces.SpiderRunOptions) (taskIds []primitive.ObjectID, err error) {
	// get node ids
	nodeIds, err := svc.getNodeIds(opts)
	if err != nil {
		return nil, err
	}

	// iterate node ids
	for _, nodeId := range nodeIds {
		// task
		t := &models2.TaskV2{
			SpiderId:   s.Id,
			NodeId:     nodeId,
			NodeIds:    opts.NodeIds,
			Mode:       opts.Mode,
			Cmd:        opts.Cmd,
			Param:      opts.Param,
			ScheduleId: opts.ScheduleId,
			Priority:   opts.Priority,
		}

		// normalize
		if t.Mode == "" {
			t.Mode = s.Mode
		}
		if t.NodeIds == nil {
			t.NodeIds = s.NodeIds
		}
		if t.Cmd == "" {
			t.Cmd = s.Cmd
		}
		if t.Param == "" {
			t.Param = s.Param
		}
		if t.Priority == 0 {
			t.Priority = s.Priority
		}

		// enqueue task
		t, err = svc.schedulerSvc.Enqueue(t, opts.UserId)
		if err != nil {
			return nil, err
		}

		// append task id
		taskIds = append(taskIds, t.Id)
	}

	return taskIds, nil
}

func (svc *Service) getNodeIds(opts *interfaces.SpiderRunOptions) (nodeIds []primitive.ObjectID, err error) {
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
	} else if opts.Mode == constants.RunTypeRandom {
		nodeIds = []primitive.ObjectID{primitive.NilObjectID}
	} else {
		return nil, errors.New("invalid run mode")
	}
	return nodeIds, nil
}

func (svc *Service) isMultiTask(opts *interfaces.SpiderRunOptions) (res bool) {
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

func newSpiderAdminService() (svc2 *Service, err error) {
	svc := &Service{
		nodeCfgSvc: config.GetNodeConfigService(),
		cfgPath:    config2.GetConfigPath(),
	}
	svc.schedulerSvc, err = scheduler.GetTaskSchedulerService()
	if err != nil {
		return nil, err
	}

	// validate node type
	if !svc.nodeCfgSvc.IsMaster() {
		return nil, errors.New("only master node can run spider admin service")
	}

	return svc, nil
}

var svc *Service
var svcOnce = new(sync.Once)

func GetSpiderAdminService() (svc2 *Service, err error) {
	if svc != nil {
		return svc, nil
	}
	svcOnce.Do(func() {
		svc, err = newSpiderAdminService()
		if err != nil {
			log2.Errorf("[GetSpiderAdminService] error: %v", err)
		}
	})
	if err != nil {
		return nil, err
	}
	return svc, nil
}
