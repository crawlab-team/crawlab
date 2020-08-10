package services

import (
	"crawlab/constants"
	"crawlab/database"
	"crawlab/entity"
	"crawlab/lib/cron"
	"crawlab/model"
	"crawlab/services/rpc"
	"crawlab/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/apex/log"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/imroc/req"
	"os/exec"
	"regexp"
	"runtime/debug"
	"sort"
	"strings"
	"sync"
)

// 系统信息 chan 映射
var SystemInfoChanMap = utils.NewChanMap()

// 从远端获取系统信息
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

// 获取系统信息
func GetSystemInfo(nodeId string) (sysInfo entity.SystemInfo, err error) {
	if IsMasterNode(nodeId) {
		sysInfo, err = rpc.GetSystemInfoServiceLocal()
	} else {
		sysInfo, err = rpc.GetSystemInfoServiceRemote(nodeId)
	}
	return
}

// 获取语言列表
func GetLangList(nodeId string) []entity.Lang {
	list := utils.GetLangList()
	for i, lang := range list {
		status, _ := GetLangInstallStatus(nodeId, lang)
		list[i].InstallStatus = status
	}
	return list
}

// 获取语言安装状态
func GetLangInstallStatus(nodeId string, lang entity.Lang) (string, error) {
	_, err := model.GetTaskByFilter(bson.M{
		"node_id": nodeId,
		"cmd":     fmt.Sprintf("sh %s", utils.GetSystemScriptPath(lang.InstallScript)),
		"status": bson.M{
			"$in": []string{constants.StatusPending, constants.StatusRunning},
		},
	})
	if err == nil {
		// 任务正在运行，正在安装
		return constants.InstallStatusInstalling, nil
	}
	if err != mgo.ErrNotFound {
		// 发生错误
		return "", err
	}
	// 获取状态
	if IsMasterNode(nodeId) {
		lang := rpc.GetLangLocal(lang)
		return lang.InstallStatus, nil
	} else {
		lang, err := rpc.GetLangRemote(nodeId, lang)
		if err != nil {
			return "", err
		}
		return lang.InstallStatus, nil
	}
}

// 是否已安装该依赖
func IsInstalledDep(installedDepList []entity.Dependency, dep entity.Dependency) bool {
	for _, _dep := range installedDepList {
		if strings.ToLower(_dep.Name) == strings.ToLower(dep.Name) {
			return true
		}
	}
	return false
}

// ========Python========

// 初始化函数
func InitDepsFetcher() error {
	c := cron.New(cron.WithSeconds())
	c.Start()
	if _, err := c.AddFunc("0 */5 * * * *", UpdatePythonDepList); err != nil {
		return err
	}

	go func() {
		UpdatePythonDepList()
	}()
	return nil
}

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

// 获取Python本地依赖列表
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

	// 获取已安装依赖列表
	var installedDepList []entity.Dependency
	if IsMasterNode(nodeId) {
		installedDepList, err = rpc.GetInstalledDepsLocal(constants.Python)
		if err != nil {
			return list, err
		}
	} else {
		installedDepList, err = rpc.GetInstalledDepsRemote(nodeId, constants.Python)
		if err != nil {
			return list, err
		}
	}

	// 根据依赖名排序
	sort.Stable(depNameList)

	// 遍历依赖名列表，取前20个
	for i, depNameDict := range depNameList {
		if i > 20 {
			break
		}
		dep := entity.Dependency{
			Name: depNameDict.Name,
		}
		dep.Installed = IsInstalledDep(installedDepList, dep)
		list = append(list, dep)
	}

	// 从依赖源获取信息
	//list, err = GetPythonDepListWithInfo(list)

	return list, nil
}

// 获取Python依赖的源数据信息
func GetPythonDepListWithInfo(depList []entity.Dependency) ([]entity.Dependency, error) {
	var goSync sync.WaitGroup
	for i, dep := range depList {
		if i > 10 {
			break
		}
		goSync.Add(1)
		go func(i int, dep entity.Dependency, depList []entity.Dependency, n *sync.WaitGroup) {
			url := fmt.Sprintf("https://pypi.org/pypi/%s/json", dep.Name)
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
			depList[i].Version = data.Info.Version
			depList[i].Description = data.Info.Summary
			n.Done()
		}(i, dep, depList, &goSync)
	}
	goSync.Wait()
	return depList, nil
}

func FetchPythonDepInfo(depName string) (entity.Dependency, error) {
	url := fmt.Sprintf("https://pypi.org/pypi/%s/json", depName)
	res, err := req.Get(url)
	if err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return entity.Dependency{}, err
	}
	var data PythonDepJsonData
	if res.Response().StatusCode == 404 {
		return entity.Dependency{}, errors.New("get depName from [https://pypi.org] error: 404")
	}
	if err := res.ToJSON(&data); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return entity.Dependency{}, err
	}
	dep := entity.Dependency{
		Name:        depName,
		Version:     data.Info.Version,
		Description: data.Info.Summary,
	}
	return dep, nil
}

// 从Redis获取Python依赖列表
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

// 从Python依赖源获取依赖列表并返回
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

// 更新Python依赖列表到Redis
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

// ========./Python========

// ========Node.js========

// 获取Nodejs本地依赖列表
func GetNodejsDepList(nodeId string, searchDepName string) (depList []entity.Dependency, err error) {
	// 执行shell命令
	cmd := exec.Command("npm", "search", "--json", searchDepName)
	outputBytes, _ := cmd.Output()

	// 获取已安装依赖列表
	var installedDepList []entity.Dependency
	if IsMasterNode(nodeId) {
		installedDepList, err = rpc.GetInstalledDepsLocal(constants.Nodejs)
		if err != nil {
			return depList, err
		}
	} else {
		installedDepList, err = rpc.GetInstalledDepsRemote(nodeId, constants.Nodejs)
		if err != nil {
			return depList, err
		}
	}

	// 反序列化
	if err := json.Unmarshal(outputBytes, &depList); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return depList, err
	}

	// 遍历安装列表
	for i, dep := range depList {
		depList[i].Installed = IsInstalledDep(installedDepList, dep)
	}

	return depList, nil
}

// ========./Node.js========
