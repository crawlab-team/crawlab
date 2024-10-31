package server

import (
	"fmt"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/entity"
	"github.com/crawlab-team/crawlab/core/grpc/middlewares"
	"github.com/crawlab-team/crawlab/core/interfaces"
	nodeconfig "github.com/crawlab-team/crawlab/core/node/config"
	grpc2 "github.com/crawlab-team/crawlab/grpc"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcauth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	errors2 "github.com/pkg/errors"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"net"
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
	if err := svr.register(); err != nil {
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
		log.Errorf("[GrpcServer] failed to listen: %v", err)
		return err
	}
	log.Infof("[GrpcServer] grpc server listens to %s", address)

	// start grpc server
	go func() {
		if err := svr.svr.Serve(svr.l); err != nil {
			if errors2.Is(err, grpc.ErrServerStopped) {
				return
			}
			log.Errorf("[GrpcServer] failed to serve: %v", err)
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
	log.Infof("[GrpcServer] grpc server stopping...")
	svr.svr.Stop()

	// close listener
	log.Infof("[GrpcServer] grpc server closing listener...")
	_ = svr.l.Close()

	// mark as stopped
	svr.stopped = true

	// log
	log.Infof("[GrpcServer] grpc server stopped")

	return nil
}

func (svr *GrpcServer) register() (err error) {
	grpc2.RegisterNodeServiceServer(svr.svr, *svr.NodeSvr)
	grpc2.RegisterModelBaseServiceV2Server(svr.svr, *svr.ModelBaseServiceSvr)
	grpc2.RegisterTaskServiceServer(svr.svr, *svr.TaskSvr)
	grpc2.RegisterDependencyServiceV2Server(svr.svr, *svr.DependencySvr)
	grpc2.RegisterMetricServiceV2Server(svr.svr, *svr.MetricSvr)

	return nil
}

func (svr *GrpcServer) recoveryHandlerFunc(p interface{}) (err error) {
	log.Errorf("[GrpcServer] recovered from panic: %v", p)
	return fmt.Errorf("recovered from panic: %v", p)
}

func NewGrpcServer() (svr *GrpcServer, err error) {
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
	recoveryOpts := []grpcrecovery.Option{
		grpcrecovery.WithRecoveryHandler(svr.recoveryHandlerFunc),
	}

	// grpc server
	svr.svr = grpc.NewServer(
		grpcmiddleware.WithUnaryServerChain(
			grpcrecovery.UnaryServerInterceptor(recoveryOpts...),
			grpcauth.UnaryServerInterceptor(middlewares.GetAuthTokenFunc(svr.nodeCfgSvc)),
		),
		grpcmiddleware.WithStreamServerChain(
			grpcrecovery.StreamServerInterceptor(recoveryOpts...),
			grpcauth.StreamServerInterceptor(middlewares.GetAuthTokenFunc(svr.nodeCfgSvc)),
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
	_serverV2, err = NewGrpcServer()
	if err != nil {
		return nil, err
	}
	return _serverV2, nil
}
