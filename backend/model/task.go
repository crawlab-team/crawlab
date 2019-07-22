package model

import (
	"crawlab/database"
	"github.com/apex/log"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"runtime/debug"
	"time"
)

type Task struct {
	Id       string        `json:"_id" bson:"_id"`
	SpiderId bson.ObjectId `json:"spider_id" bson:"spider_id"`
	StartTs  time.Time     `json:"start_ts" bson:"start_ts"`
	FinishTs time.Time     `json:"finish_ts" bson:"finish_ts"`
	Status   string        `json:"status" bson:"status"`
	NodeId   bson.ObjectId `json:"node_id" bson:"node_id"`
	LogPath  string        `json:"log_path" bson:"log_path"`
	Cmd      string        `json:"cmd" bson:"cmd"`
	Error    string        `json:"error" bson:"error"`

	// 前端数据
	SpiderName string `json:"spider_name"`
	NodeName   string `json:"node_name"`
	NumResults int    `json:"num_results"`

	CreateTs time.Time `json:"create_ts" bson:"create_ts"`
	UpdateTs time.Time `json:"update_ts" bson:"update_ts"`
}

func (t *Task) GetSpider() (Spider, error) {
	spider, err := GetSpider(t.SpiderId)
	if err != nil {
		return spider, err
	}
	return spider, nil
}

func (t *Task) GetNode() (Node, error) {
	node, err := GetNode(t.NodeId)
	if err != nil {
		return node, err
	}
	return node, nil
}

func (t *Task) Save() error {
	s, c := database.GetCol("tasks")
	defer s.Close()
	t.UpdateTs = time.Now()
	if err := c.UpdateId(t.Id, t); err != nil {
		debug.PrintStack()
		return err
	}
	return nil
}

func (t *Task) Delete() error {
	s, c := database.GetCol("tasks")
	defer s.Close()
	if err := c.RemoveId(t.Id); err != nil {
		return err
	}
	return nil
}

func (t *Task) GetResults(pageNum int, pageSize int) (results []interface{}, total int, err error) {
	spider, err := t.GetSpider()
	if err != nil {
		log.Errorf(err.Error())
		return
	}

	if spider.Col == "" {
		return
	}

	s, c := database.GetCol(spider.Col)
	defer s.Close()

	query := bson.M{
		"task_id": t.Id,
	}
	if err = c.Find(query).Skip((pageNum - 1) * pageSize).Limit(pageSize).Sort("-create_ts").All(&results); err != nil {
		return
	}

	if total, err = c.Find(query).Count(); err != nil {
		return
	}

	return
}

func GetTaskList(filter interface{}, skip int, limit int, sortKey string) ([]Task, error) {
	s, c := database.GetCol("tasks")
	defer s.Close()

	var tasks []Task
	if err := c.Find(filter).Skip(skip).Limit(limit).Sort(sortKey).All(&tasks); err != nil {
		debug.PrintStack()
		return tasks, err
	}

	for i, task := range tasks {
		// 获取爬虫名称
		spider, err := task.GetSpider()
		if err == mgo.ErrNotFound {
			// do nothing
		} else if err != nil {
			return tasks, err
		} else {
			tasks[i].SpiderName = spider.DisplayName
		}

		// 获取节点名称
		node, err := task.GetNode()
		if err == mgo.ErrNotFound {
			// do nothing
		} else if err != nil {
			return tasks, err
		} else {
			tasks[i].NodeName = node.Name
		}

		// 获取结果数
		if spider.Col == "" {
			continue
		}
		s, c := database.GetCol(spider.Col)
		tasks[i].NumResults, err = c.Find(bson.M{"task_id": task.Id}).Count()
		if err != nil {
			continue
		}
		s.Close()
	}
	return tasks, nil
}

func GetTaskListTotal(filter interface{}) (int, error) {
	s, c := database.GetCol("tasks")
	defer s.Close()

	var result int
	result, err := c.Find(filter).Count()
	if err != nil {
		return result, err
	}
	return result, nil
}

func GetTask(id string) (Task, error) {
	s, c := database.GetCol("tasks")
	defer s.Close()

	var task Task
	if err := c.FindId(id).One(&task); err != nil {
		debug.PrintStack()
		return task, err
	}
	return task, nil
}

func AddTask(item Task) error {
	s, c := database.GetCol("tasks")
	defer s.Close()

	item.CreateTs = time.Now()
	item.UpdateTs = time.Now()

	if err := c.Insert(&item); err != nil {
		return err
	}
	return nil
}

func RemoveTask(id string) error {
	s, c := database.GetCol("tasks")
	defer s.Close()

	var result Task
	if err := c.FindId(id).One(&result); err != nil {
		return err
	}

	if err := c.RemoveId(id); err != nil {
		return err
	}

	return nil
}
