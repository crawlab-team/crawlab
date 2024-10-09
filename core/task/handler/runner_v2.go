package handler

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/apex/log"
	"github.com/cenkalti/backoff/v4"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/entity"
	fs2 "github.com/crawlab-team/crawlab/core/fs"
	client2 "github.com/crawlab-team/crawlab/core/grpc/client"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/client"
	"github.com/crawlab-team/crawlab/core/models/models"
	models2 "github.com/crawlab-team/crawlab/core/models/models/v2"
	service2 "github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/core/sys_exec"
	"github.com/crawlab-team/crawlab/core/utils"
	"github.com/crawlab-team/crawlab/grpc"
	"github.com/crawlab-team/crawlab/trace"
	"github.com/shirou/gopsutil/process"
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

type RunnerV2 struct {
	// dependencies
	svc   *ServiceV2             // task handler service
	fsSvc interfaces.FsServiceV2 // task fs service

	// settings
	subscribeTimeout time.Duration
	bufferSize       int

	// internals
	cmd  *exec.Cmd                        // process command instance
	pid  int                              // process id
	tid  primitive.ObjectID               // task id
	t    *models2.TaskV2                  // task model.Task
	s    *models2.SpiderV2                // spider model.Spider
	ch   chan constants.TaskSignal        // channel to communicate between Service and RunnerV2
	err  error                            // standard process error
	envs []models.Env                     // environment variables
	cwd  string                           // working directory
	c    *client2.GrpcClientV2            // grpc client
	sub  grpc.TaskService_SubscribeClient // grpc task service stream client

	// log internals
	scannerStdout *bufio.Reader
	scannerStderr *bufio.Reader
	logBatchSize  int
}

func (r *RunnerV2) Init() (err error) {
	// update task
	if err := r.updateTask("", nil); err != nil {
		return err
	}

	// start grpc client
	if !r.c.IsStarted() {
		err := r.c.Start()
		if err != nil {
			return err
		}
	}

	// grpc task service stream client
	if err := r.initSub(); err != nil {
		return err
	}

	return nil
}

func (r *RunnerV2) Run() (err error) {
	// log task started
	log.Infof("task[%s] started", r.tid.Hex())

	// configure working directory
	r.configureCwd()

	// sync files worker nodes
	if !utils.IsMaster() {
		if err := r.syncFiles(); err != nil {
			return r.updateTask(constants.TaskStatusError, err)
		}
	}

	// configure cmd
	err = r.configureCmd()
	if err != nil {
		return r.updateTask(constants.TaskStatusError, err)
	}

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
	r.t.Pid = r.pid

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

	return err
}

func (r *RunnerV2) Cancel() (err error) {
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
			return errors.New(fmt.Sprintf("task process %d still exists", r.pid))
		}
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), r.svc.GetExitWatchDuration())
	defer cancel()
	b := backoff.WithContext(backoff.NewConstantBackOff(1*time.Second), ctx)
	if err := backoff.Retry(op, b); err != nil {
		log.Errorf("Error canceling task %s: %v", r.tid, err)
		return trace.TraceError(err)
	}

	return nil
}

// CleanUp clean up task runner
func (r *RunnerV2) CleanUp() (err error) {
	return nil
}

func (r *RunnerV2) SetSubscribeTimeout(timeout time.Duration) {
	r.subscribeTimeout = timeout
}

func (r *RunnerV2) GetTaskId() (id primitive.ObjectID) {
	return r.tid
}

func (r *RunnerV2) configureCmd() (err error) {
	var cmdStr string

	// customized spider
	if r.t.Cmd == "" {
		cmdStr = r.s.Cmd
	} else {
		cmdStr = r.t.Cmd
	}

	// parameters
	if r.t.Param != "" {
		cmdStr += " " + r.t.Param
	} else if r.s.Param != "" {
		cmdStr += " " + r.s.Param
	}

	// get cmd instance
	r.cmd, err = sys_exec.BuildCmd(cmdStr)
	if err != nil {
		log.Errorf("Error building command: %v", err)
		trace.PrintError(err)
		return err
	}

	// set working directory
	r.cmd.Dir = r.cwd

	return nil
}

func (r *RunnerV2) configureLogging() {
	// set stdout reader
	stdout, _ := r.cmd.StdoutPipe()
	r.scannerStdout = bufio.NewReaderSize(stdout, r.bufferSize)

	// set stderr reader
	stderr, _ := r.cmd.StderrPipe()
	r.scannerStderr = bufio.NewReaderSize(stderr, r.bufferSize)
}

func (r *RunnerV2) startLogging() {
	// start reading stdout
	go r.startLoggingReaderStdout()

	// start reading stderr
	go r.startLoggingReaderStderr()
}

func (r *RunnerV2) startLoggingReaderStdout() {
	for {
		line, err := r.scannerStdout.ReadString(byte('\n'))
		if err != nil {
			break
		}
		line = strings.TrimSuffix(line, "\n")
		r.writeLogLines([]string{line})
	}
}

func (r *RunnerV2) startLoggingReaderStderr() {
	for {
		line, err := r.scannerStderr.ReadString(byte('\n'))
		if err != nil {
			break
		}
		line = strings.TrimSuffix(line, "\n")
		r.writeLogLines([]string{line})
	}
}

func (r *RunnerV2) startHealthCheck() {
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

func (r *RunnerV2) configureEnv() {
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
	envs, err := client.NewModelServiceV2[models2.EnvironmentV2]().GetMany(nil, nil)
	if err != nil {
		trace.PrintError(err)
		return
	}
	for _, env := range envs {
		r.cmd.Env = append(r.cmd.Env, env.Key+"="+env.Value)
	}
}

func (r *RunnerV2) syncFiles() (err error) {
	var id string
	var workingDir string
	if r.s.GitId.IsZero() {
		id = r.s.Id.Hex()
		workingDir = ""
	} else {
		id = r.s.GitId.Hex()
		workingDir = r.s.GitRootPath
	}
	masterURL := fmt.Sprintf("%s/sync/%s", viper.GetString("api.endpoint"), id)

	// get file list from master
	resp, err := http.Get(masterURL + "/scan?path=" + workingDir)
	if err != nil {
		log.Errorf("Error getting file list from master: %v", err)
		return trace.TraceError(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Error reading response body: %v", err)
		return trace.TraceError(err)
	}
	var masterFiles map[string]entity.FsFileInfo
	err = json.Unmarshal(body, &masterFiles)
	if err != nil {
		log.Errorf("Error unmarshaling JSON: %v", err)
		return trace.TraceError(err)
	}

	// create a map for master files
	masterFilesMap := make(map[string]entity.FsFileInfo)
	for _, file := range masterFiles {
		masterFilesMap[file.Path] = file
	}

	// create working directory if not exists
	if _, err := os.Stat(r.cwd); os.IsNotExist(err) {
		if err := os.MkdirAll(r.cwd, os.ModePerm); err != nil {
			log.Errorf("Error creating worker directory: %v", err)
			return trace.TraceError(err)
		}
	}

	// get file list from worker
	workerFiles, err := utils.ScanDirectory(r.cwd)
	if err != nil {
		log.Errorf("Error scanning worker directory: %v", err)
		return trace.TraceError(err)
	}

	// delete files that are deleted on master node
	for path, workerFile := range workerFiles {
		if _, exists := masterFilesMap[path]; !exists {
			log.Infof("Deleting file: %s", path)
			err := os.Remove(workerFile.FullPath)
			if err != nil {
				log.Errorf("Error deleting file: %v", err)
			}
		}
	}

	// set up wait group and error channel
	var wg sync.WaitGroup
	pool := make(chan struct{}, 10)

	// download files that are new or modified on master node
	for path, masterFile := range masterFilesMap {
		workerFile, exists := workerFiles[path]
		if !exists || masterFile.Hash != workerFile.Hash {
			wg.Add(1)

			// acquire token
			pool <- struct{}{}

			// start goroutine to synchronize file or directory
			go func(path string, masterFile *entity.FsFileInfo) {
				defer wg.Done()

				if masterFile.IsDir {
					log.Infof("Directory needs to be synchronized: %s", path)
					_err := os.MkdirAll(filepath.Join(r.cwd, path), masterFile.Mode)
					if _err != nil {
						log.Errorf("Error creating directory: %v", _err)
						err = errors.Join(err, _err)
					}
				} else {
					log.Infof("File needs to be synchronized: %s", path)
					_err := r.downloadFile(masterURL+"/download?path="+filepath.Join(workingDir, path), filepath.Join(r.cwd, path), masterFile)
					if _err != nil {
						log.Errorf("Error downloading file: %v", _err)
						err = errors.Join(err, _err)
					}
				}

				// release token
				<-pool

			}(path, &masterFile)
		}
	}

	// wait for all files and directories to be synchronized
	wg.Wait()

	return err
}

func (r *RunnerV2) downloadFile(url string, filePath string, fileInfo *entity.FsFileInfo) error {
	// get file response
	resp, err := http.Get(url)
	if err != nil {
		log.Errorf("Error getting file response: %v", err)
		return err
	}
	if resp.StatusCode != http.StatusOK {
		log.Errorf("Error downloading file: %s", resp.Status)
		return errors.New(resp.Status)
	}
	defer resp.Body.Close()

	// create directory if not exists
	dirPath := filepath.Dir(filePath)
	utils.Exists(dirPath)
	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		log.Errorf("Error creating directory: %v", err)
		return err
	}

	// create local file
	out, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fileInfo.Mode)
	if err != nil {
		log.Errorf("Error creating file: %v", err)
		return err
	}
	defer out.Close()

	// copy file content to local file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Errorf("Error copying file: %v", err)
		return err
	}
	return nil
}

// wait for process to finish and send task signal (constants.TaskSignal)
// to task runner's channel (RunnerV2.ch) according to exit code
func (r *RunnerV2) wait() {
	// wait for process to finish
	if err := r.cmd.Wait(); err != nil {
		var exitError *exec.ExitError
		ok := errors.As(err, &exitError)
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

// updateTask update and get updated info of task (RunnerV2.t)
func (r *RunnerV2) updateTask(status string, e error) (err error) {
	if r.t != nil && status != "" {
		// update task status
		r.t.Status = status
		if e != nil {
			r.t.Error = e.Error()
		}
		if r.svc.GetNodeConfigService().IsMaster() {
			err = service2.NewModelServiceV2[models2.TaskV2]().ReplaceById(r.t.Id, *r.t)
			if err != nil {
				return err
			}
		} else {
			err = client.NewModelServiceV2[models2.TaskV2]().ReplaceById(r.t.Id, *r.t)
			if err != nil {
				return err
			}
		}

		// update stats
		r._updateTaskStat(status)
		r._updateSpiderStat(status)

		// send notification
		go r.sendNotification()
	}

	// get task
	r.t, err = r.svc.GetTaskById(r.tid)
	if err != nil {
		return err
	}

	return nil
}

func (r *RunnerV2) initSub() (err error) {
	r.sub, err = r.c.TaskClient.Subscribe(context.Background())
	if err != nil {
		return trace.TraceError(err)
	}
	return nil
}

func (r *RunnerV2) writeLogLines(lines []string) {
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

func (r *RunnerV2) _updateTaskStat(status string) {
	ts, err := client.NewModelServiceV2[models2.TaskStatV2]().GetById(r.tid)
	if err != nil {
		trace.PrintError(err)
		return
	}
	switch status {
	case constants.TaskStatusPending:
		// do nothing
	case constants.TaskStatusRunning:
		ts.StartTs = time.Now()
		ts.WaitDuration = ts.StartTs.Sub(ts.BaseModelV2.CreatedAt).Milliseconds()
	case constants.TaskStatusFinished, constants.TaskStatusError, constants.TaskStatusCancelled:
		if ts.StartTs.IsZero() {
			ts.StartTs = time.Now()
			ts.WaitDuration = ts.StartTs.Sub(ts.BaseModelV2.CreatedAt).Milliseconds()
		}
		ts.EndTs = time.Now()
		ts.RuntimeDuration = ts.EndTs.Sub(ts.StartTs).Milliseconds()
		ts.TotalDuration = ts.EndTs.Sub(ts.BaseModelV2.CreatedAt).Milliseconds()
	}
	if r.svc.GetNodeConfigService().IsMaster() {
		err = service2.NewModelServiceV2[models2.TaskStatV2]().ReplaceById(ts.Id, *ts)
		if err != nil {
			trace.PrintError(err)
			return
		}
	} else {
		err = client.NewModelServiceV2[models2.TaskStatV2]().ReplaceById(ts.Id, *ts)
		if err != nil {
			trace.PrintError(err)
			return
		}
	}
}

func (r *RunnerV2) sendNotification() {
	req := &grpc.TaskServiceSendNotificationRequest{
		NodeKey: r.svc.GetNodeConfigService().GetNodeKey(),
		TaskId:  r.tid.Hex(),
	}
	_, err := r.c.TaskClient.SendNotification(context.Background(), req)
	if err != nil {
		log.Errorf("Error sending notification: %v", err)
		trace.PrintError(err)
		return
	}
}

func (r *RunnerV2) _updateSpiderStat(status string) {
	// task stat
	ts, err := client.NewModelServiceV2[models2.TaskStatV2]().GetById(r.tid)
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
				"tasks":         1,               // task count
				"wait_duration": ts.WaitDuration, // wait duration
			},
		}
	case constants.TaskStatusFinished, constants.TaskStatusError, constants.TaskStatusCancelled:
		update = bson.M{
			"$set": bson.M{
				"last_task_id": r.tid, // last task id
			},
			"$inc": bson.M{
				"results":          ts.ResultCount,            // results
				"runtime_duration": ts.RuntimeDuration / 1000, // runtime duration
				"total_duration":   ts.TotalDuration / 1000,   // total duration
			},
		}
	default:
		log.Errorf("Invalid task status: %s", status)
		trace.PrintError(errors.New("invalid task status"))
		return
	}

	// perform update
	if r.svc.GetNodeConfigService().IsMaster() {
		err = service2.NewModelServiceV2[models2.SpiderStatV2]().UpdateById(r.s.Id, update)
		if err != nil {
			trace.PrintError(err)
			return
		}
	} else {
		err = client.NewModelServiceV2[models2.SpiderStatV2]().UpdateById(r.s.Id, update)
		if err != nil {
			trace.PrintError(err)
			return
		}
	}
}

func (r *RunnerV2) configureCwd() {
	workspacePath := viper.GetString("workspace")
	if r.s.GitId.IsZero() {
		// not git
		r.cwd = filepath.Join(workspacePath, r.s.Id.Hex())
	} else {
		// git
		r.cwd = filepath.Join(workspacePath, r.s.GitId.Hex(), r.s.GitRootPath)
	}
}

func NewTaskRunnerV2(id primitive.ObjectID, svc *ServiceV2) (r2 *RunnerV2, err error) {
	// validate options
	if id.IsZero() {
		return nil, constants.ErrInvalidOptions
	}

	// runner
	r := &RunnerV2{
		subscribeTimeout: 30 * time.Second,
		bufferSize:       1024 * 1024,
		svc:              svc,
		tid:              id,
		ch:               make(chan constants.TaskSignal),
		logBatchSize:     20,
	}

	// task
	r.t, err = svc.GetTaskById(id)
	if err != nil {
		return nil, err
	}

	// spider
	r.s, err = svc.GetSpiderById(r.t.SpiderId)
	if err != nil {
		return nil, err
	}

	// task fs service
	r.fsSvc = fs2.NewFsServiceV2(filepath.Join(viper.GetString("workspace"), r.s.Id.Hex()))

	// grpc client
	r.c = client2.GetGrpcClientV2()

	// initialize task runner
	if err := r.Init(); err != nil {
		return r, err
	}

	return r, nil
}
