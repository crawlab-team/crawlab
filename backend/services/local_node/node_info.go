package local_node

import (
	"crawlab/services/local_machine_info"
	"errors"
	"github.com/hashicorp/go-sockaddr"
	"os"
)

var localNode *LocalNode

type IdentifyType string

const (
	Ip       = IdentifyType("ip")
	Mac      = IdentifyType("mac")
	Hostname = IdentifyType("hostname")
)

type local struct {
	Ip           string
	Mac          string
	Hostname     string
	Identify     string
	IdentifyType IdentifyType
}

type LocalNode struct {
	local
	mongo
	local_machine_info.MachineInfo
}

func (l *LocalNode) Ready() error {
	err := localNode.load(true)
	if err != nil {
		return err
	}
	go localNode.watch()
	return nil
}

func NewLocalNode(ip string, identify string, identifyTypeString string) (node *LocalNode, err error) {
	addrs, err := sockaddr.GetPrivateInterfaces()
	if ip == "" {
		if err != nil {
			return node, err
		}
		if len(addrs) == 0 {
			return node, errors.New("address not found")
		}
		ipaddr := *sockaddr.ToIPAddr(addrs[0].SockAddr)
		ip = ipaddr.NetIP().String()
	}

	mac := addrs[0].HardwareAddr.String()
	hostname, err := os.Hostname()
	if err != nil {
		return node, err
	}
	local := local{Ip: ip, Mac: mac, Hostname: hostname}

	// 获取机器使用信息
	machineinfo := local_machine_info.CollectMachineInfo()

	switch IdentifyType(identifyTypeString) {
	case Ip:
		local.Identify = local.Ip
		local.IdentifyType = Ip
	case Mac:
		local.Identify = local.Mac
		local.IdentifyType = Mac
	case Hostname:
		local.Identify = local.Hostname
		local.IdentifyType = Hostname
	default:
		local.Identify = identify
		local.IdentifyType = IdentifyType(identifyTypeString)
	}
	return &LocalNode{local: local, mongo: mongo{}, MachineInfo: machineinfo}, nil
}
