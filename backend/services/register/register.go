package register

import (
	"github.com/apex/log"
	"github.com/spf13/viper"
	"net"
	"reflect"
	"runtime/debug"
)

type Register interface {
	// 注册的key类型
	GetKey() string
	// 注册的key的值，唯一标识节点
	GetValue() (string, error)
	// 注册的节点IP
	GetIp() (string, error)
}

// mac 地址注册
type MacRegister struct{}

func (mac *MacRegister) GetKey() string {
	return "mac"
}

func (mac *MacRegister) GetValue() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Errorf("get interfaces error:" + err.Error())
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

func (mac *MacRegister) GetIp() (string, error) {
	return getIp()
}

// ip 注册
type IpRegister struct {
	Ip string
}

func (ip *IpRegister) GetKey() string {
	return "ip"
}

func (ip *IpRegister) GetValue() (string, error) {
	return ip.Ip, nil
}

func (ip *IpRegister) GetIp() (string, error) {
	return ip.Ip, nil
}

// 获取本机的IP地址
// TODO: 考虑多个IP地址的情况
func getIp() (string, error) {
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

var register Register

// 获得注册器
func GetRegister() Register {
	if register != nil {
		return register
	}

	registerType := viper.GetString("server.register.type")
	if registerType == "mac" {
		register = &MacRegister{}
	} else {
		ip := viper.GetString("server.register.ip")
		if ip == "" {
			log.Error("server.register.ip is empty")
			debug.PrintStack()
			return nil
		}
		register = &IpRegister{
			Ip: ip,
		}
	}
	log.Info("register type is :" + reflect.TypeOf(register).String())
	return register
}
