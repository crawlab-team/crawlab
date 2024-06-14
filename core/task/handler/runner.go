package handler

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/apex/log"
	"github.com/cenkalti/backoff/v4"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/container"
	"github.com/crawlab-team/crawlab/core/entity"
	"github.com/crawlab-team/crawlab/core/errors"
	fs2 "github.com/crawlab-team/crawlab/core/fs"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/client"
	"github.com/crawlab-team/crawlab/core/models/delegate"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/sys_exec"
	"github.com/crawlab-team/crawlab/core/utils"
	"github.com/crawlab-team/crawlab/db/mongo"
	grpc "github.com/crawlab-team/crawlab/grpc"
	"github.com/crawlab-team/crawlab/trace"
	"github.com/shirou/gopsutil/process"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Runner struct {
	// dependencies
	svc     interfaces.TaskHandlerService // task handler service
	fsSvc   interfaces.FsServiceV2        // task fs service
	hookSvc interfaces.TaskHookService    // task hook service

	// settings
	subscribeTimeout time.Duration
	bufferSize       int

	// internals
	cmd  *exec.Cmd                        // process command instance
	pid  int                              // process id
	tid  primitive.ObjectID               // task id
	t    interfaces.Task                  // task model.Task
	s    interfaces.Spider                // spider model.Spider
	ch   chan constants.TaskSignal        // channel to communicate between Service and Runner
	err  error                            // standard process error
	envs []models.Env                     // environment variables
	cwd  string                           // working directory
	c    interfaces.GrpcClient            // grpc client
	sub  grpc.TaskService_SubscribeClient // grpc task service stream client

	// log internals
	scannerStdout *bufio.Reader
	scannerStderr *bufio.Reader
	logBatchSize  int
}

func (r *Runner) Init() (err error) {
	// update task
	if err := r.updateTask("", nil); err != nil {
		return err
	}

	// start grpc client
	if !r.c.IsStarted() {
		r.c.Start()
	}

	// working directory
	workspacePath := viper.GetString("workspace")
	r.cwd = filepath.Join(workspacePath, r.s.GetId().Hex())

	// sync files from master
	if !utils.IsMaster() {
		if err := r.syncFiles(); err != nil {
			return err
		}
	}

	// grpc task service stream client
	if err := r.initSub(); err != nil {
		return err
	}

	// pre actions
	if r.hookSvc != nil {
		if err := r.hookSvc.PreActions(r.t, r.s, r.fsSvc, r.svc); err != nil {
			return err
		}
	}

	return nil
}

func (r *Runner) Run() (err error) {
	// log task started
	log.Infof("task[%s] started", r.tid.Hex())

	// configure cmd
	r.configureCmd()

	// configure environment variables
	r.configureEnv()

	// configure logging
	r.configureLogging()

	// start process
	if err := r.cmd.Start(); err != nil {
		return r.updateTask(constants.TaskStatusError, err)
	}

	// start logging
	go r.startLogging()

	// process id
	if r.cmd.Process == nil {
		return r.updateTask(constants.TaskStatusError, constants.ErrNotExists)
	}
	r.pid = r.cmd.Process.Pid
	r.t.SetPid(r.pid)

	// update task status (processing)
	if err := r.updateTask(constants.TaskStatusRunning, nil); err != nil {
		return err
	}

	// wait for process to finish
	go r.wait()

	// start health check
	go r.startHealthCheck()

	// declare task status
	status := ""

	// wait for signal
	signal := <-r.ch
	switch signal {
	case constants.TaskSignalFinish:
		err = nil
		status = constants.TaskStatusFinished
	case constants.TaskSignalCancel:
		err = constants.ErrTaskCancelled
		status = constants.TaskStatusCancelled
	case constants.TaskSignalError:
		err = r.err
		status = constants.TaskStatusError
	case constants.TaskSignalLost:
		err = constants.ErrTaskLost
		status = constants.TaskStatusError
	default:
		err = constants.ErrInvalidSignal
		status = constants.TaskStatusError
	}

	// update task status
	if err := r.updateTask(status, err); err != nil {
		return err
	}

	// post actions
	if r.hookSvc != nil {
		if err := r.hookSvc.PostActions(r.t, r.s, r.fsSvc, r.svc); err != nil {
			return err
		}
	}

	return err
}

func (r *Runner) Cancel() (err error) {
	// kill process
	opts := &sys_exec.KillProcessOptions{
		Timeout: r.svc.GetCancelTimeout(),
		Force:   true,
	}
	if err := sys_exec.KillProcess(r.cmd, opts); err != nil {
		return err
	}

	// make sure the process does not exist
	op := func() error {
		if exists, _ := process.PidExists(int32(r.pid)); exists {
			return errors.ErrorTaskProcessStillExists
		}
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), r.svc.GetExitWatchDuration())
	defer cancel()
	b := backoff.WithContext(backoff.NewConstantBackOff(1*time.Second), ctx)
	if err := backoff.Retry(op, b); err != nil {
		return trace.TraceError(errors.ErrorTaskUnableToCancel)
	}

	return nil
}

// CleanUp clean up task runner
func (r *Runner) CleanUp() (err error) {
	return nil
}

func (r *Runner) SetSubscribeTimeout(timeout time.Duration) {
	r.subscribeTimeout = timeout
}

func (r *Runner) GetTaskId() (id primitive.ObjectID) {
	return r.tid
}

func (r *Runner) configureCmd() {
	var cmdStr string

	// customized spider
	if r.t.GetCmd() == "" {
		cmdStr = r.s.GetCmd()
	} else {
		cmdStr = r.t.GetCmd()
	}

	// parameters
	if r.t.GetParam() != "" {
		cmdStr += " " + r.t.GetParam()
	} else if r.s.GetParam() != "" {
		cmdStr += " " + r.s.GetParam()
	}

	// get cmd instance
	r.cmd = sys_exec.BuildCmd(cmdStr)

	// set working directory
	r.cmd.Dir = r.cwd

	// configure pgid to allow killing sub processes
	//sys_exec.SetPgid(r.cmd)
}

func (r *Runner) configureLogging() {
	// set stdout reader
	stdout, _ := r.cmd.StdoutPipe()
	r.scannerStdout = bufio.NewReaderSize(stdout, r.bufferSize)

	// set stderr reader
	stderr, _ := r.cmd.StderrPipe()
	r.scannerStderr = bufio.NewReaderSize(stderr, r.bufferSize)
}

func (r *Runner) startLogging() {
	// start reading stdout
	go r.startLoggingReaderStdout()

	// start reading stderr
	go r.startLoggingReaderStderr()
}

func (r *Runner) startLoggingReaderStdout() {
	for {
		line, err := r.scannerStdout.ReadString(byte('\n'))
		if err != nil {
			break
		}
		line = strings.TrimSuffix(line, "\n")
		r.writeLogLines([]string{line})
	}
}

func (r *Runner) startLoggingReaderStderr() {
	for {
		line, err := r.scannerStderr.ReadString(byte('\n'))
		if err != nil {
			break
		}
		line = strings.TrimSuffix(line, "\n")
		r.writeLogLines([]string{line})
	}
}

func (r *Runner) startHealthCheck() {
	if r.cmd.ProcessState == nil || r.cmd.ProcessState.Exited() {
		return
	}
	for {
		exists, _ := process.PidExists(int32(r.pid))
		if !exists {
			// process lost
			r.ch <- constants.TaskSignalLost
			return
		}
		time.Sleep(1 * time.Second)
	}
}

func (r *Runner) configureEnv() {
	// 默认把Node.js的全局node_modules加入环境变量
	envPath := os.Getenv("PATH")
	nodePath := "/usr/lib/node_modules"
	if !strings.Contains(envPath, nodePath) {
		_ = os.Setenv("PATH", nodePath+":"+envPath)
	}
	_ = os.Setenv("NODE_PATH", nodePath)

	// default envs
	r.cmd.Env = append(os.Environ(), "CRAWLAB_TASK_ID="+r.tid.Hex())
	if viper.GetString("grpc.address") != "" {
		r.cmd.Env = append(r.cmd.Env, "CRAWLAB_GRPC_ADDRESS="+viper.GetString("grpc.address"))
	}
	if viper.GetString("grpc.authKey") != "" {
		r.cmd.Env = append(r.cmd.Env, "CRAWLAB_GRPC_AUTH_KEY="+viper.GetString("grpc.authKey"))
	} else {
		r.cmd.Env = append(r.cmd.Env, "CRAWLAB_GRPC_AUTH_KEY="+constants.DefaultGrpcAuthKey)
	}

	// global environment variables
	envs, err := r.svc.GetModelEnvironmentService().GetEnvironmentList(nil, nil)
	if err != nil {
		trace.PrintError(err)
		return
	}
	for _, env := range envs {
		r.cmd.Env = append(r.cmd.Env, env.GetKey()+"="+env.GetValue())
	}
}

func (r *Runner) syncFiles() (err error) {
	masterURL := fmt.Sprintf("%s/sync/%s", viper.GetString("api.endpoint"), r.s.GetId().Hex())
	workspacePath := viper.GetString("workspace")
	workerDir := filepath.Join(workspacePath, r.s.GetId().Hex())

	// get file list from master
	resp, err := http.Get(masterURL + "/scan")
	if err != nil {
		fmt.Println("Error getting file list from master:", err)
		return trace.TraceError(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return trace.TraceError(err)
	}
	var masterFiles map[string]entity.FsFileInfo
	err = json.Unmarshal(body, &masterFiles)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return trace.TraceError(err)
	}

	// create a map for master files
	masterFilesMap := make(map[string]entity.FsFileInfo)
	for _, file := range masterFiles {
		masterFilesMap[file.Path] = file
	}

	// create worker directory if not exists
	if _, err := os.Stat(workerDir); os.IsNotExist(err) {
		if err := os.MkdirAll(workerDir, os.ModePerm); err != nil {
			fmt.Println("Error creating worker directory:", err)
			return trace.TraceError(err)
		}
	}

	// get file list from worker
	workerFiles, err := utils.ScanDirectory(workerDir)
	if err != nil {
		fmt.Println("Error scanning worker directory:", err)
		return trace.TraceError(err)
	}

	// set up wait group and error channel
	var wg sync.WaitGroup
	errCh := make(chan error, 1)

	// delete files that are deleted on master node
	for path, workerFile := range workerFiles {
		if _, exists := masterFilesMap[path]; !exists {
			fmt.Println("Deleting file:", path)
			err := os.Remove(workerFile.FullPath)
			if err != nil {
				fmt.Println("Error deleting file:", err)
			}
		}
	}

	// download files that are new or modified on master node
	for path, masterFile := range masterFilesMap {
		workerFile, exists := workerFiles[path]
		if !exists || masterFile.Hash != workerFile.Hash {
			wg.Add(1)
			go func(path string, masterFile entity.FsFileInfo) {
				defer wg.Done()
				logrus.Infof("File needs to be synchronized: %s", path)
				err := r.downloadFile(masterURL+"/download?path="+path, filepath.Join(workerDir, path))
				if err != nil {
					logrus.Errorf("Error downloading file: %v", err)
					select {
					case errCh <- err:
					default:
					}
				}
			}(path, masterFile)
		}
	}

	wg.Wait()
	close(errCh)
	if err := <-errCh; err != nil {
		return err
	}

	return nil
}

func (r *Runner) downloadFile(url string, filePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// wait for process to finish and send task signal (constants.TaskSignal)
// to task runner's channel (Runner.ch) according to exit code
func (r *Runner) wait() {
	// wait for process to finish
	if err := r.cmd.Wait(); err != nil {
		exitError, ok := err.(*exec.ExitError)
		if !ok {
			r.ch <- constants.TaskSignalError
			return
		}
		exitCode := exitError.ExitCode()
		if exitCode == -1 {
			// cancel error
			r.ch <- constants.TaskSignalCancel
			return
		}

		// standard error
		r.err = err
		r.ch <- constants.TaskSignalError
		return
	}

	// success
	r.ch <- constants.TaskSignalFinish
}

// updateTask update and get updated info of task (Runner.t)
func (r *Runner) updateTask(status string, e error) (err error) {
	if r.t != nil && status != "" {
		// update task status
		r.t.SetStatus(status)
		if e != nil {
			r.t.SetError(e.Error())
		}
		if r.svc.GetNodeConfigService().IsMaster() {
			if err := delegate.NewModelDelegate(r.t).Save(); err != nil {
				return err
			}
		} else {
			if err := client.NewModelDelegate(r.t, client.WithDelegateConfigPath(r.svc.GetConfigPath())).Save(); err != nil {
				return err
			}
		}

		// send notification
		go r.sendNotification()

		// update stats
		go func() {
			r._updateTaskStat(status)
			r._updateSpiderStat(status)
		}()
	}

	// get task
	r.t, err = r.svc.GetTaskById(r.tid)
	if err != nil {
		return err
	}

	return nil
}

func (r *Runner) initSub() (err error) {
	r.sub, err = r.c.GetTaskClient().Subscribe(context.Background())
	if err != nil {
		return trace.TraceError(err)
	}
	return nil
}

func (r *Runner) writeLogLines(lines []string) {
	data, err := json.Marshal(&entity.StreamMessageTaskData{
		TaskId: r.tid,
		Logs:   lines,
	})
	if err != nil {
		trace.PrintError(err)
		return
	}
	msg := &grpc.StreamMessage{
		Code: grpc.StreamMessageCode_INSERT_LOGS,
		Data: data,
	}
	if err := r.sub.Send(msg); err != nil {
		trace.PrintError(err)
		return
	}
}

func (r *Runner) _updateTaskStat(status string) {
	ts, err := r.svc.GetModelTaskStatService().GetTaskStatById(r.tid)
	if err != nil {
		trace.PrintError(err)
		return
	}
	switch status {
	case constants.TaskStatusPending:
		// do nothing
	case constants.TaskStatusRunning:
		ts.SetStartTs(time.Now())
		ts.SetWaitDuration(ts.GetStartTs().Sub(ts.GetCreateTs()).Milliseconds())
	case constants.TaskStatusFinished, constants.TaskStatusError, constants.TaskStatusCancelled:
		ts.SetEndTs(time.Now())
		ts.SetRuntimeDuration(ts.GetEndTs().Sub(ts.GetStartTs()).Milliseconds())
		ts.SetTotalDuration(ts.GetEndTs().Sub(ts.GetCreateTs()).Milliseconds())
	}
	if r.svc.GetNodeConfigService().IsMaster() {
		if err := delegate.NewModelDelegate(ts).Save(); err != nil {
			trace.PrintError(err)
			return
		}
	} else {
		if err := client.NewModelDelegate(ts, client.WithDelegateConfigPath(r.svc.GetConfigPath())).Save(); err != nil {
			trace.PrintError(err)
			return
		}
	}
}

func (r *Runner) sendNotification() {
	data, err := json.Marshal(r.t)
	if err != nil {
		trace.PrintError(err)
		return
	}
	req := &grpc.Request{
		NodeKey: r.svc.GetNodeConfigService().GetNodeKey(),
		Data:    data,
	}
	_, err = r.c.GetTaskClient().SendNotification(context.Background(), req)
	if err != nil {
		trace.PrintError(err)
		return
	}
}

func (r *Runner) _updateSpiderStat(status string) {
	// task stat
	ts, err := r.svc.GetModelTaskStatService().GetTaskStatById(r.tid)
	if err != nil {
		trace.PrintError(err)
		return
	}

	// update
	var update bson.M
	switch status {
	case constants.TaskStatusPending, constants.TaskStatusRunning:
		update = bson.M{
			"$set": bson.M{
				"last_task_id": r.tid, // last task id
			},
			"$inc": bson.M{
				"tasks":         1,                    // task count
				"wait_duration": ts.GetWaitDuration(), // wait duration
			},
		}
	case constants.TaskStatusFinished, constants.TaskStatusError, constants.TaskStatusCancelled:
		update = bson.M{
			"$inc": bson.M{
				"results":          ts.GetResultCount(),            // results
				"runtime_duration": ts.GetRuntimeDuration() / 1000, // runtime duration
				"total_duration":   ts.GetTotalDuration() / 1000,   // total duration
			},
		}
	default:
		trace.PrintError(errors.ErrorTaskInvalidType)
		return
	}

	// perform update
	if r.svc.GetNodeConfigService().IsMaster() {
		if err := mongo.GetMongoCol(interfaces.ModelColNameSpiderStat).UpdateId(r.s.GetId(), update); err != nil {
			trace.PrintError(err)
			return
		}
	} else {
		modelSvc, err := client.NewBaseServiceDelegate(
			client.WithBaseServiceModelId(interfaces.ModelIdSpiderStat),
			client.WithBaseServiceConfigPath(r.svc.GetConfigPath()),
		)
		if err != nil {
			trace.PrintError(err)
			return
		}
		if err := modelSvc.UpdateById(r.s.GetId(), update); err != nil {
			trace.PrintError(err)
			return
		}
	}

}

func NewTaskRunner(id primitive.ObjectID, svc interfaces.TaskHandlerService, opts ...RunnerOption) (r2 interfaces.TaskRunner, err error) {
	// validate options
	if id.IsZero() {
		return nil, constants.ErrInvalidOptions
	}

	// runner
	r := &Runner{
		subscribeTimeout: 30 * time.Second,
		bufferSize:       1024 * 1024,
		svc:              svc,
		tid:              id,
		ch:               make(chan constants.TaskSignal),
		logBatchSize:     20,
	}

	// apply options
	for _, opt := range opts {
		opt(r)
	}

	// task
	r.t, err = svc.GetTaskById(id)
	if err != nil {
		return nil, err
	}

	// spider
	r.s, err = svc.GetSpiderById(r.t.GetSpiderId())
	if err != nil {
		return nil, err
	}

	// task fs service
	r.fsSvc = fs2.NewFsServiceV2(filepath.Join(viper.GetString("workspace"), r.s.GetId().Hex()))

	// dependency injection
	if err := container.GetContainer().Invoke(func(
		c interfaces.GrpcClient,
	) {
		r.c = c
	}); err != nil {
		return nil, trace.TraceError(err)
	}

	_ = container.GetContainer().Invoke(func(hookSvc interfaces.TaskHookService) {
		r.hookSvc = hookSvc
	})

	// initialize task runner
	if err := r.Init(); err != nil {
		return r, err
	}

	return r, nil
}
