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
	GetType() string
	// 注册的key的值，唯一标识节点
	GetKey() (string, error)
	// 注册的节点IP
	GetIp() (string, error)
	// 注册节点的mac地址
	GetMac() (string, error)
}

// ===================== mac 地址注册 =====================
type MacRegister struct{}

func (mac *MacRegister) GetType() string {
	return "mac"
}

func (mac *MacRegister) GetKey() (string, error) {
	return mac.GetMac()
}

func (mac *MacRegister) GetMac() (string, error) {
	return getMac()
}

func (mac *MacRegister) GetIp() (string, error) {
	return getIp()
}

// ===================== ip 地址注册 =====================
type IpRegister struct {
	Ip string
}

func (ip *IpRegister) GetType() string {
	return "ip"
}

func (ip *IpRegister) GetKey() (string, error) {
	return ip.Ip, nil
}

func (ip *IpRegister) GetIp() (string, error) {
	return ip.Ip, nil
}

func (ip *IpRegister) GetMac() (string, error) {
	return getMac()
}

// ===================== 公共方法 =====================
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

func getMac() (string, error) {
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

// ===================== 获得注册简单工厂 =====================
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
