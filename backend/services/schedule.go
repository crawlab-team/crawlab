package services

import (
	"crawlab/constants"
	"crawlab/lib/cron"
	"crawlab/model"
	"github.com/apex/log"
	uuid "github.com/satori/go.uuid"
	"runtime/debug"
)

var Sched *Scheduler

type Scheduler struct {
	cron *cron.Cron
}

func AddTask(s model.Schedule) func() {
	return func() {
		nodeId := s.NodeId

		// 生成任务ID
		id := uuid.NewV4()

		// 生成任务模型
		t := model.Task{
			Id:       id.String(),
			SpiderId: s.SpiderId,
			NodeId:   nodeId,
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
	eid, err := s.cron.AddFunc(spec, AddTask(job))
	if err != nil {
		log.Errorf("add func task error: %s", err.Error())
		debug.PrintStack()
		return err
	}

	// 更新EntryID
	job.EntryId = eid
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

		// 添加到定时任务
		if err := s.AddJob(job); err != nil {
			log.Errorf("add job error: %s", err.Error())
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
