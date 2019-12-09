package services

import (
	"crawlab/constants"
	"crawlab/lib/cron"
	"crawlab/model"
	"errors"
	"github.com/apex/log"
	"github.com/globalsign/mgo/bson"
	"github.com/satori/go.uuid"
	"runtime/debug"
)

var Sched *Scheduler

type Scheduler struct {
	cron *cron.Cron
}

func AddScheduleTask(s model.Schedule) func() {
	return func() {
		node, err := model.GetNodeByKey(s.NodeKey)
		if err != nil || node.Id.Hex() == "" {
			log.Errorf("get node by key error: %s", err.Error())
			debug.PrintStack()
			return
		}

		spider := model.GetSpiderByName(s.SpiderName)
		if spider == nil || spider.Id.Hex() == "" {
			log.Errorf("get spider by name error: %s", err.Error())
			debug.PrintStack()
			return
		}

		// 同步ID到定时任务
		s.SyncNodeIdAndSpiderId(node, *spider)

		// 生成任务ID
		id := uuid.NewV4()

		// 生成任务模型
		t := model.Task{
			Id:       id.String(),
			SpiderId: spider.Id,
			NodeId:   node.Id,
			Status:   constants.StatusPending,
			Param:    s.Param,
		}

		// 将任务存入数据库
		if err := model.AddTask(t); err != nil {
			log.Errorf(err.Error())
			debug.PrintStack()
			return
		}

		// 加入任务队列
		if err := AssignTask(t); err != nil {
			log.Errorf(err.Error())
			debug.PrintStack()
			return
		}
	}
}

func UpdateSchedules() {
	if err := Sched.Update(); err != nil {
		log.Errorf(err.Error())
		return
	}
}

func (s *Scheduler) Start() error {
	exec := cron.New(cron.WithSeconds())

	// 启动cron服务
	s.cron.Start()

	// 更新任务列表
	if err := s.Update(); err != nil {
		log.Errorf("update scheduler error: %s", err.Error())
		debug.PrintStack()
		return err
	}

	// 每30秒更新一次任务列表
	spec := "*/30 * * * * *"
	if _, err := exec.AddFunc(spec, UpdateSchedules); err != nil {
		log.Errorf("add func update schedulers error: %s", err.Error())
		debug.PrintStack()
		return err
	}

	return nil
}

func (s *Scheduler) AddJob(job model.Schedule) error {
	spec := job.Cron

	// 添加任务
	eid, err := s.cron.AddFunc(spec, AddScheduleTask(job))
	if err != nil {
		log.Errorf("add func task error: %s", err.Error())
		debug.PrintStack()
		return err
	}

	// 更新EntryID
	job.EntryId = eid
	job.Status = constants.ScheduleStatusRunning
	if err := job.Save(); err != nil {
		log.Errorf("job save error: %s", err.Error())
		debug.PrintStack()
		return err
	}

	return nil
}

func (s *Scheduler) RemoveAll() {
	entries := s.cron.Entries()
	for i := 0; i < len(entries); i++ {
		s.cron.Remove(entries[i].ID)
	}
}

// 验证cron表达式是否正确
func ParserCron(spec string) error {
	parser := cron.NewParser(
		cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)

	if _, err := parser.Parse(spec); err != nil {
		return err
	}
	return nil
}

// 停止定时任务
func (s *Scheduler) Stop(id bson.ObjectId) error {
	schedule, err := model.GetSchedule(id)
	if err != nil {
		return err
	}
	if schedule.EntryId == 0 {
		return errors.New("entry id not found")
	}
	s.cron.Remove(schedule.EntryId)
	// 更新状态
	schedule.Status = constants.ScheduleStatusStop
	if err = schedule.Save(); err != nil {
		return err
	}
	return nil
}

// 运行任务
func (s *Scheduler) Run(id bson.ObjectId) error {
	schedule, err := model.GetSchedule(id)
	if err != nil {
		return err
	}
	if err := s.AddJob(schedule); err != nil {
		return err
	}
	return nil
}

func (s *Scheduler) Update() error {
	// 删除所有定时任务
	s.RemoveAll()

	// 获取所有定时任务
	sList, err := model.GetScheduleList(nil)
	if err != nil {
		log.Errorf("get scheduler list error: %s", err.Error())
		debug.PrintStack()
		return err
	}

	// 遍历任务列表
	for i := 0; i < len(sList); i++ {
		// 单个任务
		job := sList[i]

		if job.Status == constants.ScheduleStatusStop {
			continue
		}

		// 添加到定时任务
		if err := s.AddJob(job); err != nil {
			log.Errorf("add job error: %s, job: %s, cron: %s", err.Error(), job.Name, job.Cron)
			debug.PrintStack()
			return err
		}
	}

	return nil
}

func InitScheduler() error {
	Sched = &Scheduler{
		cron: cron.New(cron.WithSeconds()),
	}
	if err := Sched.Start(); err != nil {
		log.Errorf("start scheduler error: %s", err.Error())
		debug.PrintStack()
		return err
	}
	return nil
}
