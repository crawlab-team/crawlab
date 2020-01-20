package services

import (
	"crawlab/constants"
	"crawlab/lib/cron"
	"crawlab/model"
	"errors"
	"github.com/apex/log"
	"github.com/globalsign/mgo/bson"
	uuid "github.com/satori/go.uuid"
	"runtime/debug"
)

var Sched *Scheduler

type Scheduler struct {
	cron *cron.Cron
}

func AddScheduleTask(s model.Schedule) func() {
	return func() {
		// 生成任务ID
		id := uuid.NewV4()

		if s.RunType == constants.RunTypeAllNodes {
			// 所有节点
			nodes, err := model.GetNodeList(nil)
			if err != nil {
				return
			}
			for _, node := range nodes {
				t := model.Task{
					Id:       id.String(),
					SpiderId: s.SpiderId,
					NodeId:   node.Id,
					Param:    s.Param,
					UserId:   s.UserId,
				}

				if err := AddTask(t); err != nil {
					return
				}
			}
		} else if s.RunType == constants.RunTypeRandom {
			// 随机
			t := model.Task{
				Id:       id.String(),
				SpiderId: s.SpiderId,
				Param:    s.Param,
				UserId:   s.UserId,
			}
			if err := AddTask(t); err != nil {
				log.Errorf(err.Error())
				debug.PrintStack()
				return
			}
		} else if s.RunType == constants.RunTypeSelectedNodes {
			// 指定节点
			for _, nodeId := range s.NodeIds {
				t := model.Task{
					Id:       id.String(),
					SpiderId: s.SpiderId,
					NodeId:   nodeId,
					Param:    s.Param,
					UserId:   s.UserId,
				}

				if err := AddTask(t); err != nil {
					return
				}
			}
		} else {
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

	// 添加定时任务
	eid, err := s.cron.AddFunc(spec, AddScheduleTask(job))
	if err != nil {
		log.Errorf("add func task error: %s", err.Error())
		debug.PrintStack()
		return err
	}

	// 更新EntryID
	job.EntryId = eid

	// 更新状态
	job.Status = constants.ScheduleStatusRunning
	job.Enabled = true

	// 保存定时任务
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

// 禁用定时任务
func (s *Scheduler) Disable(id bson.ObjectId) error {
	schedule, err := model.GetSchedule(id)
	if err != nil {
		return err
	}
	if schedule.EntryId == 0 {
		return errors.New("entry id not found")
	}

	// 从cron服务中删除该任务
	s.cron.Remove(schedule.EntryId)

	// 更新状态
	schedule.Status = constants.ScheduleStatusStop
	schedule.Enabled = false

	if err = schedule.Save(); err != nil {
		return err
	}
	return nil
}

// 启用定时任务
func (s *Scheduler) Enable(id bson.ObjectId) error {
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
