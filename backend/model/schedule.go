package model

import (
	"crawlab/constants"
	"crawlab/database"
	"crawlab/lib/cron"
	"github.com/apex/log"
	"github.com/globalsign/mgo/bson"
	"runtime/debug"
	"time"
)

type Schedule struct {
	Id          bson.ObjectId `json:"_id" bson:"_id"`
	Name        string        `json:"name" bson:"name"`
	Description string        `json:"description" bson:"description"`
	SpiderId    bson.ObjectId `json:"spider_id" bson:"spider_id"`
	NodeId      bson.ObjectId `json:"node_id" bson:"node_id"`
	Cron        string        `json:"cron" bson:"cron"`
	EntryId     cron.EntryID  `json:"entry_id" bson:"entry_id"`
	Param       string        `json:"param" bson:"param"`

	// 前端展示
	SpiderName string `json:"spider_name" bson:"spider_name"`
	NodeName   string `json:"node_name" bson:"node_name"`

	CreateTs time.Time `json:"create_ts" bson:"create_ts"`
	UpdateTs time.Time `json:"update_ts" bson:"update_ts"`
}

func (sch *Schedule) Save() error {
	s, c := database.GetCol("schedules")
	defer s.Close()
	sch.UpdateTs = time.Now()
	if err := c.UpdateId(sch.Id, sch); err != nil {
		return err
	}
	return nil
}

func GetScheduleList(filter interface{}) ([]Schedule, error) {
	s, c := database.GetCol("schedules")
	defer s.Close()

	var schedules []Schedule
	if err := c.Find(filter).All(&schedules); err != nil {
		return schedules, err
	}

	for i, schedule := range schedules {
		// 获取节点名称
		if schedule.NodeId == bson.ObjectIdHex(constants.ObjectIdNull) {
			// 选择所有节点
			schedules[i].NodeName = "All Nodes"
		} else {
			// 选择单一节点
			node, err := GetNode(schedule.NodeId)
			if err != nil {
				log.Errorf(err.Error())
				continue
			}
			schedules[i].NodeName = node.Name
		}

		// 获取爬虫名称
		spider, err := GetSpider(schedule.SpiderId)
		if err != nil {
			log.Errorf("get spider by id: %s, error: %s", schedule.SpiderId.Hex(), err.Error())
			debug.PrintStack()
			continue
		}
		schedules[i].SpiderName = spider.Name
	}
	return schedules, nil
}

func GetSchedule(id bson.ObjectId) (Schedule, error) {
	s, c := database.GetCol("schedules")
	defer s.Close()

	var result Schedule
	if err := c.FindId(id).One(&result); err != nil {
		return result, err
	}
	return result, nil
}

func UpdateSchedule(id bson.ObjectId, item Schedule) error {
	s, c := database.GetCol("schedules")
	defer s.Close()

	var result Schedule
	if err := c.FindId(id).One(&result); err != nil {
		return err
	}

	if err := item.Save(); err != nil {
		return err
	}
	return nil
}

func AddSchedule(item Schedule) error {
	s, c := database.GetCol("schedules")
	defer s.Close()

	item.Id = bson.NewObjectId()
	item.CreateTs = time.Now()
	item.UpdateTs = time.Now()

	if err := c.Insert(&item); err != nil {
		debug.PrintStack()
		log.Errorf(err.Error())
		return err
	}
	return nil
}

func RemoveSchedule(id bson.ObjectId) error {
	s, c := database.GetCol("schedules")
	defer s.Close()

	var result Schedule
	if err := c.FindId(id).One(&result); err != nil {
		return err
	}

	if err := c.RemoveId(id); err != nil {
		return err
	}

	return nil
}

func GetScheduleCount() (int, error) {
	s, c := database.GetCol("schedules")
	defer s.Close()

	count, err := c.Count()
	if err != nil {
		return 0, err
	}

	return count, nil
}
