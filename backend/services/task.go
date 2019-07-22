package services

import (
	"crawlab/constants"
	"crawlab/database"
	"crawlab/lib/cron"
	"crawlab/model"
	"crawlab/utils"
	"encoding/json"
	"errors"
	"github.com/apex/log"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strconv"
	"time"
)

var Exec *Executor

// 任务执行锁
var LockList []bool

// 任务消息
type TaskMessage struct {
	Id  string
	Cmd string
}

// 序列化任务消息
func (m *TaskMessage) ToString() (string, error) {
	data, err := json.Marshal(&m)
	if err != nil {
		return "", err
	}
	return string(data), err
}

// 任务执行器
type Executor struct {
	Cron *cron.Cron
}

// 启动任务执行器
func (ex *Executor) Start() error {
	// 启动cron服务
	ex.Cron.Start()

	// 加入执行器到定时任务
	spec := "0/1 * * * * *" // 每秒执行一次
	for i := 0; i < viper.GetInt("task.workers"); i++ {
		// WorkerID
		id := i

		// 初始化任务锁
		LockList = append(LockList, false)

		// 加入定时任务
		_, err := ex.Cron.AddFunc(spec, GetExecuteTaskFunc(id))
		if err != nil {
			return err
		}
	}

	return nil
}

var TaskExecChanMap = utils.NewChanMap()

// 派发任务
func AssignTask(task model.Task) error {
	// 生成任务信息
	msg := TaskMessage{
		Id: task.Id,
	}

	// 序列化
	msgStr, err := msg.ToString()
	if err != nil {
		return err
	}

	// 队列名称
	var queue string
	if utils.IsObjectIdNull(task.NodeId) {
		queue = "tasks:public"
	} else {
		queue = "tasks:node:" + task.NodeId.Hex()
	}

	// 任务入队
	if err := database.RedisClient.RPush(queue, msgStr); err != nil {
		return err
	}
	return nil
}

// 执行shell命令
func ExecuteShellCmd(cmdStr string, cwd string, t model.Task, s model.Spider) (err error) {
	log.Infof("cwd: " + cwd)
	log.Infof("cmd: " + cmdStr)

	// 生成执行命令
	var cmd *exec.Cmd
	if runtime.GOOS == constants.Windows {
		cmd = exec.Command("cmd", "/C", cmdStr)
	} else {
		cmd = exec.Command("sh", "-c", cmdStr)
	}

	// 工作目录
	cmd.Dir = cwd

	// 指定stdout, stderr日志位置
	fLog, err := os.Create(t.LogPath)
	if err != nil {
		HandleTaskError(t, err)
		return err
	}
	defer fLog.Close()
	cmd.Stdout = fLog
	cmd.Stderr = fLog

	// 添加环境变量
	cmd.Env = append(cmd.Env, "CRAWLAB_TASK_ID="+t.Id)
	cmd.Env = append(cmd.Env, "CRAWLAB_COLLECTION="+s.Col)

	// 起一个goroutine来监控进程
	ch := TaskExecChanMap.ChanBlocked(t.Id)
	go func() {
		// 传入信号，此处阻塞
		signal := <-ch

		if signal == constants.TaskCancel {
			// 取消进程
			if err := cmd.Process.Kill(); err != nil {
				log.Errorf(err.Error())
				debug.PrintStack()
				return
			}
			t.Status = constants.StatusCancelled
		} else if signal == constants.TaskFinish {
			// 完成进程
			t.Status = constants.StatusFinished
		}

		// 保存任务
		t.FinishTs = time.Now()
		if err := t.Save(); err != nil {
			log.Infof(err.Error())
			debug.PrintStack()
			return
		}
	}()

	// 开始执行
	if err := cmd.Run(); err != nil {
		HandleTaskError(t, err)
		return err
	}
	ch <- constants.TaskFinish

	return nil
}

// 生成日志目录
func MakeLogDir(t model.Task) (fileDir string, err error) {
	// 日志目录
	fileDir = filepath.Join(viper.GetString("log.path"), t.SpiderId.Hex())

	// 如果日志目录不存在，生成该目录
	if !utils.Exists(fileDir) {
		if err := os.MkdirAll(fileDir, 0777); err != nil {
			debug.PrintStack()
			return "", err
		}
	}

	return fileDir, nil
}

// 获取日志文件路径
func GetLogFilePaths(fileDir string) (filePath string) {
	// 时间戳
	ts := time.Now()
	tsStr := ts.Format("20060102150405")

	// stdout日志文件
	filePath = filepath.Join(fileDir, tsStr+".log")

	return filePath
}

// 生成执行任务方法
func GetExecuteTaskFunc(id int) func() {
	return func() {
		ExecuteTask(id)
	}
}

func GetWorkerPrefix(id int) string {
	return "[Worker " + strconv.Itoa(id) + "] "
}

// 执行任务
func ExecuteTask(id int) {
	if LockList[id] {
		log.Debugf(GetWorkerPrefix(id) + "正在执行任务...")
		return
	}

	// 上锁
	LockList[id] = true

	// 解锁（延迟执行）
	defer func() {
		LockList[id] = false
	}()

	// 开始计时
	tic := time.Now()

	// 获取当前节点
	node, err := GetCurrentNode()
	if err != nil {
		log.Errorf(GetWorkerPrefix(id) + err.Error())
		return
	}

	// 公共队列
	queuePub := "tasks:public"

	// 节点队列
	queueCur := "tasks:node:" + node.Id.Hex()

	// 节点队列任务
	var msg string
	msg, err = database.RedisClient.LPop(queueCur)
	if err != nil {
		if msg == "" {
			// 节点队列没有任务，获取公共队列任务
			msg, err = database.RedisClient.LPop(queuePub)
			if err != nil {
				if msg == "" {
					// 公共队列没有任务
					log.Debugf(GetWorkerPrefix(id) + "没有任务...")
					return
				} else {
					log.Errorf(GetWorkerPrefix(id) + err.Error())
					debug.PrintStack()
					return
				}
			}
		} else {
			log.Errorf(GetWorkerPrefix(id) + err.Error())
			debug.PrintStack()
			return
		}
	}

	// 反序列化
	tMsg := TaskMessage{}
	if err := json.Unmarshal([]byte(msg), &tMsg); err != nil {
		log.Errorf(GetWorkerPrefix(id) + err.Error())
		debug.PrintStack()
		return
	}

	// 获取任务
	t, err := model.GetTask(tMsg.Id)
	if err != nil {
		log.Errorf(GetWorkerPrefix(id) + err.Error())
		return
	}

	// 获取爬虫
	spider, err := t.GetSpider()
	if err != nil {
		log.Errorf(GetWorkerPrefix(id) + err.Error())
		return
	}

	// 创建日志目录
	fileDir, err := MakeLogDir(t)
	if err != nil {
		log.Errorf(GetWorkerPrefix(id) + err.Error())
		return
	}

	// 获取日志文件路径
	t.LogPath = GetLogFilePaths(fileDir)

	// 创建日志目录文件夹
	fileStdoutDir := filepath.Dir(t.LogPath)
	if !utils.Exists(fileStdoutDir) {
		if err := os.MkdirAll(fileStdoutDir, os.ModePerm); err != nil {
			log.Errorf(GetWorkerPrefix(id) + err.Error())
			return
		}
	}

	// 工作目录
	cwd := filepath.Join(
		viper.GetString("spider.path"),
		spider.Name,
	)

	// 执行命令
	cmd := spider.Cmd
	if t.Cmd != "" {
		cmd = t.Cmd
	}

	// 任务赋值
	t.NodeId = node.Id                 // 任务节点信息
	t.StartTs = time.Now()             // 任务开始时间
	t.Status = constants.StatusRunning // 任务状态

	// 开始执行任务
	log.Infof(GetWorkerPrefix(id) + "开始执行任务(ID:" + t.Id + ")")

	// 储存任务
	if err := t.Save(); err != nil {
		log.Errorf(err.Error())
		HandleTaskError(t, err)
		return
	}

	// 执行Shell命令
	if err := ExecuteShellCmd(cmd, cwd, t, spider); err != nil {
		log.Errorf(GetWorkerPrefix(id) + err.Error())
		return
	}

	// 结束计时
	toc := time.Now()

	// 统计时长
	duration := toc.Sub(tic).Seconds()
	durationStr := strconv.FormatFloat(duration, 'f', 6, 64)
	log.Infof(GetWorkerPrefix(id) + "任务(ID:" + t.Id + ")" + "执行完毕. 消耗时间:" + durationStr + "秒")
}

func GetTaskLog(id string) (logStr string, err error) {
	task, err := model.GetTask(id)
	if err != nil {
		return "", err
	}

	logStr = ""
	if IsMaster() {
		// 若为主节点，获取本机日志
		logBytes, err := GetLocalLog(task.LogPath)
		logStr = string(logBytes)
		if err != nil {
			log.Errorf(err.Error())
			return "", err
		}
		logStr = string(logBytes)
	} else {
		// 若不为主节点，获取远端日志
		logStr, err = GetRemoteLog(task)
		if err != nil {
			log.Errorf(err.Error())
			return "", err
		}
	}

	return logStr, nil
}

func CancelTask(id string) (err error) {
	// 获取任务
	task, err := model.GetTask(id)
	if err != nil {
		return err
	}

	// 如果任务状态不为pending或running，返回错误
	if task.Status != constants.StatusPending && task.Status != constants.StatusRunning {
		return errors.New("task is not cancellable")
	}

	// 获取当前节点（默认当前节点为主节点）
	node, err := GetCurrentNode()
	if err != nil {
		return err
	}

	if node.Id == task.NodeId {
		// 任务节点为主节点

		// 获取任务执行频道
		ch := TaskExecChanMap.ChanBlocked(id)

		// 发出取消进程信号
		ch <- constants.TaskCancel
	} else {
		// 任务节点为工作节点

		// 序列化消息
		msg := NodeMessage{
			Type:   constants.MsgTypeCancelTask,
			TaskId: id,
		}
		msgBytes, err := json.Marshal(&msg)
		if err != nil {
			return err
		}

		// 发布消息
		if err := database.Publish("nodes:"+task.NodeId.Hex(), string(msgBytes)); err != nil {
			return err
		}
	}

	return nil
}

func HandleTaskError(t model.Task, err error) {
	t.Status = constants.StatusError
	t.Error = err.Error()
	t.FinishTs = time.Now()
	if err := t.Save(); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return
	}
	debug.PrintStack()
}

func InitTaskExecutor() error {
	c := cron.New(cron.WithSeconds())
	Exec = &Executor{
		Cron: c,
	}
	if err := Exec.Start(); err != nil {
		return err
	}
	return nil
}
