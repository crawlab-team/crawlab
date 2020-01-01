package services

import (
	"crawlab/constants"
	"crawlab/database"
	"crawlab/entity"
	"crawlab/lib/cron"
	"crawlab/model"
	"crawlab/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/apex/log"
	"github.com/imroc/req"
	"os/exec"
	"regexp"
	"runtime/debug"
	"sort"
	"strings"
	"sync"
)

type PythonDepJsonData struct {
	Info PythonDepJsonDataInfo `json:"info"`
}

type PythonDepJsonDataInfo struct {
	Name    string `json:"name"`
	Summary string `json:"summary"`
	Version string `json:"version"`
}

type PythonDepNameDict struct {
	Name   string `json:"name"`
	Weight int    `json:"weight"`
}

type PythonDepNameDictSlice []PythonDepNameDict

func (s PythonDepNameDictSlice) Len() int           { return len(s) }
func (s PythonDepNameDictSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s PythonDepNameDictSlice) Less(i, j int) bool { return s[i].Weight > s[j].Weight }

var SystemInfoChanMap = utils.NewChanMap()

func GetRemoteSystemInfo(nodeId string) (sysInfo entity.SystemInfo, err error) {
	// 发送消息
	msg := entity.NodeMessage{
		Type:   constants.MsgTypeGetSystemInfo,
		NodeId: nodeId,
	}

	// 序列化
	msgBytes, _ := json.Marshal(&msg)
	if _, err := database.RedisClient.Publish("nodes:"+nodeId, utils.BytesToString(msgBytes)); err != nil {
		return entity.SystemInfo{}, err
	}

	// 通道
	ch := SystemInfoChanMap.ChanBlocked(nodeId)

	// 等待响应，阻塞
	sysInfoStr := <-ch

	// 反序列化
	if err := json.Unmarshal([]byte(sysInfoStr), &sysInfo); err != nil {
		return sysInfo, err
	}

	return sysInfo, nil
}

func GetSystemInfo(nodeId string) (sysInfo entity.SystemInfo, err error) {
	if IsMasterNode(nodeId) {
		sysInfo, err = model.GetLocalSystemInfo()
	} else {
		sysInfo, err = GetRemoteSystemInfo(nodeId)
	}
	return
}

func GetLangList(nodeId string) []entity.Lang {
	list := []entity.Lang{
		{Name: "Python", ExecutableName: "python", ExecutablePath: "/usr/local/bin/python", DepExecutablePath: "/usr/local/bin/pip"},
		{Name: "NodeJS", ExecutableName: "node", ExecutablePath: "/usr/local/bin/node"},
		{Name: "Java", ExecutableName: "java", ExecutablePath: "/usr/local/bin/java"},
	}
	for i, lang := range list {
		list[i].Installed = IsInstalledLang(nodeId, lang)
	}
	return list
}

func GetLangFromLangName(nodeId string, name string) entity.Lang {
	langList := GetLangList(nodeId)
	for _, lang := range langList {
		if lang.ExecutableName == name {
			return lang
		}
	}
	return entity.Lang{}
}

func GetPythonDepList(nodeId string, searchDepName string) ([]entity.Dependency, error) {
	var list []entity.Dependency

	// 先从 Redis 获取
	depList, err := GetPythonDepListFromRedis()
	if err != nil {
		return list, err
	}

	// 过滤相似的依赖
	var depNameList PythonDepNameDictSlice
	for _, depName := range depList {
		if strings.HasPrefix(strings.ToLower(depName), strings.ToLower(searchDepName)) {
			var weight int
			if strings.ToLower(depName) == strings.ToLower(searchDepName) {
				weight = 3
			} else if strings.HasPrefix(strings.ToLower(depName), strings.ToLower(searchDepName)) {
				weight = 2
			} else {
				weight = 1
			}
			depNameList = append(depNameList, PythonDepNameDict{
				Name:   depName,
				Weight: weight,
			})
		}
	}

	// 获取已安装依赖
	installedDepList, err := GetPythonInstalledDepList(nodeId)
	if err != nil {
		return list, err
	}

	// 从依赖源获取数据
	var goSync sync.WaitGroup
	sort.Stable(depNameList)
	for i, depNameDict := range depNameList {
		if i > 10 {
			break
		}
		goSync.Add(1)
		go func(depName string, n *sync.WaitGroup) {
			url := fmt.Sprintf("https://pypi.org/pypi/%s/json", depName)
			res, err := req.Get(url)
			if err != nil {
				n.Done()
				return
			}
			var data PythonDepJsonData
			if err := res.ToJSON(&data); err != nil {
				n.Done()
				return
			}
			dep := entity.Dependency{
				Name:        depName,
				Version:     data.Info.Version,
				Description: data.Info.Summary,
			}
			dep.Installed = IsInstalledDep(installedDepList, dep)
			list = append(list, dep)
			n.Done()
		}(depNameDict.Name, &goSync)
	}
	goSync.Wait()

	return list, nil
}

func GetPythonDepListFromRedis() ([]string, error) {
	var list []string

	// 从 Redis 获取字符串
	rawData, err := database.RedisClient.HGet("system", "deps:python")
	if err != nil {
		return list, err
	}

	// 反序列化
	if err := json.Unmarshal([]byte(rawData), &list); err != nil {
		return list, err
	}

	// 如果为空，则从依赖源获取列表
	if len(list) == 0 {
		UpdatePythonDepList()
	}

	return list, nil
}

func IsInstalledLang(nodeId string, lang entity.Lang) bool {
	sysInfo, err := GetSystemInfo(nodeId)
	if err != nil {
		return false
	}
	for _, exec := range sysInfo.Executables {
		if exec.Path == lang.ExecutablePath {
			return true
		}
	}
	return false
}

func IsInstalledDep(installedDepList []entity.Dependency, dep entity.Dependency) bool {
	for _, _dep := range installedDepList {
		if strings.ToLower(_dep.Name) == strings.ToLower(dep.Name) {
			return true
		}
	}
	return false
}

func FetchPythonDepList() ([]string, error) {
	// 依赖URL
	url := "https://pypi.tuna.tsinghua.edu.cn/simple"

	// 输出列表
	var list []string

	// 请求URL
	res, err := req.Get(url)
	if err != nil {
		log.Error(err.Error())
		debug.PrintStack()
		return list, err
	}

	// 获取响应数据
	text, err := res.ToString()
	if err != nil {
		log.Error(err.Error())
		debug.PrintStack()
		return list, err
	}

	// 从响应数据中提取依赖名
	regex := regexp.MustCompile("<a href=\".*/\">(.*)</a>")
	for _, line := range strings.Split(text, "\n") {
		arr := regex.FindStringSubmatch(line)
		if len(arr) < 2 {
			continue
		}
		list = append(list, arr[1])
	}

	// 赋值给列表
	return list, nil
}

func UpdatePythonDepList() {
	// 从依赖源获取列表
	list, _ := FetchPythonDepList()

	// 序列化
	listBytes, err := json.Marshal(list)
	if err != nil {
		log.Error(err.Error())
		debug.PrintStack()
		return
	}

	// 设置Redis
	if err := database.RedisClient.HSet("system", "deps:python", string(listBytes)); err != nil {
		log.Error(err.Error())
		debug.PrintStack()
		return
	}
}

func GetPythonInstalledDepList(nodeId string) ([]entity.Dependency, error){
	var list []entity.Dependency

	lang := GetLangFromLangName(nodeId, constants.Python)
	if !IsInstalledLang(nodeId, lang) {
		return list, errors.New("python is not installed")
	}
	cmd := exec.Command("pip", "freeze")
	outputBytes, err := cmd.Output()
	if err != nil {
		debug.PrintStack()
		return list, err
	}

	for _, line := range strings.Split(string(outputBytes), "\n") {
		arr := strings.Split(line, "==")
		if len(arr) < 2 {
			continue
		}
		dep := entity.Dependency{
			Name:      strings.ToLower(arr[0]),
			Version:   arr[1],
			Installed: true,
		}
		list = append(list, dep)
	}

	return list, nil
}

func InitDepsFetcher() error {
	c := cron.New(cron.WithSeconds())
	c.Start()
	if _, err := c.AddFunc("0 */5 * * * *", UpdatePythonDepList); err != nil {
		return err
	}
	return nil
}
