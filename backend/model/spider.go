package model

import (
	"crawlab/database"
	"crawlab/entity"
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
	Type        string        `json:"type" bson:"type"`                 // 爬虫类别
	FileId      bson.ObjectId `json:"file_id" bson:"file_id"`           // GridFS文件ID
	Col         string        `json:"col" bson:"col"`                   // 结果储存位置
	Site        string        `json:"site" bson:"site"`                 // 爬虫网站
	Envs        []Env         `json:"envs" bson:"envs"`                 // 环境变量
	Remark      string        `json:"remark" bson:"remark"`             // 备注

	// 自定义爬虫
	Src string `json:"src" bson:"src"` // 源码位置
	Cmd string `json:"cmd" bson:"cmd"` // 执行命令

	// 前端展示
	LastRunTs  time.Time `json:"last_run_ts"` // 最后一次执行时间
	LastStatus string    `json:"last_status"` // 最后执行状态

	// 时间
	CreateTs time.Time `json:"create_ts" bson:"create_ts"`
	UpdateTs time.Time `json:"update_ts" bson:"update_ts"`
}

// 更新爬虫
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

// 新增爬虫
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

// 获取爬虫的任务
func (spider *Spider) GetTasks() ([]Task, error) {
	tasks, err := GetTaskList(bson.M{"spider_id": spider.Id}, 0, 10, "-create_ts")
	if err != nil {
		return tasks, err
	}
	return tasks, nil
}

// 爬虫最新的任务
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

// 删除爬虫
func (spider *Spider) Delete() error {
	s, c := database.GetCol("spiders")
	defer s.Close()
	return c.RemoveId(spider.Id)
}

// 获取爬虫列表
func GetSpiderList(filter interface{}, skip int, limit int) ([]Spider, int, error) {
	s, c := database.GetCol("spiders")
	defer s.Close()

	// 获取爬虫列表
	var spiders []Spider
	if err := c.Find(filter).Skip(skip).Limit(limit).Sort("+name").All(&spiders); err != nil {
		debug.PrintStack()
		return spiders, 0, err
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

	count, _ := c.Find(filter).Count()

	return spiders, count, nil
}

// 获取爬虫(根据FileId)
func GetSpiderByFileId(fileId bson.ObjectId) *Spider {
	s, c := database.GetCol("spiders")
	defer s.Close()

	var result *Spider
	if err := c.Find(bson.M{"file_id": fileId}).One(&result); err != nil {
		log.Errorf("get spider error: %s, file_id: %s", err.Error(), fileId.Hex())
		debug.PrintStack()
		return nil
	}
	return result
}

// 获取爬虫(根据名称)
func GetSpiderByName(name string) *Spider {
	s, c := database.GetCol("spiders")
	defer s.Close()

	var result *Spider
	if err := c.Find(bson.M{"name": name}).One(&result); err != nil {
		log.Errorf("get spider error: %s, spider_name: %s", err.Error(), name)
		//debug.PrintStack()
		return nil
	}
	return result
}

// 获取爬虫(根据ID)
func GetSpider(id bson.ObjectId) (Spider, error) {
	s, c := database.GetCol("spiders")
	defer s.Close()

	var result Spider
	if err := c.FindId(id).One(&result); err != nil {
		if err != mgo.ErrNotFound {
			log.Errorf("get spider error: %s, id: %id", err.Error(), id.Hex())
			debug.PrintStack()
		}
		return result, err
	}
	return result, nil
}

// 更新爬虫
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

// 删除爬虫
func RemoveSpider(id bson.ObjectId) error {
	s, c := database.GetCol("spiders")
	defer s.Close()

	var result Spider
	if err := c.FindId(id).One(&result); err != nil {
		return err
	}

	if err := c.RemoveId(id); err != nil {
		log.Errorf("remove spider error: %s, id:%s", err.Error(), id.Hex())
		debug.PrintStack()
		return err
	}

	// gf上的文件
	s, gf := database.GetGridFs("files")
	defer s.Close()

	if err := gf.RemoveId(result.FileId); err != nil {
		log.Error("remove file error, id:" + result.FileId.Hex())
		debug.PrintStack()
		return err
	}

	return nil
}

// 删除所有爬虫
func RemoveAllSpider() error {
	s, c := database.GetCol("spiders")
	defer s.Close()

	var spiders []Spider
	err := c.Find(nil).All(&spiders)
	if err != nil {
		log.Error("get all spiders error:" + err.Error())
		return err
	}
	for _, spider := range spiders {
		if err := RemoveSpider(spider.Id); err != nil {
			log.Error("remove spider error:" + err.Error())
		}
	}
	return nil
}

// 获取爬虫总数
func GetSpiderCount() (int, error) {
	s, c := database.GetCol("spiders")
	defer s.Close()

	count, err := c.Count()
	if err != nil {
		return 0, err
	}
	return count, nil
}

// 获取爬虫类型
func GetSpiderTypes() ([]*entity.SpiderType, error) {
	s, c := database.GetCol("spiders")
	defer s.Close()

	group := bson.M{
		"$group": bson.M{
			"_id":   "$type",
			"count": bson.M{"$sum": 1},
		},
	}
	var types []*entity.SpiderType
	if err := c.Pipe([]bson.M{group}).All(&types); err != nil {
		log.Errorf("get spider types error: %s", err.Error())
		debug.PrintStack()
		return nil, err
	}

	return types, nil
}
