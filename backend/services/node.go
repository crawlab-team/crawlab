package services

import (
	"crawlab/constants"
	"crawlab/database"
	"crawlab/lib/cron"
	"crawlab/model"
	"encoding/json"
	"fmt"
	"github.com/apex/log"
	"github.com/globalsign/mgo/bson"
	"github.com/spf13/viper"
	"net"
	"runtime/debug"
	"time"
)

type Data struct {
	Mac          string    `json:"mac"`
	Ip           string    `json:"ip"`
	Master       bool      `json:"master"`
	UpdateTs     time.Time `json:"update_ts"`
	UpdateTsUnix int64     `json:"update_ts_unix"`
}

type NodeMessage struct {
	// 通信类别
	Type string `json:"type"`

	// 任务相关
	TaskId string `json:"task_id"` // 任务ID

	// 节点相关
	NodeId string `json:"node_id"` // 节点ID

	// 日志相关
	LogPath string `json:"log_path"` // 日志路径
	Log     string `json:"log"`      // 日志

	// 系统信息
	SysInfo model.SystemInfo `json:"sys_info"`

	// 错误相关
	Error string `json:"error"`
}

const (
	Yes = "Y"
	No  = "N"
)

// 获取本机的IP地址
// TODO: 考虑多个IP地址的情况
func GetIp() (string, error) {
	addrList, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, value := range addrList {
		if ipNet, ok := value.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}
	return "", nil
}

// 获取本机的MAC地址
func GetMac() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		debug.PrintStack()
		return "", err
	}
	for _, inter := range interfaces {
		if inter.HardwareAddr != nil {
			mac := inter.HardwareAddr.String()
			return mac, nil
		}
	}
	return "", nil
}

// 获取本机节点
func GetCurrentNode() (model.Node, error) {
	// 获取本机MAC地址
	mac, err := GetMac()
	if err != nil {
		debug.PrintStack()
		return model.Node{}, err
	}

	// 从数据库中获取当前节点
	var node model.Node
	errNum := 0
	for {
		// 如果错误次数超过10次
		if errNum >= 10 {
			panic("cannot get current node")
		}

		// 尝试获取节点
		node, err = model.GetNodeByMac(mac)

		// 如果获取失败
		if err != nil {
			// 如果为主节点，表示为第一次注册，插入节点信息
			if IsMaster() {
				// 获取本机IP地址
				ip, err := GetIp()
				if err != nil {
					debug.PrintStack()
					return model.Node{}, err
				}
				// 生成节点
				node = model.Node{
					Id:       bson.NewObjectId(),
					Ip:       ip,
					Name:     mac,
					Mac:      mac,
					IsMaster: true,
				}
				if err := node.Add(); err != nil {
					return node, err
				}
				return node, nil
			}
			// 增加错误次数
			errNum++

			// 5秒后重试
			time.Sleep(5 * time.Second)
			continue
		}

		// 跳出循环
		break
	}

	return node, nil
}

// 当前节点是否为主节点
func IsMaster() bool {
	return viper.GetString("server.master") == Yes
}

// 该ID的节点是否为主节点
func IsMasterNode(id string) bool {
	curNode, _ := GetCurrentNode()
	node, _ := model.GetNode(bson.ObjectIdHex(id))
	return curNode.Id == node.Id
}

// 获取节点数据
func GetNodeData() (Data, error) {
	mac, err := GetMac()
	if err != nil {
		return Data{}, err
	}

	value, err := database.RedisClient.HGet("nodes", mac)
	data := Data{}
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return data, err
	}
	return data, err
}

// 更新所有节点状态
func UpdateNodeStatus() {
	// 从Redis获取节点keys
	list, err := database.RedisClient.HKeys("nodes")
	if err != nil {
		log.Errorf(err.Error())
		return
	}

	// 遍历节点keys
	for _, mac := range list {
		// 获取节点数据
		value, err := database.RedisClient.HGet("nodes", mac)
		if err != nil {
			log.Errorf(err.Error())
			return
		}

		// 解析节点列表数据
		var data Data
		if err := json.Unmarshal([]byte(value), &data); err != nil {
			log.Errorf(err.Error())
			return
		}

		// 如果记录的更新时间超过60秒，该节点被认为离线
		if time.Now().Unix()-data.UpdateTsUnix > 60 {
			// 在Redis中删除该节点
			if err := database.RedisClient.HDel("nodes", data.Mac); err != nil {
				log.Errorf(err.Error())
				return
			}

			// 在MongoDB中该节点设置状态为离线
			s, c := database.GetCol("nodes")
			defer s.Close()
			var node model.Node
			if err := c.Find(bson.M{"mac": mac}).One(&node); err != nil {
				log.Errorf(err.Error())
				debug.PrintStack()
				return
			}
			node.Status = constants.StatusOffline
			if err := node.Save(); err != nil {
				log.Errorf(err.Error())
				return
			}
			continue
		}

		// 更新节点信息到数据库
		s, c := database.GetCol("nodes")
		defer s.Close()
		var node model.Node
		if err := c.Find(bson.M{"mac": mac}).One(&node); err != nil {
			// 数据库不存在该节点
			node = model.Node{
				Name:   data.Mac,
				Ip:     data.Ip,
				Port:   "8000",
				Mac:    data.Mac,
				Status: constants.StatusOnline,
			}
			if err := node.Add(); err != nil {
				log.Errorf(err.Error())
				return
			}
		} else {
			// 数据库存在该节点
			node.Status = constants.StatusOnline
			if err := node.Save(); err != nil {
				log.Errorf(err.Error())
				return
			}
		}
	}

	// 遍历数据库中的节点列表
	nodes, err := model.GetNodeList(nil)
	for _, node := range nodes {
		hasNode := false
		for _, mac := range list {
			if mac == node.Mac {
				hasNode = true
				break
			}
		}
		if !hasNode {
			node.Status = constants.StatusOffline
			if err := node.Save(); err != nil {
				log.Errorf(err.Error())
				return
			}
			continue
		}
	}
}

// 更新节点数据
func UpdateNodeData() {
	// 获取MAC地址
	mac, err := GetMac()
	if err != nil {
		log.Errorf(err.Error())
		return
	}

	// 获取IP地址
	ip, err := GetIp()
	if err != nil {
		log.Errorf(err.Error())
		return
	}

	// 构造节点数据
	data := Data{
		Mac:          mac,
		Ip:           ip,
		Master:       IsMaster(),
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
	if err := database.RedisClient.HSet("nodes", mac, string(dataBytes)); err != nil {
		log.Errorf(err.Error())
		return
	}
}

func MasterNodeCallback(channel string, msgStr string) {
	// 反序列化
	var msg NodeMessage
	if err := json.Unmarshal([]byte(msgStr), &msg); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return
	}

	if msg.Type == constants.MsgTypeGetLog {
		// 获取日志
		fmt.Println(msg)
		time.Sleep(10 * time.Millisecond)
		ch := TaskLogChanMap.ChanBlocked(msg.TaskId)
		ch <- msg.Log
	} else if msg.Type == constants.MsgTypeGetSystemInfo {
		// 获取系统信息
		fmt.Println(msg)
		time.Sleep(10 * time.Millisecond)
		ch := SystemInfoChanMap.ChanBlocked(msg.NodeId)
		sysInfoBytes, _ := json.Marshal(&msg.SysInfo)
		ch <- string(sysInfoBytes)
	}
}

func WorkerNodeCallback(channel string, msgStr string) {
	// 反序列化
	msg := NodeMessage{}
	fmt.Println(msgStr)
	if err := json.Unmarshal([]byte(msgStr), &msg); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return
	}

	if msg.Type == constants.MsgTypeGetLog {
		// 消息类型为获取日志

		// 发出的消息
		msgSd := NodeMessage{
			Type:   constants.MsgTypeGetLog,
			TaskId: msg.TaskId,
		}

		// 获取本地日志
		logStr, err := GetLocalLog(msg.LogPath)
		if err != nil {
			log.Errorf(err.Error())
			debug.PrintStack()
			msgSd.Error = err.Error()
		}
		msgSd.Log = string(logStr)

		// 序列化
		msgSdBytes, err := json.Marshal(&msgSd)
		if err != nil {
			log.Errorf(err.Error())
			debug.PrintStack()
			return
		}

		// 发布消息给主节点
		fmt.Println(msgSd)
		if err := database.Publish("nodes:master", string(msgSdBytes)); err != nil {
			log.Errorf(err.Error())
			return
		}
	} else if msg.Type == constants.MsgTypeCancelTask {
		// 取消任务
		ch := TaskExecChanMap.ChanBlocked(msg.TaskId)
		ch <- constants.TaskCancel
	} else if msg.Type == constants.MsgTypeGetSystemInfo {
		// 获取环境信息
		sysInfo, err := GetLocalSystemInfo()
		if err != nil {
			log.Errorf(err.Error())
			return
		}
		msgSd := NodeMessage{
			Type:    constants.MsgTypeGetSystemInfo,
			NodeId:  msg.NodeId,
			SysInfo: sysInfo,
		}
		msgSdBytes, err := json.Marshal(&msgSd)
		if err != nil {
			log.Errorf(err.Error())
			debug.PrintStack()
			return
		}
		fmt.Println(msgSd)
		if err := database.Publish("nodes:master", string(msgSdBytes)); err != nil {
			log.Errorf(err.Error())
			return
		}
	}
}

// 初始化节点服务
func InitNodeService() error {
	// 构造定时任务
	c := cron.New(cron.WithSeconds())

	// 每5秒更新一次本节点信息
	spec := "0/5 * * * * *"
	if _, err := c.AddFunc(spec, UpdateNodeData); err != nil {
		debug.PrintStack()
		return err
	}

	// 首次更新节点数据（注册到Redis）
	UpdateNodeData()

	// 消息订阅
	var sub database.Subscriber
	sub.Connect()

	// 获取当前节点
	node, err := GetCurrentNode()
	if err != nil {
		log.Errorf(err.Error())
		return err
	}

	if IsMaster() {
		// 如果为主节点，订阅主节点通信频道
		channel := "nodes:master"
		sub.Subscribe(channel, MasterNodeCallback)
	} else {
		// 若为工作节点，订阅单独指定通信频道
		channel := "nodes:" + node.Id.Hex()
		sub.Subscribe(channel, WorkerNodeCallback)
	}

	// 如果为主节点，每30秒刷新所有节点信息
	if IsMaster() {
		spec := "*/10 * * * * *"
		if _, err := c.AddFunc(spec, UpdateNodeStatus); err != nil {
			debug.PrintStack()
			return err
		}
	}

	c.Start()
	return nil
}
