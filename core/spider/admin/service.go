package admin

import (
	"context"
	"github.com/apex/log"
	config2 "github.com/crawlab-team/crawlab/core/config"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/container"
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/core/utils"
	"github.com/crawlab-team/crawlab/trace"
	"github.com/crawlab-team/crawlab/vcs"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"
)

type Service struct {
	// dependencies
	nodeCfgSvc   interfaces.NodeConfigService
	modelSvc     service.ModelService
	schedulerSvc interfaces.TaskSchedulerService
	cron         *cron.Cron
	syncLock     bool

	// settings
	cfgPath string
}

func (svc *Service) GetConfigPath() (path string) {
	return svc.cfgPath
}

func (svc *Service) SetConfigPath(path string) {
	svc.cfgPath = path
}

func (svc *Service) Start() (err error) {
	return svc.SyncGit()
}

func (svc *Service) Schedule(id primitive.ObjectID, opts *interfaces.SpiderRunOptions) (taskIds []primitive.ObjectID, err error) {
	// spider
	s, err := svc.modelSvc.GetSpiderById(id)
	if err != nil {
		return nil, err
	}

	// assign tasks
	return svc.scheduleTasks(s, opts)
}

func (svc *Service) Clone(id primitive.ObjectID, opts *interfaces.SpiderCloneOptions) (err error) {
	// TODO: implement
	return nil
}

func (svc *Service) Delete(id primitive.ObjectID) (err error) {
	panic("implement me")
}

func (svc *Service) SyncGit() (err error) {
	if _, err = svc.cron.AddFunc("* * * * *", svc.syncGit); err != nil {
		return trace.TraceError(err)
	}
	svc.cron.Start()
	return nil
}

func (svc *Service) SyncGitOne(g interfaces.Git) (err error) {
	svc.syncGitOne(g)
	return nil
}

func (svc *Service) Export(id primitive.ObjectID) (filePath string, err error) {
	// spider fs
	workspacePath := viper.GetString("workspace")
	spiderFolderPath := filepath.Join(workspacePath, id.Hex())

	// zip files in workspace
	dirPath := spiderFolderPath
	zipFilePath := path.Join(os.TempDir(), uuid.New().String()+".zip")
	if err := utils.ZipDirectory(dirPath, zipFilePath); err != nil {
		return "", trace.TraceError(err)
	}

	return zipFilePath, nil
}

func (svc *Service) scheduleTasks(s *models.Spider, opts *interfaces.SpiderRunOptions) (taskIds []primitive.ObjectID, err error) {
	// main task
	mainTask := &models.Task{
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
			t := &models.Task{
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
			t2, err := svc.schedulerSvc.Enqueue(t)
			if err != nil {
				return nil, err
			}
			taskIds = append(taskIds, t2.GetId())
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
		t2, err := svc.schedulerSvc.Enqueue(mainTask)
		if err != nil {
			return nil, err
		}
		taskIds = append(taskIds, t2.GetId())
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
		nodes, err := svc.modelSvc.GetNodeList(query, nil)
		if err != nil {
			return nil, err
		}
		for _, node := range nodes {
			nodeIds = append(nodeIds, node.GetId())
		}
	} else if opts.Mode == constants.RunTypeSelectedNodes {
		nodeIds = opts.NodeIds
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
		nodes, err := svc.modelSvc.GetNodeList(query, nil)
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

func (svc *Service) syncGit() {
	if svc.syncLock {
		log.Infof("[SpiderAdminService] sync git is locked, skip")
		return
	}
	log.Infof("[SpiderAdminService] start to sync git")

	svc.syncLock = true
	defer func() {
		svc.syncLock = false
	}()

	// spiders
	spiders, err := svc.modelSvc.GetSpiderList(nil, nil)
	if err != nil {
		trace.PrintError(err)
		return
	}

	// spider ids
	var spiderIds []primitive.ObjectID
	for _, s := range spiders {
		spiderIds = append(spiderIds, s.Id)
	}

	if len(spiderIds) > 0 {
		// gits
		gits, err := svc.modelSvc.GetGitList(bson.M{
			"_id": bson.M{
				"$in": spiderIds,
			},
			"auto_pull": true,
		}, nil)
		if err != nil {
			trace.PrintError(err)
			return
		}

		wg := sync.WaitGroup{}
		wg.Add(len(gits))
		for _, g := range gits {
			go func(g models.Git) {
				svc.syncGitOne(&g)
				wg.Done()
			}(g)
		}
		wg.Wait()
	}

	log.Infof("[SpiderAdminService] finished sync git")
}

func (svc *Service) syncGitOne(g interfaces.Git) {
	log.Infof("[SpiderAdminService] sync git %s", g.GetId())

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// git client
	workspacePath := viper.GetString("workspace")
	gitClient, err := vcs.NewGitClient(vcs.WithPath(filepath.Join(workspacePath, g.GetId().Hex())))
	if err != nil {
		return
	}

	// set auth
	utils.InitGitClientAuth(g, gitClient)

	// check if remote has changes
	ok, err := gitClient.IsRemoteChanged()
	if err != nil {
		trace.PrintError(err)
		return
	}
	if !ok {
		// no change
		return
	}

	// pull and sync to workspace
	if err := gitClient.Reset(); err != nil {
		trace.PrintError(err)
		return
	}
	if err := gitClient.Pull(); err != nil {
		trace.PrintError(err)
		return
	}

	// wait for context to end
	<-ctx.Done()
}

func NewSpiderAdminService(opts ...Option) (svc2 interfaces.SpiderAdminService, err error) {
	svc := &Service{
		cfgPath: config2.GetConfigPath(),
	}

	// apply options
	for _, opt := range opts {
		opt(svc)
	}

	// dependency injection
	if err := container.GetContainer().Invoke(func(nodeCfgSvc interfaces.NodeConfigService, modelSvc service.ModelService, schedulerSvc interfaces.TaskSchedulerService) {
		svc.nodeCfgSvc = nodeCfgSvc
		svc.modelSvc = modelSvc
		svc.schedulerSvc = schedulerSvc
	}); err != nil {
		return nil, trace.TraceError(err)
	}

	// cron
	svc.cron = cron.New()

	// validate node type
	if !svc.nodeCfgSvc.IsMaster() {
		return nil, trace.TraceError(errors.ErrorSpiderForbidden)
	}

	return svc, nil
}

var _service interfaces.SpiderAdminService

func GetSpiderAdminService() (svc2 interfaces.SpiderAdminService, err error) {
	if _service != nil {
		return _service, nil
	}

	_service, err = NewSpiderAdminService()
	if err != nil {
		return nil, err
	}

	return _service, nil
}
