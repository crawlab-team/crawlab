package model

import (
	"crawlab/constants"
	"crawlab/database"
	"crawlab/lib/cron"
	"github.com/apex/log"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"runtime/debug"
	"time"
)

type Schedule struct {
	Id             bson.ObjectId   `json:"_id" bson:"_id"`
	Name           string          `json:"name" bson:"name"`
	Description    string          `json:"description" bson:"description"`
	SpiderId       bson.ObjectId   `json:"spider_id" bson:"spider_id"`
	Cron           string          `json:"cron" bson:"cron"`
	EntryId        cron.EntryID    `json:"entry_id" bson:"entry_id"`
	Param          string          `json:"param" bson:"param"`
	RunType        string          `json:"run_type" bson:"run_type"`
	NodeIds        []bson.ObjectId `json:"node_ids" bson:"node_ids"`
	Status         string          `json:"status" bson:"status"`
	Enabled        bool            `json:"enabled" bson:"enabled"`
	UserId         bson.ObjectId   `json:"user_id" bson:"user_id"`
	ScrapySpider   string          `json:"scrapy_spider" bson:"scrapy_spider"`
	ScrapyLogLevel string          `json:"scrapy_log_level" bson:"scrapy_log_level"`

	// 前端展示
	SpiderName string `json:"spider_name" bson:"spider_name"`
	Username   string `json:"user_name" bson:"user_name"`
	Nodes      []Node `json:"nodes" bson:"nodes"`
	Message    string `json:"message" bson:"message"`

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

func (sch *Schedule) Delete() error {
	s, c := database.GetCol("schedules")
	defer s.Close()
	return c.RemoveId(sch.Id)
}

func GetScheduleList(filter interface{}) ([]Schedule, error) {
	s, c := database.GetCol("schedules")
	defer s.Close()

	var schedules []Schedule
	if err := c.Find(filter).All(&schedules); err != nil {
		return schedules, err
	}

	var schs []Schedule
	for _, schedule := range schedules {
		// 获取节点名称
		schedule.Nodes = []Node{}
		if schedule.RunType == constants.RunTypeSelectedNodes {
			for _, nodeId := range schedule.NodeIds {
				// 选择单一节点
				node, _ := GetNode(nodeId)
				schedule.Nodes = append(schedule.Nodes, node)
			}
		}

		// 获取爬虫名称
		spider, err := GetSpider(schedule.SpiderId)
		if err != nil {
			log.Errorf("get spider by id: %s, error: %s", schedule.SpiderId.Hex(), err.Error())
			schedule.Status = constants.ScheduleStatusError
			if err == mgo.ErrNotFound {
				schedule.Message = constants.ScheduleStatusErrorNotFoundSpider
			} else {
				schedule.Message = err.Error()
			}
		} else {
			schedule.SpiderName = spider.Name
		}

		// 获取用户名称
		user, _ := GetUser(schedule.UserId)
		schedule.Username = user.Username

		schs = append(schs, schedule)
	}
	return schs, nil
}

func GetSchedule(id bson.ObjectId) (Schedule, error) {
	s, c := database.GetCol("schedules")
	defer s.Close()

	var schedule Schedule
	if err := c.FindId(id).One(&schedule); err != nil {
		return schedule, err
	}

	// 获取用户名称
	user, _ := GetUser(schedule.UserId)
	schedule.Username = user.Username

	return schedule, nil
}

func UpdateSchedule(id bson.ObjectId, item Schedule) error {
	s, c := database.GetCol("schedules")
	defer s.Close()

	var result Schedule
	if err := c.FindId(id).One(&result); err != nil {
		return err
	}

	item.UpdateTs = time.Now()
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

func GetScheduleCount(filter interface{}) (int, error) {
	s, c := database.GetCol("schedules")
	defer s.Close()

	count, err := c.Find(filter).Count()
	if err != nil {
		return 0, err
	}

	return count, nil
}
