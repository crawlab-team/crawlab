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
	"github.com/apex/log"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strconv"
	"sync"
	"syscall"
	"time"
)

var Exec *Executor

// 任务执行锁
//Added by cloud: 2019/09/04,solve data race
var LockList sync.Map

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
	return utils.BytesToString(data), err
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
		LockList.Store(id, false)

		// 加入定时任务
		_, err := ex.Cron.AddFunc(spec, GetExecuteTaskFunc(id))
		if err != nil {
			return err
		}
	}

	return nil
}

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

// 设置环境变量
func SetEnv(cmd *exec.Cmd, envs []model.Env, taskId string, dataCol string) *exec.Cmd {
	// 默认环境变量
	cmd.Env = append(os.Environ(), "CRAWLAB_TASK_ID="+taskId)
	cmd.Env = append(cmd.Env, "CRAWLAB_COLLECTION="+dataCol)
	cmd.Env = append(cmd.Env, "PYTHONUNBUFFERED=0")
	cmd.Env = append(cmd.Env, "PYTHONIOENCODING=utf-8")
	cmd.Env = append(cmd.Env, "TZ=Asia/Shanghai")

	//任务环境变量
	for _, env := range envs {
		cmd.Env = append(cmd.Env, env.Name+"="+env.Value)
	}

	// TODO 全局环境变量
	return cmd
}

func SetLogConfig(cmd *exec.Cmd, path string) error {
	fLog, err := os.Create(path)
	if err != nil {
		log.Errorf("create task log file error: %s", path)
		debug.PrintStack()
		return err
	}
	cmd.Stdout = fLog
	cmd.Stderr = fLog
	return nil
}

func FinishOrCancelTask(ch chan string, cmd *exec.Cmd, t model.Task) {
	// 传入信号，此处阻塞
	signal := <-ch
	log.Infof("process received signal: %s", signal)

	if signal == constants.TaskCancel && cmd.Process != nil {
		// 取消进程
		if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); err != nil {
			log.Errorf("process kill error: %s", err.Error())
			debug.PrintStack()

			t.Error = "kill process error: " + err.Error()
			t.Status = constants.StatusError
		} else {
			t.Error = "user kill the process ..."
			t.Status = constants.StatusCancelled
		}
	} else {
		// 保存任务
		t.Status = constants.StatusFinished
	}

	t.FinishTs = time.Now()
	_ = t.Save()
}

func StartTaskProcess(cmd *exec.Cmd, t model.Task) error {
	if err := cmd.Start(); err != nil {
		log.Errorf("start spider error:{}", err.Error())
		debug.PrintStack()

		t.Error = "start task error: " + err.Error()
		t.Status = constants.StatusError
		t.FinishTs = time.Now()
		_ = t.Save()
		return err
	}
	return nil
}

func WaitTaskProcess(cmd *exec.Cmd, t model.Task) error {
	if err := cmd.Wait(); err != nil {
		log.Errorf("wait process finish error: %s", err.Error())
		debug.PrintStack()

		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode := exitError.ExitCode()
			log.Errorf("exit error, exit code: %d", exitCode)

			// 非kill 的错误类型
			if exitCode != -1 {
				// 非手动kill保存为错误状态
				t.Error = err.Error()
				t.FinishTs = time.Now()
				t.Status = constants.StatusError
				_ = t.Save()
			}
		}

		return err
	}
	return nil
}

// 执行shell命令
func ExecuteShellCmd(cmdStr string, cwd string, t model.Task, s model.Spider) (err error) {
	log.Infof("cwd: %s", cwd)
	log.Infof("cmd: %s", cmdStr)

	// 生成执行命令
	var cmd *exec.Cmd
	if runtime.GOOS == constants.Windows {
		cmd = exec.Command("cmd", "/C", cmdStr)
	} else {
		cmd = exec.Command("")
		cmd = exec.Command("sh", "-c", cmdStr)
	}

	// 工作目录
	cmd.Dir = cwd

	// 日志配置
	if err := SetLogConfig(cmd, t.LogPath); err != nil {
		return err
	}

	// 环境变量配置
	envs := s.Envs
	if s.Type == constants.Configurable {
		envs = append(envs, model.Env{Name: "CRAWLAB_MONGO_HOST", Value: viper.GetString("mongo.host")})
		envs = append(envs, model.Env{Name: "CRAWLAB_MONGO_PORT", Value: viper.GetString("mongo.port")})
		envs = append(envs, model.Env{Name: "CRAWLAB_MONGO_DB", Value: viper.GetString("mongo.db")})
		envs = append(envs, model.Env{Name: "CRAWLAB_MONGO_USERNAME", Value: viper.GetString("mongo.username")})
		envs = append(envs, model.Env{Name: "CRAWLAB_MONGO_PASSWORD", Value: viper.GetString("mongo.password")})
		envs = append(envs, model.Env{Name: "CRAWLAB_MONGO_AUTHSOURCE", Value: viper.GetString("mongo.authSource")})
	}
	cmd = SetEnv(cmd, envs, t.Id, s.Col)

	// 起一个goroutine来监控进程
	ch := utils.TaskExecChanMap.ChanBlocked(t.Id)

	go FinishOrCancelTask(ch, cmd, t)

	// kill的时候，可以kill所有的子进程
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	// 启动进程
	if err := StartTaskProcess(cmd, t); err != nil {
		return err
	}

	// 同步等待进程完成
	if err := WaitTaskProcess(cmd, t); err != nil {
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
			log.Errorf("execute task, make log dir error: %s", err.Error())
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

// 统计任务结果数
func SaveTaskResultCount(id string) func() {
	return func() {
		if err := model.UpdateTaskResultCount(id); err != nil {
			log.Errorf(err.Error())
			debug.PrintStack()
			return
		}
	}
}

// 执行任务
func ExecuteTask(id int) {
	if flag, ok := LockList.Load(id); ok {
		if flag.(bool) {
			log.Debugf(GetWorkerPrefix(id) + "正在执行任务...")
			return
		}

	}

	// 上锁
	LockList.Store(id, true)

	// 解锁（延迟执行）
	defer func() {
		LockList.Delete(id)
		LockList.Store(id, false)
	}()

	// 开始计时
	tic := time.Now()

	// 获取当前节点
	node, err := model.GetCurrentNode()
	if err != nil {
		log.Errorf("execute task get current node error: %s", err.Error())
		debug.PrintStack()
		return
	}

	// 节点队列
	queueCur := "tasks:node:" + node.Id.Hex()
	// 节点队列任务
	var msg string
	if msg, err = database.RedisClient.LPop(queueCur); err != nil {
		// 节点队列没有任务，获取公共队列任务
		queuePub := "tasks:public"
		if msg, err = database.RedisClient.LPop(queuePub); err != nil {
		}
	}

	if msg == "" {
		return
	}

	// 反序列化
	tMsg := TaskMessage{}
	if err := json.Unmarshal([]byte(msg), &tMsg); err != nil {
		log.Errorf("json string to struct error: %s", err.Error())
		return
	}

	// 获取任务
	t, err := model.GetTask(tMsg.Id)
	if err != nil {
		log.Errorf("execute task, get task error: %s", err.Error())
		return
	}

	// 获取爬虫
	spider, err := t.GetSpider()
	if err != nil {
		log.Errorf("execute task, get spider error: %s", err.Error())
		return
	}

	// 创建日志目录
	var fileDir string
	if fileDir, err = MakeLogDir(t); err != nil {
		return
	}
	// 获取日志文件路径
	t.LogPath = GetLogFilePaths(fileDir)

	// 工作目录
	cwd := filepath.Join(
		viper.GetString("spider.path"),
		spider.Name,
	)

	// 执行命令
	var cmd string
	if spider.Type == constants.Configurable {
		// 可配置爬虫命令
		cmd = "scrapy crawl config_spider"
	} else {
		// 自定义爬虫命令
		cmd = spider.Cmd
	}

	// 加入参数
	if t.Param != "" {
		cmd += " " + t.Param
	}

	// 任务赋值
	t.NodeId = node.Id                                   // 任务节点信息
	t.StartTs = time.Now()                               // 任务开始时间
	t.Status = constants.StatusRunning                   // 任务状态
	t.WaitDuration = t.StartTs.Sub(t.CreateTs).Seconds() // 等待时长

	// 开始执行任务
	log.Infof(GetWorkerPrefix(id) + "开始执行任务(ID:" + t.Id + ")")

	// 储存任务
	if err := t.Save(); err != nil {
		log.Errorf(err.Error())
		HandleTaskError(t, err)
		return
	}

	// 起一个cron执行器来统计任务结果数
	if spider.Col != "" {
		cronExec := cron.New(cron.WithSeconds())
		_, err = cronExec.AddFunc("*/5 * * * * *", SaveTaskResultCount(t.Id))
		if err != nil {
			log.Errorf(GetWorkerPrefix(id) + err.Error())
			return
		}
		cronExec.Start()
		defer cronExec.Stop()
	}

	// 执行Shell命令
	if err := ExecuteShellCmd(cmd, cwd, t, spider); err != nil {
		log.Errorf(GetWorkerPrefix(id) + err.Error())
		return
	}

	// 更新任务结果数
	if spider.Col != "" {
		if err := model.UpdateTaskResultCount(t.Id); err != nil {
			log.Errorf(GetWorkerPrefix(id) + err.Error())
			return
		}
	}

	// 完成进程
	t, err = model.GetTask(t.Id)
	if err != nil {
		log.Errorf(GetWorkerPrefix(id) + err.Error())
		return
	}
	t.Status = constants.StatusFinished                     // 任务状态: 已完成
	t.FinishTs = time.Now()                                 // 结束时间
	t.RuntimeDuration = t.FinishTs.Sub(t.StartTs).Seconds() // 运行时长
	t.TotalDuration = t.FinishTs.Sub(t.CreateTs).Seconds()  // 总时长

	// 保存任务
	if err := t.Save(); err != nil {
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
		return
	}

	if IsMasterNode(task.NodeId.Hex()) {
		if !utils.Exists(task.LogPath) {
			fileDir, err := MakeLogDir(task)

			if err != nil {
				log.Errorf(err.Error())
			}

			fileP := GetLogFilePaths(fileDir)

			// 获取日志文件路径
			fLog, err := os.Create(fileP)
			defer fLog.Close()
			if err != nil {
				log.Errorf("create task log file error: %s", fileP)
				debug.PrintStack()
			}
			task.LogPath = fileP
			if err := task.Save(); err != nil {
				log.Errorf(err.Error())
				debug.PrintStack()
			}

		}
		// 若为主节点，获取本机日志
		logBytes, err := model.GetLocalLog(task.LogPath)
		if err != nil {
			log.Errorf(err.Error())
			logStr = err.Error()
		} else {
			logStr = utils.BytesToString(logBytes)
		}
		return logStr, err
	}
	// 若不为主节点，获取远端日志
	logStr, err = GetRemoteLog(task)
	if err != nil {
		log.Errorf(err.Error())

	}
	return logStr, err
}

func CancelTask(id string) (err error) {
	// 获取任务
	task, err := model.GetTask(id)
	if err != nil {
		log.Errorf("task not found, task id : %s, error: %s", id, err.Error())
		debug.PrintStack()
		return err
	}

	// 如果任务状态不为pending或running，返回错误
	if task.Status != constants.StatusPending && task.Status != constants.StatusRunning {
		return errors.New("task is not cancellable")
	}

	// 获取当前节点（默认当前节点为主节点）
	node, err := model.GetCurrentNode()
	if err != nil {
		log.Errorf("get current node error: %s", err.Error())
		debug.PrintStack()
		return err
	}

	log.Infof("current node id is: %s", node.Id.Hex())
	log.Infof("task node id is: %s", task.NodeId.Hex())

	if node.Id == task.NodeId {
		// 任务节点为主节点

		// 获取任务执行频道
		ch := utils.TaskExecChanMap.ChanBlocked(id)
		if ch != nil {
			// 发出取消进程信号
			ch <- constants.TaskCancel
		} else {
			if err := model.UpdateTaskToAbnormal(node.Id); err != nil {
				log.Errorf("update task to abnormal : {}", err.Error())
				debug.PrintStack()
				return err
			}
		}
	} else {
		// 任务节点为工作节点

		// 序列化消息
		msg := entity.NodeMessage{
			Type:   constants.MsgTypeCancelTask,
			TaskId: id,
		}
		msgBytes, err := json.Marshal(&msg)
		if err != nil {
			return err
		}

		// 发布消息
		if _, err := database.RedisClient.Publish("nodes:"+task.NodeId.Hex(), utils.BytesToString(msgBytes)); err != nil {
			return err
		}
	}

	return nil
}

func HandleTaskError(t model.Task, err error) {
	log.Error("handle task error:" + err.Error())
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
