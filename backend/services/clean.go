package services

import (
	"crawlab/constants"
	"crawlab/model"
	"github.com/apex/log"
	"github.com/globalsign/mgo/bson"
	"runtime/debug"
)

func InitTaskCleanUserIds() {
	adminUser, err := GetAdminUser()
	if err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return
	}
	tasks, err := model.GetTaskList(nil, 0, constants.Infinite, "+_id")
	if err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return
	}
	for _, t := range tasks {
		if !t.ScheduleId.Valid() {
			t.ScheduleId = bson.ObjectIdHex(constants.ObjectIdNull)
			if err := t.Save(); err != nil {
				log.Errorf(err.Error())
				debug.PrintStack()
				continue
			}
		}

		if !t.UserId.Valid() {
			t.UserId = adminUser.Id
			if err := t.Save(); err != nil {
				log.Errorf(err.Error())
				debug.PrintStack()
				continue
			}
		}
	}
}

func InitProjectCleanUserIds() {
	adminUser, err := GetAdminUser()
	if err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return
	}
	projects, err := model.GetProjectList(nil, "+_id")
	if err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return
	}
	for _, p := range projects {
		if !p.UserId.Valid() {
			p.UserId = adminUser.Id
			if err := p.Save(); err != nil {
				log.Errorf(err.Error())
				debug.PrintStack()
				continue
			}
		}
	}
}

func InitSpiderCleanUserIds() {
	adminUser, err := GetAdminUser()
	if err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return
	}
	spiders, _ := model.GetSpiderAllList(nil)
	for _, s := range spiders {
		if !s.UserId.Valid() {
			s.UserId = adminUser.Id
			if err := s.Save(); err != nil {
				log.Errorf(err.Error())
				debug.PrintStack()
				continue
			}
		}
	}
}

func InitScheduleCleanUserIds() {
	adminUser, err := GetAdminUser()
	if err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return
	}
	schedules, _ := model.GetScheduleList(nil)
	for _, s := range schedules {
		if !s.UserId.Valid() {
			s.UserId = adminUser.Id
			if err := s.Save(); err != nil {
				log.Errorf(err.Error())
				debug.PrintStack()
				continue
			}
		}
	}
}

func InitCleanService() error {
	if model.IsMaster() {
		// 清理任务UserIds
		InitTaskCleanUserIds()
		// 清理项目UserIds
		InitProjectCleanUserIds()
		// 清理爬虫UserIds
		InitSpiderCleanUserIds()
		// 清理定时任务UserIds
		InitScheduleCleanUserIds()
	}
	return nil
}
