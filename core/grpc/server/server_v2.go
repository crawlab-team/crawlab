package server

import (
	"fmt"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/entity"
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/grpc/middlewares"
	"github.com/crawlab-team/crawlab/core/interfaces"
	nodeconfig "github.com/crawlab-team/crawlab/core/node/config"
	grpc2 "github.com/crawlab-team/crawlab/grpc"
	"github.com/crawlab-team/crawlab/trace"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	errors2 "github.com/pkg/errors"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"net"
	"sync"
)

var (
	subsV2      = map[string]interfaces.GrpcSubscribe{}
	mutexSubsV2 = &sync.Mutex{}
)

type GrpcServer struct {
	// settings
	cfgPath string
	address interfaces.Address

	// internals
	svr     *grpc.Server
	l       net.Listener
	stopped bool

	// dependencies
	nodeCfgSvc interfaces.NodeConfigService

	// servers
	NodeSvr             *NodeServiceServer
	TaskSvr             *TaskServiceServer
	ModelBaseServiceSvr *ModelBaseServiceServerV2
	DependencySvr       *DependencyServiceServer
	MetricSvr           *MetricServiceServer
}

func (svr *GrpcServer) GetConfigPath() (path string) {
	return svr.cfgPath
}

func (svr *GrpcServer) SetConfigPath(path string) {
	svr.cfgPath = path
}

func (svr *GrpcServer) Init() (err error) {
	// register
	if err := svr.Register(); err != nil {
		return err
	}

	return nil
}

func (svr *GrpcServer) Start() (err error) {
	// grpc server binding address
	address := svr.address.String()

	// listener
	svr.l, err = net.Listen("tcp", address)
	if err != nil {
		_ = trace.TraceError(err)
		return errors.ErrorGrpcServerFailedToListen
	}
	log.Infof("grpc server listens to %s", address)

	// start grpc server
	go func() {
		if err := svr.svr.Serve(svr.l); err != nil {
			if errors2.Is(err, grpc.ErrServerStopped) {
				return
			}
			trace.PrintError(err)
			log.Error(errors.ErrorGrpcServerFailedToServe.Error())
		}
	}()

	return nil
}

func (svr *GrpcServer) Stop() (err error) {
	// skip if listener is nil
	if svr.l == nil {
		return nil
	}

	// graceful stop
	log.Infof("grpc server stopping...")
	svr.svr.Stop()

	// close listener
	log.Infof("grpc server closing listener...")
	_ = svr.l.Close()

	// mark as stopped
	svr.stopped = true

	// log
	log.Infof("grpc server stopped")

	return nil
}

func (svr *GrpcServer) Register() (err error) {
	grpc2.RegisterNodeServiceServer(svr.svr, *svr.NodeSvr)
	grpc2.RegisterModelBaseServiceV2Server(svr.svr, *svr.ModelBaseServiceSvr)
	grpc2.RegisterTaskServiceServer(svr.svr, *svr.TaskSvr)
	grpc2.RegisterDependencyServiceV2Server(svr.svr, *svr.DependencySvr)
	grpc2.RegisterMetricServiceV2Server(svr.svr, *svr.MetricSvr)

	return nil
}

func (svr *GrpcServer) recoveryHandlerFunc(p interface{}) (err error) {
	err = errors.NewError(errors.ErrorPrefixGrpc, fmt.Sprintf("%v", p))
	trace.PrintError(err)
	return err
}

func (svr *GrpcServer) SetAddress(address interfaces.Address) {

}

func (svr *GrpcServer) GetSubscribe(key string) (sub interfaces.GrpcSubscribe, err error) {
	mutexSubsV2.Lock()
	defer mutexSubsV2.Unlock()
	sub, ok := subsV2[key]
	if !ok {
		return nil, errors.ErrorGrpcSubscribeNotExists
	}
	return sub, nil
}

func (svr *GrpcServer) SetSubscribe(key string, sub interfaces.GrpcSubscribe) {
	mutexSubsV2.Lock()
	defer mutexSubsV2.Unlock()
	subsV2[key] = sub
}

func (svr *GrpcServer) DeleteSubscribe(key string) {
	mutexSubsV2.Lock()
	defer mutexSubsV2.Unlock()
	delete(subsV2, key)
}

func (svr *GrpcServer) IsStopped() (res bool) {
	return svr.stopped
}

func NewGrpcServerV2() (svr *GrpcServer, err error) {
	// server
	svr = &GrpcServer{
		address: entity.NewAddress(&entity.AddressOptions{
			Host: constants.DefaultGrpcServerHost,
			Port: constants.DefaultGrpcServerPort,
		}),
	}

	if viper.GetString("grpc.server.address") != "" {
		svr.address, err = entity.NewAddressFromString(viper.GetString("grpc.server.address"))
		if err != nil {
			return nil, err
		}
	}

	svr.nodeCfgSvc = nodeconfig.GetNodeConfigService()

	svr.NodeSvr, err = NewNodeServiceServer()
	if err != nil {
		return nil, err
	}
	svr.ModelBaseServiceSvr = NewModelBaseServiceV2Server()
	svr.TaskSvr, err = NewTaskServiceServer()
	if err != nil {
		return nil, err
	}
	svr.DependencySvr = GetDependencyServerV2()
	svr.MetricSvr = GetMetricsServerV2()

	// recovery options
	recoveryOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(svr.recoveryHandlerFunc),
	}

	// grpc server
	svr.svr = grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_recovery.UnaryServerInterceptor(recoveryOpts...),
			grpc_auth.UnaryServerInterceptor(middlewares.GetAuthTokenFunc(svr.nodeCfgSvc)),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_recovery.StreamServerInterceptor(recoveryOpts...),
			grpc_auth.StreamServerInterceptor(middlewares.GetAuthTokenFunc(svr.nodeCfgSvc)),
		),
	)

	// initialize
	if err := svr.Init(); err != nil {
		return nil, err
	}

	return svr, nil
}

var _serverV2 *GrpcServer

func GetGrpcServerV2() (svr *GrpcServer, err error) {
	if _serverV2 != nil {
		return _serverV2, nil
	}
	_serverV2, err = NewGrpcServerV2()
	if err != nil {
		return nil, err
	}
	return _serverV2, nil
}
