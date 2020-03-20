package model

import (
	"crawlab/constants"
	"crawlab/database"
	"github.com/apex/log"
	"github.com/globalsign/mgo/bson"
	"runtime/debug"
	"time"
)

type Project struct {
	Id          bson.ObjectId `json:"_id" bson:"_id"`
	Name        string        `json:"name" bson:"name"`
	Description string        `json:"description" bson:"description"`
	Tags        []string      `json:"tags" bson:"tags"`

	// 前端展示
	Spiders  []Spider `json:"spiders" bson:"spiders"`
	Username string   `json:"username" bson:"username"`

	UserId   bson.ObjectId `json:"user_id" bson:"user_id"`
	CreateTs time.Time     `json:"create_ts" bson:"create_ts"`
	UpdateTs time.Time     `json:"update_ts" bson:"update_ts"`
}

func (p *Project) Save() error {
	s, c := database.GetCol("projects")
	defer s.Close()

	p.UpdateTs = time.Now()

	if err := c.UpdateId(p.Id, p); err != nil {
		debug.PrintStack()
		return err
	}
	return nil
}

func (p *Project) Add() error {
	s, c := database.GetCol("projects")
	defer s.Close()

	p.Id = bson.NewObjectId()
	p.UpdateTs = time.Now()
	p.CreateTs = time.Now()
	if err := c.Insert(p); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return err
	}

	return nil
}

func (p *Project) GetSpiders() ([]Spider, error) {
	s, c := database.GetCol("spiders")
	defer s.Close()

	var query interface{}
	if p.Id.Hex() == constants.ObjectIdNull {
		query = bson.M{
			"$or": []bson.M{
				{"project_id": p.Id},
				{"project_id": bson.M{"$exists": false}},
			},
		}
	} else {
		query = bson.M{"project_id": p.Id}
	}

	var spiders []Spider
	if err := c.Find(query).All(&spiders); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return spiders, err
	}

	return spiders, nil
}

func GetProject(id bson.ObjectId) (Project, error) {
	s, c := database.GetCol("projects")
	defer s.Close()
	var p Project
	if err := c.Find(bson.M{"_id": id}).One(&p); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return p, err
	}
	return p, nil
}

func GetProjectList(filter interface{}, sortKey string) ([]Project, error) {
	s, c := database.GetCol("projects")
	defer s.Close()

	var projects []Project
	if err := c.Find(filter).Sort(sortKey).All(&projects); err != nil {
		debug.PrintStack()
		return projects, err
	}

	for i, p := range projects {
		// 获取用户名称
		user, _ := GetUser(p.UserId)
		projects[i].Username = user.Username
	}
	return projects, nil
}

func GetProjectListTotal(filter interface{}) (int, error) {
	s, c := database.GetCol("projects")
	defer s.Close()

	var result int
	result, err := c.Find(filter).Count()
	if err != nil {
		return result, err
	}
	return result, nil
}

func UpdateProject(id bson.ObjectId, item Project) error {
	s, c := database.GetCol("projects")
	defer s.Close()

	var result Project
	if err := c.FindId(id).One(&result); err != nil {
		debug.PrintStack()
		return err
	}

	if err := item.Save(); err != nil {
		return err
	}
	return nil
}

func RemoveProject(id bson.ObjectId) error {
	s, c := database.GetCol("projects")
	defer s.Close()

	var result User
	if err := c.FindId(id).One(&result); err != nil {
		return err
	}

	if err := c.RemoveId(id); err != nil {
		return err
	}

	return nil
}
