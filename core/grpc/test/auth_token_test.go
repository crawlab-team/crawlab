package test

import (
	"context"
	"encoding/json"
	"github.com/crawlab-team/crawlab/core/entity"
	"github.com/crawlab-team/crawlab/core/grpc/client"
	"github.com/crawlab-team/crawlab/core/grpc/server"
	"github.com/crawlab-team/crawlab/core/node/config"
	grpc "github.com/crawlab-team/crawlab/grpc"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestAuthToken(t *testing.T) {
	var err error

	// auth key
	authKey := "test-auth-key"

	// auth key (invalid)
	authKeyInvalid := "test-auth-key-invalid"

	// tmp dir
	tmpDir := os.TempDir()

	// master config
	masterConfigPath := path.Join(tmpDir, "config-master.json")
	masterConfig := config.Config{
		Key:      "master",
		IsMaster: true,
		AuthKey:  authKey,
	}
	masterConfigData, err := json.Marshal(&masterConfig)
	require.Nil(t, err)
	err = ioutil.WriteFile(masterConfigPath, masterConfigData, os.FileMode(0777))

	// worker config
	workerConfigPath := path.Join(tmpDir, "config-worker.json")
	workerConfig := config.Config{
		Key:      "worker",
		IsMaster: false,
		AuthKey:  authKey,
	}
	workerConfigData, err := json.Marshal(&workerConfig)
	require.Nil(t, err)
	err = ioutil.WriteFile(workerConfigPath, workerConfigData, os.FileMode(0777))

	// worker config (invalid)
	workerInvalidConfigPath := path.Join(tmpDir, "worker-invalid")
	workerInvalidConfig := config.Config{
		Key:      "worker",
		IsMaster: false,
		AuthKey:  authKeyInvalid,
	}
	workerInvalidConfigData, err := json.Marshal(&workerInvalidConfig)
	require.Nil(t, err)
	err = ioutil.WriteFile(workerInvalidConfigPath, workerInvalidConfigData, os.FileMode(0777))

	// server
	svr, err := server.NewServer(
		server.WithConfigPath(masterConfigPath),
		server.WithAddress(entity.NewAddress(&entity.AddressOptions{
			Host: "0.0.0.0",
			Port: "9999",
		})),
	)
	require.Nil(t, err)
	err = svr.Start()
	require.Nil(t, err)

	// client
	c, err := client.GetClient(workerConfigPath, client.WithAddress(entity.NewAddress(&entity.AddressOptions{
		Host: "localhost",
		Port: "9999",
	})))
	require.Nil(t, err)
	err = c.Start()
	require.Nil(t, err)
	_, err = c.GetNodeClient().Ping(context.Background(), &grpc.Request{NodeKey: workerConfig.Key})
	require.Nil(t, err)

	// client (invalid)
	ci, err := client.GetClient(workerInvalidConfigPath, client.WithAddress(entity.NewAddress(&entity.AddressOptions{
		Host: "localhost",
		Port: "9999",
	})))
	require.Nil(t, err)
	err = ci.Start()
	require.Nil(t, err)
	_, err = ci.GetNodeClient().Ping(context.Background(), &grpc.Request{NodeKey: workerInvalidConfig.Key})
	require.NotNil(t, err)
}
