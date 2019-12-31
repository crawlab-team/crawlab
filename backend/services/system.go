package services

import (
	"crawlab/constants"
	"crawlab/database"
	"crawlab/entity"
	"crawlab/lib/cron"
	"crawlab/model"
	"crawlab/utils"
	"encoding/json"
	"fmt"
	"github.com/apex/log"
	"github.com/imroc/req"
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
		list[i].Installed = isInstalledLang(nodeId, lang)
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

func GetDepList(nodeId string, langExecutableName string, searchDepName string) ([]entity.Dependency, error) {
	// TODO: support other languages
	// 获取语言
	lang := GetLangFromLangName(nodeId, langExecutableName)

	// 如果没有依赖列表，先获取
	if len(DepList) == 0 {
		FetchDepList()
	}

	// 过滤相似的依赖
	var depNameList PythonDepNameDictSlice
	for _, depName := range DepList {
		if strings.Contains(strings.ToLower(depName), strings.ToLower(searchDepName)) {
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

	var depList []entity.Dependency
	var goSync sync.WaitGroup
	sort.Stable(depNameList)
	for i, depNameDict := range depNameList {
		if i > 20 {
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
				Lang:        lang.ExecutableName,
				Installed:   false,
			}
			depList = append(depList, dep)
			n.Done()
		}(depNameDict.Name, &goSync)
	}
	goSync.Wait()

	return depList, nil
}

func isInstalledLang(nodeId string, lang entity.Lang) bool {
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

func FetchDepList() {
	url := "https://pypi.tuna.tsinghua.edu.cn/simple"
	res, err := req.Get(url)
	if err != nil {
		log.Error(err.Error())
		debug.PrintStack()
		return
	}
	text, err := res.ToString()
	if err != nil {
		log.Error(err.Error())
		debug.PrintStack()
		return
	}
	var list []string
	regex := regexp.MustCompile("<a href=\".*/\">(.*)</a>")
	for _, line := range strings.Split(text, "\n") {
		arr := regex.FindStringSubmatch(line)
		if len(arr) < 2 {
			continue
		}
		list = append(list, arr[1])
	}
	DepList = list
}

var DepList []string

func InitDepsFetcher() error {
	c := cron.New(cron.WithSeconds())
	c.Start()
	if _, err := c.AddFunc("0 */5 * * * *", FetchDepList); err != nil {
		return err
	}
	return nil
}
