package client

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/cenkalti/backoff/v4"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/entity"
	"github.com/crawlab-team/crawlab/core/grpc/middlewares"
	"github.com/crawlab-team/crawlab/core/interfaces"
	nodeconfig "github.com/crawlab-team/crawlab/core/node/config"
	"github.com/crawlab-team/crawlab/core/utils"
	grpc2 "github.com/crawlab-team/crawlab/grpc"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	// dependencies
	nodeCfgSvc interfaces.NodeConfigService

	// settings
	address interfaces.Address
	timeout time.Duration

	// internals
	conn    *grpc.ClientConn
	err     error
	once    sync.Once
	stopped bool
	stop    chan struct{}

	// clients
	NodeClient               grpc2.NodeServiceClient
	TaskClient               grpc2.TaskServiceClient
	ModelBaseServiceV2Client grpc2.ModelBaseServiceV2Client
	DependencyClient         grpc2.DependencyServiceV2Client
	MetricClient             grpc2.MetricServiceV2Client
}

func (c *GrpcClient) Start() (err error) {
	c.once.Do(func() {
		// connect
		err = c.connect()
		if err != nil {
			return
		}

		// register rpc services
		c.register()
	})

	return err
}

func (c *GrpcClient) Stop() (err error) {
	// set stopped flag
	c.stopped = true
	c.stop <- struct{}{}
	log.Infof("[GrpcClient] stopped")

	// skip if connection is nil
	if c.conn == nil {
		return nil
	}

	// grpc server address
	address := c.address.String()

	// close connection
	if err := c.conn.Close(); err != nil {
		return err
	}
	log.Infof("grpc client disconnected from %s", address)

	return nil
}

func (c *GrpcClient) register() {
	c.NodeClient = grpc2.NewNodeServiceClient(c.conn)
	c.ModelBaseServiceV2Client = grpc2.NewModelBaseServiceV2Client(c.conn)
	c.TaskClient = grpc2.NewTaskServiceClient(c.conn)
	c.DependencyClient = grpc2.NewDependencyServiceV2Client(c.conn)
	c.MetricClient = grpc2.NewMetricServiceV2Client(c.conn)
}

func (c *GrpcClient) Context() (ctx context.Context, cancel context.CancelFunc) {
	return context.WithTimeout(context.Background(), c.timeout)
}

func (c *GrpcClient) IsStarted() (res bool) {
	return c.conn != nil
}

func (c *GrpcClient) IsClosed() (res bool) {
	if c.conn != nil {
		return c.conn.GetState() == connectivity.Shutdown
	}
	return false
}

func (c *GrpcClient) getRequestData(d interface{}) (data []byte) {
	if d == nil {
		return data
	}
	switch d.(type) {
	case []byte:
		data = d.([]byte)
	default:
		var err error
		data, err = json.Marshal(d)
		if err != nil {
			panic(err)
		}
	}
	return data
}

func (c *GrpcClient) connect() (err error) {
	op := func() error {
		// grpc server address
		address := c.address.String()

		// connection options
		opts := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithChainUnaryInterceptor(middlewares.GetAuthTokenUnaryChainInterceptor(c.nodeCfgSvc)),
			grpc.WithChainStreamInterceptor(middlewares.GetAuthTokenStreamChainInterceptor(c.nodeCfgSvc)),
		}

		// create new client connection
		c.conn, err = grpc.NewClient(address, opts...)
		if err != nil {
			log.Errorf("[GrpcClient] grpc client failed to connect to %s: %v", address, err)
			return err
		}

		log.Infof("[GrpcClient] grpc client connected to %s", address)

		return nil
	}
	return backoff.RetryNotify(op, backoff.NewExponentialBackOff(), utils.BackoffErrorNotify("grpc client connect"))
}

func newGrpcClient() (c *GrpcClient) {
	client := &GrpcClient{
		address: entity.NewAddress(&entity.AddressOptions{
			Host: constants.DefaultGrpcClientRemoteHost,
			Port: constants.DefaultGrpcClientRemotePort,
		}),
		timeout: 10 * time.Second,
		stop:    make(chan struct{}),
	}
	client.nodeCfgSvc = nodeconfig.GetNodeConfigService()

	if viper.GetString("grpc.address") != "" {
		address, err := entity.NewAddressFromString(viper.GetString("grpc.address"))
		if err != nil {
			log.Errorf("failed to parse grpc address: %s", viper.GetString("grpc.address"))
			panic(err)
		}
		client.address = address
	}

	return client
}

var clientV2 *GrpcClient
var clientV2Once sync.Once

func GetGrpcClient() *GrpcClient {
	clientV2Once.Do(func() {
		clientV2 = newGrpcClient()
	})
	return clientV2
}
