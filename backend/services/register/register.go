package register

import (
	"bytes"
	"crawlab/constants"
	"fmt"
	"github.com/apex/log"
	"github.com/spf13/viper"
	"net"
	"os/exec"
	"reflect"
	"runtime/debug"
	"strings"
	"sync"
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
	// 注册节点的Hostname
	GetHostname() (string, error)
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

func (mac *MacRegister) GetHostname() (string, error) {
	return getHostname()
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

func (ip *IpRegister) GetHostname() (string, error) {
	return getHostname()
}

// ===================== mac 地址注册 =====================
type HostnameRegister struct{}

func (h *HostnameRegister) GetType() string {
	return "mac"
}

func (h *HostnameRegister) GetKey() (string, error) {
	return h.GetHostname()
}

func (h *HostnameRegister) GetMac() (string, error) {
	return getMac()
}

func (h *HostnameRegister) GetIp() (string, error) {
	return getIp()
}

func (h *HostnameRegister) GetHostname() (string, error) {
	return getHostname()
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

func getHostname() (string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("hostname")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		log.Errorf(err.Error())
		log.Errorf(fmt.Sprintf("error: %s", stderr.String()))
		debug.PrintStack()
		return "", err
	}

	return strings.Replace(stdout.String(), "\n", "", -1), nil
}

// ===================== 获得注册简单工厂 =====================
var register Register

// 获得注册器
var once sync.Once

func GetRegister() Register {
	once.Do(func() {
		registerType := viper.GetString("server.register.type")
		if registerType == constants.RegisterTypeMac {
			register = &MacRegister{}
		} else if registerType == constants.RegisterTypeIp {
			ip := viper.GetString("server.register.ip")
			if ip == "" {
				log.Error("server.register.ip is empty")
				debug.PrintStack()
				register = nil
			}
			register = &IpRegister{
				Ip: ip,
			}
		} else if registerType == constants.RegisterTypeHostname {
			register = &HostnameRegister{}
		}
		log.Info("register type is :" + reflect.TypeOf(register).String())

	})
	return register
}
