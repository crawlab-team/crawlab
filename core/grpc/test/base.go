package test

import (
	"github.com/crawlab-team/crawlab/core/entity"
	"github.com/crawlab-team/crawlab/core/grpc/client"
	"github.com/crawlab-team/crawlab/core/grpc/server"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/node/test"
	"testing"
	"time"
)

type Test struct {
	Server interfaces.GrpcServer
	Client interfaces.GrpcClient

	MasterNodeInfo *entity.NodeInfo
	WorkerNodeInfo *entity.NodeInfo
}

func (t *Test) Setup(t2 *testing.T) {
	test.T.Cleanup()
	t2.Cleanup(t.Cleanup)
	if !T.Client.IsStarted() {
		_ = T.Client.Start()
	} else if T.Client.IsClosed() {
		_ = T.Client.Restart()
	}
}

func (t *Test) Cleanup() {
	_ = t.Client.Stop()
	_ = t.Server.Stop()
	test.T.Cleanup()

	// wait to avoid caching
	time.Sleep(200 * time.Millisecond)
}

var T *Test

func NewTest() (res *Test, err error) {
	// test
	t := &Test{}

	// server
	t.Server, err = server.NewServer(
		server.WithConfigPath(test.T.MasterSvc.GetConfigPath()),
		server.WithAddress(test.T.MasterSvc.GetAddress()),
	)
	if err != nil {
		return nil, err
	}
	if err := t.Server.Start(); err != nil {
		return nil, err
	}

	// client
	t.Client, err = client.GetClient(test.T.WorkerSvc.GetConfigPath())
	if err != nil {
		return nil, err
	}

	// master node info
	t.MasterNodeInfo = &entity.NodeInfo{
		Key:      "master",
		IsMaster: true,
	}

	// worker node info
	t.WorkerNodeInfo = &entity.NodeInfo{
		Key:      "worker",
		IsMaster: false,
	}

	return t, nil
}
