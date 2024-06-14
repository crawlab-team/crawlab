package test

import (
	"context"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/models/delegate"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/node/test"
	grpc "github.com/crawlab-team/crawlab/grpc"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGrpcServer_Register(t *testing.T) {
	var err error

	T, _ = NewTest()
	T.Setup(t)

	// register
	register(t)

	// validate
	workerNodeKey := T.WorkerNodeInfo.Key
	workerNode, err := test.T.ModelSvc.GetNodeByKey(workerNodeKey, nil)
	require.Nil(t, err)
	require.Equal(t, workerNodeKey, workerNode.Key)
	require.Equal(t, constants.NodeStatusRegistered, workerNode.Status)
}

func TestGrpcServer_Register_Existing(t *testing.T) {
	var err error

	T, _ = NewTest()
	T.Setup(t)

	// add to db
	node := &models.Node{
		Key:      T.WorkerNodeInfo.Key,
		IsMaster: false,
		Status:   constants.NodeStatusUnregistered,
	}
	nodeD := delegate.NewModelDelegate(node)
	err = nodeD.Add()
	require.Nil(t, err)

	// register
	register(t)

	// validate
	workerNodeKey := T.WorkerNodeInfo.Key
	workerNode, err := test.T.ModelSvc.GetNodeByKey(workerNodeKey, nil)
	require.Nil(t, err)
	require.Equal(t, workerNodeKey, workerNode.Key)
	require.Equal(t, constants.NodeStatusRegistered, workerNode.Status)
}

func TestGrpcServer_SendHeartbeat(t *testing.T) {
	var err error

	T, _ = NewTest()
	T.Setup(t)

	// register
	register(t)

	// send heartbeat
	sendHeartbeat(t)

	// validate
	workerNodeKey := T.WorkerNodeInfo.Key
	workerNode, err := test.T.ModelSvc.GetNodeByKey(workerNodeKey, nil)
	require.Nil(t, err)
	require.Equal(t, workerNodeKey, workerNode.Key)
	require.Equal(t, constants.NodeStatusOnline, workerNode.Status)
}

func TestGrpcServer_Subscribe(t *testing.T) {
	var err error

	T, _ = NewTest()
	T.Setup(t)

	// register
	register(t)

	// handle client message
	go handleClientMessage(t)

	time.Sleep(1 * time.Second)

	// server PING client
	sub, err := T.Server.GetSubscribe("node:" + T.WorkerNodeInfo.Key)
	require.Nil(t, err)
	require.NotNil(t, sub)
	err = sub.GetStream().Send(&grpc.StreamMessage{
		Code:    grpc.StreamMessageCode_PING,
		NodeKey: T.MasterNodeInfo.Key,
	})
	require.Nil(t, err)

	// wait
	time.Sleep(1 * time.Second)

	// validate
	workerNode, err := test.T.ModelSvc.GetNodeByKey(T.WorkerNodeInfo.Key, nil)
	require.Nil(t, err)
	require.Equal(t, constants.NodeStatusOnline, workerNode.Status)
}

func register(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := T.Client.GetNodeClient().Register(ctx, T.Client.NewRequest(T.WorkerNodeInfo))
	require.Nil(t, err)
	require.Equal(t, grpc.ResponseCode_OK, res.Code)
}

func sendHeartbeat(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := T.Client.GetNodeClient().SendHeartbeat(ctx, T.Client.NewRequest(T.WorkerNodeInfo))
	require.Nil(t, err)
	require.Equal(t, grpc.ResponseCode_OK, res.Code)
}

func handleClientMessage(t *testing.T) {
	msgCh := T.Client.GetMessageChannel()
	for {
		msg := <-msgCh
		switch msg.Code {
		case grpc.StreamMessageCode_PING:
			require.Equal(t, T.MasterNodeInfo.Key, msg.NodeKey)
			res, err := T.Client.GetNodeClient().SendHeartbeat(context.Background(), T.Client.NewRequest(T.WorkerNodeInfo))
			require.Nil(t, err)
			require.NotNil(t, res)
		}
	}
}
