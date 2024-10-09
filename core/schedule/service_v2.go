package schedule

import (
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab/core/config"
	"github.com/crawlab-team/crawlab/core/interfaces"
	models2 "github.com/crawlab-team/crawlab/core/models/models/v2"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/core/spider/admin"
	"github.com/crawlab-team/crawlab/core/utils"
	"github.com/crawlab-team/crawlab/trace"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
	"time"
)

type ServiceV2 struct {
	// dependencies
	interfaces.WithConfigPath
	modelSvc *service.ModelServiceV2[models2.ScheduleV2]
	adminSvc *admin.ServiceV2

	// settings variables
	loc            *time.Location
	delay          bool
	skip           bool
	updateInterval time.Duration

	// internals
	cron      *cron.Cron
	logger    cron.Logger
	schedules []models2.ScheduleV2
	stopped   bool
	mu        sync.Mutex
}

func (svc *ServiceV2) GetLocation() (loc *time.Location) {
	return svc.loc
}

func (svc *ServiceV2) SetLocation(loc *time.Location) {
	svc.loc = loc
}

func (svc *ServiceV2) GetDelay() (delay bool) {
	return svc.delay
}

func (svc *ServiceV2) SetDelay(delay bool) {
	svc.delay = delay
}

func (svc *ServiceV2) GetSkip() (skip bool) {
	return svc.skip
}

func (svc *ServiceV2) SetSkip(skip bool) {
	svc.skip = skip
}

func (svc *ServiceV2) GetUpdateInterval() (interval time.Duration) {
	return svc.updateInterval
}

func (svc *ServiceV2) SetUpdateInterval(interval time.Duration) {
	svc.updateInterval = interval
}

func (svc *ServiceV2) Init() (err error) {
	return svc.fetch()
}

func (svc *ServiceV2) Start() {
	svc.cron.Start()
	go svc.Update()
}

func (svc *ServiceV2) Wait() {
	utils.DefaultWait()
	svc.Stop()
}

func (svc *ServiceV2) Stop() {
	svc.stopped = true
	svc.cron.Stop()
}

func (svc *ServiceV2) Enable(s models2.ScheduleV2, by primitive.ObjectID) (err error) {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	id, err := svc.cron.AddFunc(s.Cron, svc.schedule(s.Id))
	if err != nil {
		return trace.TraceError(err)
	}
	s.Enabled = true
	s.EntryId = id
	s.SetUpdated(by)
	return svc.modelSvc.ReplaceById(s.Id, s)
}

func (svc *ServiceV2) Disable(s models2.ScheduleV2, by primitive.ObjectID) (err error) {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	svc.cron.Remove(s.EntryId)
	s.Enabled = false
	s.EntryId = -1
	s.SetUpdated(by)
	return svc.modelSvc.ReplaceById(s.Id, s)
}

func (svc *ServiceV2) Update() {
	for {
		if svc.stopped {
			return
		}

		svc.update()

		time.Sleep(svc.updateInterval)
	}
}

func (svc *ServiceV2) GetCron() (c *cron.Cron) {
	return svc.cron
}

func (svc *ServiceV2) update() {
	// fetch enabled schedules
	if err := svc.fetch(); err != nil {
		trace.PrintError(err)
		return
	}

	// entry id map
	entryIdsMap := svc.getEntryIdsMap()

	// iterate enabled schedules
	for _, s := range svc.schedules {
		_, ok := entryIdsMap[s.EntryId]
		if ok {
			entryIdsMap[s.EntryId] = true
		} else {
			if !s.Enabled {
				err := svc.Enable(s, s.GetCreatedBy())
				if err != nil {
					trace.PrintError(err)
					continue
				}
			}
		}
	}

	// remove non-existent entries
	for id, ok := range entryIdsMap {
		if !ok {
			svc.cron.Remove(id)
		}
	}
}

func (svc *ServiceV2) getEntryIdsMap() (res map[cron.EntryID]bool) {
	res = map[cron.EntryID]bool{}
	for _, e := range svc.cron.Entries() {
		res[e.ID] = false
	}
	return res
}

func (svc *ServiceV2) fetch() (err error) {
	query := bson.M{
		"enabled": true,
	}
	svc.schedules, err = svc.modelSvc.GetMany(query, nil)
	if err != nil {
		return err
	}
	return nil
}

func (svc *ServiceV2) schedule(id primitive.ObjectID) (fn func()) {
	return func() {
		// schedule
		s, err := svc.modelSvc.GetById(id)
		if err != nil {
			trace.PrintError(err)
			return
		}

		// spider
		spider, err := service.NewModelServiceV2[models2.SpiderV2]().GetById(s.SpiderId)
		if err != nil {
			trace.PrintError(err)
			return
		}

		// options
		opts := &interfaces.SpiderRunOptions{
			Mode:       s.Mode,
			NodeIds:    s.NodeIds,
			Cmd:        s.Cmd,
			Param:      s.Param,
			Priority:   s.Priority,
			ScheduleId: s.Id,
			UserId:     s.GetCreatedBy(),
		}

		// normalize options
		if opts.Mode == "" {
			opts.Mode = spider.Mode
		}
		if len(opts.NodeIds) == 0 {
			opts.NodeIds = spider.NodeIds
		}
		if opts.Cmd == "" {
			opts.Cmd = spider.Cmd
		}
		if opts.Param == "" {
			opts.Param = spider.Param
		}
		if opts.Priority == 0 {
			if spider.Priority > 0 {
				opts.Priority = spider.Priority
			} else {
				opts.Priority = 5
			}
		}

		// schedule or assign a task in the task queue
		if _, err := svc.adminSvc.Schedule(s.SpiderId, opts); err != nil {
			trace.PrintError(err)
		}
	}
}

func NewScheduleServiceV2() (svc2 *ServiceV2, err error) {
	// service
	svc := &ServiceV2{
		WithConfigPath: config.NewConfigPathService(),
		loc:            time.Local,
		// TODO: implement delay and skip
		delay:          false,
		skip:           false,
		updateInterval: 1 * time.Minute,
	}
	svc.adminSvc, err = admin.GetSpiderAdminServiceV2()
	if err != nil {
		return nil, err
	}
	svc.modelSvc = service.NewModelServiceV2[models2.ScheduleV2]()

	// logger
	svc.logger = NewLogger()

	// cron
	svc.cron = cron.New(
		cron.WithLogger(svc.logger),
		cron.WithLocation(svc.loc),
		cron.WithChain(cron.Recover(svc.logger)),
	)

	// initialize
	if err := svc.Init(); err != nil {
		return nil, err
	}

	return svc, nil
}

var svcV2 *ServiceV2
var svcV2Once = new(sync.Once)

func GetScheduleServiceV2() (res *ServiceV2, err error) {
	if svcV2 != nil {
		return svcV2, nil
	}
	svcV2Once.Do(func() {
		svcV2, err = NewScheduleServiceV2()
		if err != nil {
			log.Errorf("failed to get schedule service: %v", err)
		}
	})
	if err != nil {
		return nil, err
	}
	return svcV2, nil
}
