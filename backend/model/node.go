package model

import (
	"crawlab/constants"
	"crawlab/database"
	"crawlab/errors"
	"crawlab/services/local_machine_info"
	"github.com/apex/log"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/spf13/viper"
	"runtime/debug"
	"time"
)

type Node struct {
	Id          bson.ObjectId                  `json:"_id" bson:"_id"`
	Name        string                         `json:"name" bson:"name"`
	Status      string                         `json:"status" bson:"status"`
	Ip          string                         `json:"ip" bson:"ip"`
	Port        string                         `json:"port" bson:"port"`
	Mac         string                         `json:"mac" bson:"mac"`
	Hostname    string                         `json:"hostname" bson:"hostname"`
	Description string                         `json:"description" bson:"description"`
	Usage       local_machine_info.MachineInfo `json:"usage" bson:"usage"`
	// 用于唯一标识节点，可能是mac地址，可能是ip地址
	Key string `json:"key" bson:"key"`

	// 前端展示
	IsMaster bool `json:"is_master" bson:"is_master"`

	UpdateTs     time.Time `json:"update_ts" bson:"update_ts"`
	CreateTs     time.Time `json:"create_ts" bson:"create_ts"`
	UpdateTsUnix int64     `json:"update_ts_unix" bson:"update_ts_unix"`
}

const (
	Yes = "Y"
)

// 当前节点是否为主节点
func IsMaster() bool {
	return viper.GetString("server.master") == Yes
}

func (n *Node) Save() error {
	s, c := database.GetCol("nodes")
	defer s.Close()
	n.UpdateTs = time.Now()
	if err := c.UpdateId(n.Id, n); err != nil {
		return err
	}
	return nil
}

func (n *Node) Add() error {
	s, c := database.GetCol("nodes")
	defer s.Close()
	n.Id = bson.NewObjectId()
	n.UpdateTs = time.Now()
	n.UpdateTsUnix = time.Now().Unix()
	n.CreateTs = time.Now()
	if err := c.Insert(&n); err != nil {
		debug.PrintStack()
		return err
	}
	return nil
}

func (n *Node) Delete() error {
	s, c := database.GetCol("nodes")
	defer s.Close()
	if err := c.RemoveId(n.Id); err != nil {
		debug.PrintStack()
		return err
	}
	return nil
}

func (n *Node) GetTasks() ([]Task, error) {
	tasks, err := GetTaskList(bson.M{"node_id": n.Id}, 0, 10, "-create_ts")
	//tasks, err := GetTaskList(nil, 0, 10, "-create_ts")
	if err != nil {
		debug.PrintStack()
		return []Task{}, err
	}

	return tasks, nil
}

// 节点列表
func GetNodeList(filter interface{}) ([]Node, error) {
	s, c := database.GetCol("nodes")
	defer s.Close()

	var results []Node
	if err := c.Find(filter).All(&results); err != nil {
		log.Error("get node list error: " + err.Error())
		debug.PrintStack()
		return results, err
	}
	return results, nil
}

// 节点信息
func GetNode(id bson.ObjectId) (Node, error) {
	var node Node

	if id.Hex() == "" {
		log.Infof("id is empty")
		debug.PrintStack()
		return node, errors.ErrIpEmpty
	}

	s, c := database.GetCol("nodes")
	defer s.Close()

	if err := c.FindId(id).One(&node); err != nil {
		//log.Errorf("get node error: %s, id: %s", err.Error(), id.Hex())
		//debug.PrintStack()
		return node, err
	}
	return node, nil
}

// 节点信息
func GetNodeByKey(key string) (Node, error) {
	s, c := database.GetCol("nodes")
	defer s.Close()

	var node Node
	if err := c.Find(bson.M{"key": key}).One(&node); err != nil {
		if err != mgo.ErrNotFound {
			log.Errorf(err.Error())
			debug.PrintStack()
		}
		return node, err
	}
	return node, nil
}

// 更新节点
func UpdateNode(id bson.ObjectId, item Node) error {
	s, c := database.GetCol("nodes")
	defer s.Close()

	var node Node
	if err := c.FindId(id).One(&node); err != nil {
		return err
	}

	if err := item.Save(); err != nil {
		return err
	}
	return nil
}

// 任务列表
func GetNodeTaskList(id bson.ObjectId) ([]Task, error) {
	node, err := GetNode(id)
	if err != nil {
		return []Task{}, err
	}
	tasks, err := node.GetTasks()
	if err != nil {
		return []Task{}, err
	}
	return tasks, nil
}

// 节点数
func GetNodeCount(query interface{}) (int, error) {
	s, c := database.GetCol("nodes")
	defer s.Close()

	count, err := c.Find(query).Count()
	if err != nil {
		return 0, err
	}

	return count, nil
}

// 根据redis的key值，重置node节点为offline
func ResetNodeStatusToOffline(list []string) {
	nodes, _ := GetNodeList(nil)
	for _, node := range nodes {
		hasNode := false
		for _, key := range list {
			if key == node.Key {
				hasNode = true
				break
			}
		}
		if !hasNode || node.Status == "" {
			node.Status = constants.StatusOffline
			if err := node.Save(); err != nil {
				log.Errorf(err.Error())
				return
			}
			continue
		}
	}
}

func UpdateMasterNodeInfo(key string, ip string, mac string, hostname string) error {
	s, c := database.GetCol("nodes")
	defer s.Close()
	c.UpdateAll(bson.M{
		"is_master": true,
	}, bson.M{
		"is_master": false,
	})
	_, err := c.Upsert(bson.M{
		"key": key,
	}, bson.M{
		"$set": bson.M{
			"ip":             ip,
			"port":           "8000",
			"mac":            mac,
			"hostname":       hostname,
			"is_master":      true,
			"update_ts":      time.Now(),
			"update_ts_unix": time.Now().Unix(),
		},
		"$setOnInsert": bson.M{
			"key":       key,
			"name":      key,
			"create_ts": time.Now(),
			"_id":       bson.NewObjectId(),
		},
	})
	return err
}

func ChoiceNodes(nodes []Node) (Node, error) {
	// 挑选出最有节点
	// 无可用节点、返回异常

	// 挑选指标--加权平均
	// (x * y * z) / (x1 * y1 * z1)
	var betterRate = 0.0
	var cpuW = 3.0
	var dhW = 1.0
	var memW = 1.0
	betterNode := Node{}

	for _, node := range nodes {
		if node.Status != "online" {
			continue
		}
		if viper.GetString("setting.runOnMaster") == "N" && node.IsMaster {
			continue
		}
		cpu := (node.Usage.LoadAvg1 / float64(node.Usage.CpuNums)) * 2
		dhUse := node.Usage.DhUsed
		memUse := node.Usage.MemUsed
		curRate := cpu*cpuW + dhUse*dhW + memW*memUse
		if betterRate == 0 {
			betterRate = curRate
			betterNode = node
			continue
		}
		if curRate < betterRate {
			betterRate = curRate
			betterNode = node
		}
	}
	if betterNode == (Node{}) {
		return betterNode, errors.ErrGotNodeEmpty
	}
	return betterNode, nil
}
