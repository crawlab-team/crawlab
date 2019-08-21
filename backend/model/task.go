package model

import (
	"runtime/debug"
	"time"

	"crawlab/constants"
	"crawlab/database"

	"github.com/apex/log"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type Task struct {
	Id              string        `json:"_id" bson:"_id"`
	SpiderId        bson.ObjectId `json:"spider_id" bson:"spider_id"`
	StartTs         time.Time     `json:"start_ts" bson:"start_ts"`
	FinishTs        time.Time     `json:"finish_ts" bson:"finish_ts"`
	Status          string        `json:"status" bson:"status"`
	NodeId          bson.ObjectId `json:"node_id" bson:"node_id"`
	LogPath         string        `json:"log_path" bson:"log_path"`
	Cmd             string        `json:"cmd" bson:"cmd"`
	Error           string        `json:"error" bson:"error"`
	ResultCount     int           `json:"result_count" bson:"result_count"`
	WaitDuration    float64       `json:"wait_duration" bson:"wait_duration"`
	RuntimeDuration float64       `json:"runtime_duration" bson:"runtime_duration"`
	TotalDuration   float64       `json:"total_duration" bson:"total_duration"`

	// 前端数据
	SpiderName string `json:"spider_name"`
	NodeName   string `json:"node_name"`

	CreateTs time.Time `json:"create_ts" bson:"create_ts"`
	UpdateTs time.Time `json:"update_ts" bson:"update_ts"`
}

type TaskDailyItem struct {
	Date               string  `json:"date" bson:"_id"`
	TaskCount          int     `json:"task_count" bson:"task_count"`
	AvgRuntimeDuration float64 `json:"avg_runtime_duration" bson:"avg_runtime_duration"`
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

func GetTaskCount(query interface{}) (int, error) {
	s, c := database.GetCol("tasks")
	defer s.Close()

	count, err := c.Find(query).Count()
	if err != nil {
		return 0, err
	}

	return count, nil
}

func GetDailyTaskStats(query bson.M) ([]TaskDailyItem, error) {
	s, c := database.GetCol("tasks")
	defer s.Close()

	// 起始日期
	startDate := time.Now().Add(-30 * 24 * time.Hour)
	endDate := time.Now()

	// query
	query["create_ts"] = bson.M{
		"$gte": startDate,
		"$lt":  endDate,
	}

	// match
	op1 := bson.M{
		"$match": query,
	}

	// project
	op2 := bson.M{
		"$project": bson.M{
			"date": bson.M{
				"$dateToString": bson.M{
					"format":   "%Y%m%d",
					"date":     "$create_ts",
					"timezone": "Asia/Shanghai",
				},
			},
			"success_count": bson.M{
				"$cond": []interface{}{
					bson.M{
						"$eq": []string{
							"$status",
							constants.StatusFinished,
						},
					},
					1,
					0,
				},
			},
			"runtime_duration": "$runtime_duration",
		},
	}

	// group
	op3 := bson.M{
		"$group": bson.M{
			"_id":              "$date",
			"task_count":       bson.M{"$sum": 1},
			"runtime_duration": bson.M{"$sum": "$runtime_duration"},
		},
	}

	op4 := bson.M{
		"$project": bson.M{
			"task_count": "$task_count",
			"date":       "$date",
			"avg_runtime_duration": bson.M{
				"$divide": []string{"$runtime_duration", "$task_count"},
			},
		},
	}

	// run aggregation
	var items []TaskDailyItem
	if err := c.Pipe([]bson.M{op1, op2, op3, op4}).All(&items); err != nil {
		return items, err
	}

	// 缓存每日数据
	dict := make(map[string]TaskDailyItem)
	for _, item := range items {
		dict[item.Date] = item
	}

	// 遍历日期
	var dailyItems []TaskDailyItem
	for date := startDate; endDate.Sub(date) > 0; date = date.Add(24 * time.Hour) {
		dateStr := date.Format("20060102")
		dailyItems = append(dailyItems, TaskDailyItem{
			Date:               dateStr,
			TaskCount:          dict[dateStr].TaskCount,
			AvgRuntimeDuration: dict[dateStr].AvgRuntimeDuration,
		})
	}

	return dailyItems, nil
}

func UpdateTaskResultCount(id string) (err error) {
	// 获取任务
	task, err := GetTask(id)
	if err != nil {
		log.Errorf(err.Error())
		return err
	}

	// 获取爬虫
	spider, err := GetSpider(task.SpiderId)
	if err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return err
	}

	// 获取结果数量
	s, c := database.GetCol(spider.Col)
	defer s.Close()
	resultCount, err := c.Find(bson.M{"task_id": task.Id}).Count()
	if err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return err
	}

	// 保存结果数量
	task.ResultCount = resultCount
	if err := task.Save(); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return err
	}
	return nil
}
