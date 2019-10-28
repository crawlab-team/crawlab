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
	NodeKey     string        `json:"node_key" bson:"node_key"`
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

func (sch *Schedule) Delete() error {
	s, c := database.GetCol("schedules")
	defer s.Close()
	return c.RemoveId(sch.Id)
}

func (sch *Schedule) SyncNodeIdAndSpiderId(node Node, spider Spider) {
	sch.syncNodeId(node)
	sch.syncSpiderId(spider)
}

func (sch *Schedule) syncNodeId(node Node) {
	if node.Id.Hex() == sch.NodeId.Hex() {
		return
	}
	sch.NodeId = node.Id
	_ = sch.Save()
}

func (sch *Schedule) syncSpiderId(spider Spider) {
	if spider.Id.Hex() == sch.SpiderId.Hex() {
		return
	}
	sch.SpiderId = spider.Id
	_ = sch.Save()
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
		if schedule.NodeId == bson.ObjectIdHex(constants.ObjectIdNull) {
			// 选择所有节点
			schedule.NodeName = "All Nodes"
		} else {
			// 选择单一节点
			node, err := GetNode(schedule.NodeId)
			if err != nil {
				log.Errorf(err.Error())
				continue
			}
			schedule.NodeName = node.Name
		}

		// 获取爬虫名称
		spider, err := GetSpider(schedule.SpiderId)
		if err != nil {
			log.Errorf("get spider by id: %s, error: %s", schedule.SpiderId.Hex(), err.Error())
			debug.PrintStack()
			_ = schedule.Delete()
			continue
		}
		schedule.SpiderName = spider.Name
		schs = append(schs, schedule)
	}
	return schs, nil
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
	node, err := GetNode(item.NodeId)
	if err != nil {
		return err
	}

	item.NodeKey = node.Key
	if err := item.Save(); err != nil {
		return err
	}
	return nil
}

func AddSchedule(item Schedule) error {
	s, c := database.GetCol("schedules")
	defer s.Close()

	node, err := GetNode(item.NodeId)
	if err != nil {
		return err
	}

	item.Id = bson.NewObjectId()
	item.CreateTs = time.Now()
	item.UpdateTs = time.Now()
	item.NodeKey = node.Key

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
