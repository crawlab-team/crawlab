package test

import (
	"github.com/crawlab-team/crawlab/core/entity"
	"github.com/crawlab-team/crawlab/core/interfaces"
	service2 "github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/core/node/service"
	"github.com/crawlab-team/crawlab/core/utils"
	"github.com/spf13/viper"
	"go.uber.org/dig"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"
)

func init() {
	var err error
	T, err = NewTest()
	if err != nil {
		panic(err)
	}
}

var T *Test

type Test struct {
	DefaultSvc       interfaces.NodeMasterService
	MasterSvc        interfaces.NodeMasterService
	WorkerSvc        interfaces.NodeWorkerService
	MasterSvcMonitor interfaces.NodeMasterService
	WorkerSvcMonitor interfaces.NodeWorkerService
	ModelSvc         service2.ModelService
}

func NewTest() (res *Test, err error) {
	// test
	t := &Test{}

	// recreate config directory path
	_ = os.RemoveAll(viper.GetString("metadata"))
	_ = os.MkdirAll(viper.GetString("metadata"), os.FileMode(0766))

	// master config and settings
	masterNodeConfigName := "config-master.json"
	masterNodeConfigPath := path.Join(viper.GetString("metadata"), masterNodeConfigName)
	if err := ioutil.WriteFile(masterNodeConfigPath, []byte("{\"key\":\"master\",\"is_master\":true}"), os.FileMode(0766)); err != nil {
		return nil, err
	}
	masterHost := "0.0.0.0"
	masterPort := "9667"

	// worker config and settings
	workerNodeConfigName := "config-worker.json"
	workerNodeConfigPath := path.Join(viper.GetString("metadata"), workerNodeConfigName)
	if err = ioutil.WriteFile(workerNodeConfigPath, []byte("{\"key\":\"worker\",\"is_master\":false}"), os.FileMode(0766)); err != nil {
		return nil, err
	}
	workerHost := "localhost"
	workerPort := masterPort

	// master for monitor config and settings
	masterNodeMonitorConfigName := "config-master-monitor.json"
	masterNodeMonitorConfigPath := path.Join(viper.GetString("metadata"), masterNodeMonitorConfigName)
	if err := ioutil.WriteFile(masterNodeMonitorConfigPath, []byte("{\"key\":\"master-monitor\",\"is_master\":true}"), os.FileMode(0766)); err != nil {
		return nil, err
	}
	masterMonitorHost := masterHost
	masterMonitorPort := "9668"

	// worker for monitor config and settings
	workerNodeMonitorConfigName := "config-worker-monitor.json"
	workerNodeMonitorConfigPath := path.Join(viper.GetString("metadata"), workerNodeMonitorConfigName)
	if err := ioutil.WriteFile(workerNodeMonitorConfigPath, []byte("{\"key\":\"worker-monitor\",\"is_master\":false}"), os.FileMode(0766)); err != nil {
		return nil, err
	}
	workerMonitorHost := workerHost
	workerMonitorPort := masterMonitorPort

	// dependency injection
	c := dig.New()
	if err := c.Provide(service.ProvideMasterService(
		masterNodeConfigPath,
		service.WithMonitorInterval(3*time.Second),
		service.WithAddress(entity.NewAddress(&entity.AddressOptions{
			Host: masterHost,
			Port: masterPort,
		})),
	)); err != nil {
		return nil, err
	}
	if err := c.Provide(service.ProvideWorkerService(
		workerNodeConfigPath,
		service.WithHeartbeatInterval(1*time.Second),
		service.WithAddress(entity.NewAddress(&entity.AddressOptions{
			Host: workerHost,
			Port: workerPort,
		})),
	)); err != nil {
		return nil, err
	}
	if err := c.Provide(service2.NewService); err != nil {
		return nil, err
	}
	if err := c.Invoke(func(masterSvc interfaces.NodeMasterService, workerSvc interfaces.NodeWorkerService, modelSvc service2.ModelService) {
		t.MasterSvc = masterSvc
		t.WorkerSvc = workerSvc
		t.ModelSvc = modelSvc
	}); err != nil {
		return nil, err
	}

	// default service
	t.DefaultSvc, err = service.NewMasterService()
	if err != nil {
		return nil, err
	}

	// master and worker for monitor
	t.MasterSvcMonitor, err = service.NewMasterService(
		service.WithConfigPath(masterNodeMonitorConfigPath),
		service.WithAddress(entity.NewAddress(&entity.AddressOptions{
			Host: masterMonitorHost,
			Port: masterMonitorPort,
		})),
		service.WithMonitorInterval(3*time.Second),
		service.WithStopOnError(),
	)
	if err != nil {
		return nil, err
	}
	t.WorkerSvcMonitor, err = service.NewWorkerService(
		service.WithConfigPath(workerNodeMonitorConfigPath),
		service.WithAddress(entity.NewAddress(&entity.AddressOptions{
			Host: workerMonitorHost,
			Port: workerMonitorPort,
		})),
		service.WithHeartbeatInterval(1*time.Second),
		service.WithStopOnError(),
	)
	if err != nil {
		return nil, err
	}

	// removed all data in db
	_ = t.ModelSvc.DropAll()

	// visualize dependencies
	if err := utils.VisualizeContainer(c); err != nil {
		return nil, err
	}

	return t, nil
}

func (t *Test) Setup(t2 *testing.T) {
	if err := t.ModelSvc.DropAll(); err != nil {
		panic(err)
	}
	_ = os.RemoveAll(viper.GetString("metadata"))
	t2.Cleanup(t.Cleanup)
}

func (t *Test) Cleanup() {
	if err := t.ModelSvc.DropAll(); err != nil {
		panic(err)
	}
	_ = os.RemoveAll(viper.GetString("metadata"))
}

func (t *Test) StartMasterWorker() {
	startMasterWorker()
}

func (t *Test) StopMasterWorker() {
	stopMasterWorker()
}

func startMasterWorker() {
	go T.MasterSvc.Start()
	time.Sleep(1 * time.Second)
	go T.WorkerSvc.Start()
	time.Sleep(1 * time.Second)
}

func stopMasterWorker() {
	go T.MasterSvc.Stop()
	time.Sleep(1 * time.Second)
	go T.WorkerSvc.Stop()
	time.Sleep(1 * time.Second)
}

func startMasterWorkerMonitor() {
	go T.MasterSvcMonitor.Start()
	time.Sleep(1 * time.Second)
	go T.WorkerSvcMonitor.Start()
	time.Sleep(1 * time.Second)
}

func stopMasterWorkerMonitor() {
	go T.MasterSvcMonitor.Stop()
	time.Sleep(1 * time.Second)
	go T.WorkerSvcMonitor.Stop()
	time.Sleep(1 * time.Second)
}
