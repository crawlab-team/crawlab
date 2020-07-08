package services

import (
	"crawlab/constants"
	"crawlab/database"
	"crawlab/model"
	"crawlab/services/local_node"
	"crawlab/utils"
	"encoding/json"
	"github.com/apex/log"
	"github.com/globalsign/mgo/bson"
	"runtime/debug"
	"time"
)

type Data struct {
	Key      string `json:"key"`
	Mac      string `json:"mac"`
	Ip       string `json:"ip"`
	Hostname string `json:"hostname"`
	Name     string `json:"name"`
	NameType string `json:"name_type"`

	Master       bool      `json:"master"`
	UpdateTs     time.Time `json:"update_ts"`
	UpdateTsUnix int64     `json:"update_ts_unix"`
}

// 所有调用IsMasterNode的方法，都永远会在master节点执行，所以GetCurrentNode方法返回永远是master节点
// 该ID的节点是否为主节点
func IsMasterNode(id string) bool {
	curNode := local_node.CurrentNode()
	//curNode, _ := model.GetCurrentNode()
	node, _ := model.GetNode(bson.ObjectIdHex(id))
	return curNode.Id == node.Id
}

// 获取节点数据
func GetNodeData() (Data, error) {
	localNode := local_node.GetLocalNode()
	key := localNode.Identify
	if key == "" {
		return Data{}, nil
	}

	value, err := database.RedisClient.HGet("nodes", key)
	data := Data{}
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return data, err
	}
	return data, err
}
func GetRedisNode(key string) (*Data, error) {
	// 获取节点数据
	value, err := database.RedisClient.HGet("nodes", key)
	if err != nil {
		log.Errorf(err.Error())
		return nil, err
	}

	// 解析节点列表数据
	var data Data
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		log.Errorf(err.Error())
		return nil, err
	}
	return &data, nil
}

// 更新所有节点状态
func UpdateNodeStatus() {
	// 从Redis获取节点keys
	list, err := database.RedisClient.HScan("nodes")
	if err != nil {
		log.Errorf("get redis node keys error: %s", err.Error())
		return
	}
	var offlineKeys []string
	// 遍历节点keys
	for _, dataStr := range list {
		var data Data
		if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
			log.Errorf(err.Error())
			continue
		}
		// 如果记录的更新时间超过60秒，该节点被认为离线
		if time.Now().Unix()-data.UpdateTsUnix > 60 {
			offlineKeys = append(offlineKeys, data.Key)
			// 在Redis中删除该节点
			if err := database.RedisClient.HDel("nodes", data.Key); err != nil {
				log.Errorf("delete redis node key error:%s, key:%s", err.Error(), data.Key)
			}
			continue
		}

		// 处理node信息
		if err = UpdateNodeInfo(&data); err != nil {
			log.Errorf(err.Error())
			continue
		}
	}
	if len(offlineKeys) > 0 {
		s, c := database.GetCol("nodes")
		defer s.Close()
		_, err = c.UpdateAll(bson.M{
			"key": bson.M{
				"$in": offlineKeys,
			},
		}, bson.M{
			"$set": bson.M{
				"status":         constants.StatusOffline,
				"update_ts":      time.Now(),
				"update_ts_unix": time.Now().Unix(),
			},
		})
		if err != nil {
			log.Errorf(err.Error())
		}
	}
}

// 处理节点信息
func UpdateNodeInfo(data *Data) (err error) {
	// 更新节点信息到数据库
	s, c := database.GetCol("nodes")
	defer s.Close()

	_, err = c.Upsert(bson.M{"key": data.Key}, bson.M{
		"$set": bson.M{
			"status":         constants.StatusOnline,
			"key":            data.Key,
			"name_type":      data.NameType,
			"ip":             data.Ip,
			"port":           "8000",
			"mac":            data.Mac,
			"is_master":      data.Master,
			"update_ts":      time.Now(),
			"update_ts_unix": time.Now().Unix(),
		},
		"$setOnInsert": bson.M{
			"name": data.Name,
			"_id":  bson.NewObjectId(),
		},
	})
	return err
}

// 更新节点数据
func UpdateNodeData() {
	localNode := local_node.GetLocalNode()
	key := localNode.Identify
	// 构造节点数据
	data := Data{
		Key:          key,
		Mac:          localNode.Mac,
		Ip:           localNode.Ip,
		Hostname:     localNode.Hostname,
		Name:         localNode.Identify,
		NameType:     string(localNode.IdentifyType),
		Master:       model.IsMaster(),
		UpdateTs:     time.Now(),
		UpdateTsUnix: time.Now().Unix(),
	}

	// 注册节点到Redis
	dataBytes, err := json.Marshal(&data)
	if err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return
	}

	if err := database.RedisClient.HSet("nodes", key, utils.BytesToString(dataBytes)); err != nil {
		log.Errorf(err.Error())
		return
	}
}

// 发送心跳信息到Redis，每5秒发送一次
func SendHeartBeat() {
	for {
		UpdateNodeData()
		time.Sleep(5 * time.Second)
	}
}

// 每10秒刷新一次节点信息
func UpdateNodeStatusPeriodically() {
	for {
		UpdateNodeStatus()
		time.Sleep(10 * time.Second)
	}
}

// 初始化节点服务
func InitNodeService() error {
	node, err := local_node.InitLocalNode()
	if err != nil {
		return err
	}

	// 每5秒更新一次本节点信息
	go SendHeartBeat()

	// 首次更新节点数据（注册到Redis）
	UpdateNodeData()
	if model.IsMaster() {
		err = model.UpdateMasterNodeInfo(node.Identify, node.Ip, node.Mac, node.Hostname)
		if err != nil {
			return err
		}
	}

	// 节点准备完毕
	if err = node.Ready(); err != nil {
		return err
	}

	// 如果为主节点，每10秒刷新所有节点信息
	if model.IsMaster() {
		go UpdateNodeStatusPeriodically()
	}

	// 更新在当前节点执行中的任务状态为：abnormal
	if err := model.UpdateTaskToAbnormal(node.Current().Id); err != nil {
		debug.PrintStack()
		return err
	}

	return nil
}
