package model

import (
	"crawlab/database"
	"github.com/apex/log"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"runtime/debug"
	"time"
)

type Env struct {
	Name  string `json:"name" bson:"name"`
	Value string `json:"value" bson:"value"`
}

type Spider struct {
	Id          bson.ObjectId `json:"_id" bson:"_id"`                   // 爬虫ID
	Name        string        `json:"name" bson:"name"`                 // 爬虫名称（唯一）
	DisplayName string        `json:"display_name" bson:"display_name"` // 爬虫显示名称
	Type        string        `json:"type"`                             // 爬虫类别
	FileId      bson.ObjectId `json:"file_id" bson:"file_id"`           // GridFS文件ID
	Col         string        `json:"col"`                              // 结果储存位置
	Site        string        `json:"site"`                             // 爬虫网站
	Envs        []Env         `json:"envs" bson:"envs"`                 // 环境变量
	Remark      string        `json:"remark"`                           // 备注
	// 自定义爬虫
	Src string `json:"src" bson:"src"` // 源码位置
	Cmd string `json:"cmd" bson:"cmd"` // 执行命令

	// 前端展示
	LastRunTs  time.Time `json:"last_run_ts"` // 最后一次执行时间
	LastStatus string    `json:"last_status"` // 最后执行状态

	// TODO: 可配置爬虫
	//Fields                 []interface{} `json:"fields"`
	//DetailFields           []interface{} `json:"detail_fields"`
	//CrawlType              string        `json:"crawl_type"`
	//StartUrl               string        `json:"start_url"`
	//UrlPattern             string        `json:"url_pattern"`
	//ItemSelector           string        `json:"item_selector"`
	//ItemSelectorType       string        `json:"item_selector_type"`
	//PaginationSelector     string        `json:"pagination_selector"`
	//PaginationSelectorType string        `json:"pagination_selector_type"`

	CreateTs time.Time `json:"create_ts" bson:"create_ts"`
	UpdateTs time.Time `json:"update_ts" bson:"update_ts"`
}

func (spider *Spider) Save() error {
	s, c := database.GetCol("spiders")
	defer s.Close()

	spider.UpdateTs = time.Now()

	if err := c.UpdateId(spider.Id, spider); err != nil {
		debug.PrintStack()
		return err
	}
	return nil
}

func (spider *Spider) Add() error {
	s, c := database.GetCol("spiders")
	defer s.Close()

	spider.Id = bson.NewObjectId()
	spider.CreateTs = time.Now()
	spider.UpdateTs = time.Now()

	if err := c.Insert(&spider); err != nil {
		return err
	}
	return nil
}

func (spider *Spider) GetTasks() ([]Task, error) {
	tasks, err := GetTaskList(bson.M{"spider_id": spider.Id}, 0, 10, "-create_ts")
	if err != nil {
		return tasks, err
	}
	return tasks, nil
}

func (spider *Spider) GetLastTask() (Task, error) {
	tasks, err := GetTaskList(bson.M{"spider_id": spider.Id}, 0, 1, "-create_ts")
	if err != nil {
		return Task{}, err
	}
	if tasks == nil {
		return Task{}, nil
	}
	return tasks[0], nil
}

func GetSpiderList(filter interface{}, skip int, limit int) ([]Spider, error) {
	s, c := database.GetCol("spiders")
	defer s.Close()

	// 获取爬虫列表
	spiders := []Spider{}
	if err := c.Find(filter).Skip(skip).Limit(limit).Sort("name asc").All(&spiders); err != nil {
		debug.PrintStack()
		return spiders, err
	}

	// 遍历爬虫列表
	for i, spider := range spiders {
		// 获取最后一次任务
		task, err := spider.GetLastTask()
		if err != nil {
			log.Errorf(err.Error())
			debug.PrintStack()
			continue
		}

		// 赋值
		spiders[i].LastRunTs = task.CreateTs
		spiders[i].LastStatus = task.Status
	}

	return spiders, nil
}

func GetSpider(id bson.ObjectId) (Spider, error) {
	s, c := database.GetCol("spiders")
	defer s.Close()

	var result Spider
	if err := c.FindId(id).One(&result); err != nil {
		if err != mgo.ErrNotFound {
			debug.PrintStack()
		}
		return result, err
	}
	return result, nil
}

func UpdateSpider(id bson.ObjectId, item Spider) error {
	s, c := database.GetCol("spiders")
	defer s.Close()

	var result Spider
	if err := c.FindId(id).One(&result); err != nil {
		debug.PrintStack()
		return err
	}

	if err := item.Save(); err != nil {
		return err
	}
	return nil
}

func RemoveSpider(id bson.ObjectId) error {
	s, c := database.GetCol("spiders")
	defer s.Close()

	var result Spider
	if err := c.FindId(id).One(&result); err != nil {
		return err
	}

	if err := c.RemoveId(id); err != nil {
		return err
	}

	// gf上的文件
	s, gf := database.GetGridFs("files")
	defer s.Close()

	if err := gf.RemoveId(result.FileId); err != nil {
		log.Error("remove file error, id:" + result.FileId.Hex())
		return err
	}

	return nil
}

func GetSpiderCount() (int, error) {
	s, c := database.GetCol("spiders")
	defer s.Close()

	count, err := c.Count()
	if err != nil {
		return 0, err
	}

	return count, nil
}
