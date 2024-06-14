package server

import (
	"encoding/json"
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
	"go/types"
	"google.golang.org/grpc"
	"net"
	"sync"
)

var (
	subsV2      = map[string]interfaces.GrpcSubscribe{}
	mutexSubsV2 = &sync.Mutex{}
)

type GrpcServerV2 struct {
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
	nodeSvr             *NodeServerV2
	taskSvr             *TaskServerV2
	modelBaseServiceSvr *ModelBaseServiceServerV2
	dependenciesSvr     *DependenciesServerV2
}

func (svr *GrpcServerV2) GetConfigPath() (path string) {
	return svr.cfgPath
}

func (svr *GrpcServerV2) SetConfigPath(path string) {
	svr.cfgPath = path
}

func (svr *GrpcServerV2) Init() (err error) {
	// register
	if err := svr.Register(); err != nil {
		return err
	}

	return nil
}

func (svr *GrpcServerV2) Start() (err error) {
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

func (svr *GrpcServerV2) Stop() (err error) {
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

func (svr *GrpcServerV2) Register() (err error) {
	grpc2.RegisterNodeServiceServer(svr.svr, *svr.nodeSvr) // node service
	grpc2.RegisterModelBaseServiceV2Server(svr.svr, *svr.modelBaseServiceSvr)
	grpc2.RegisterTaskServiceServer(svr.svr, *svr.taskSvr)

	return nil
}

func (svr *GrpcServerV2) recoveryHandlerFunc(p interface{}) (err error) {
	err = errors.NewError(errors.ErrorPrefixGrpc, fmt.Sprintf("%v", p))
	trace.PrintError(err)
	return err
}

func (svr *GrpcServerV2) SetAddress(address interfaces.Address) {

}

func (svr *GrpcServerV2) GetSubscribe(key string) (sub interfaces.GrpcSubscribe, err error) {
	mutexSubsV2.Lock()
	defer mutexSubsV2.Unlock()
	sub, ok := subsV2[key]
	if !ok {
		return nil, errors.ErrorGrpcSubscribeNotExists
	}
	return sub, nil
}

func (svr *GrpcServerV2) SetSubscribe(key string, sub interfaces.GrpcSubscribe) {
	mutexSubsV2.Lock()
	defer mutexSubsV2.Unlock()
	subsV2[key] = sub
}

func (svr *GrpcServerV2) DeleteSubscribe(key string) {
	mutexSubsV2.Lock()
	defer mutexSubsV2.Unlock()
	delete(subsV2, key)
}

func (svr *GrpcServerV2) SendStreamMessage(key string, code grpc2.StreamMessageCode) (err error) {
	return svr.SendStreamMessageWithData(key, code, nil)
}

func (svr *GrpcServerV2) SendStreamMessageWithData(key string, code grpc2.StreamMessageCode, d interface{}) (err error) {
	var data []byte
	switch d.(type) {
	case types.Nil:
		// do nothing
	case []byte:
		data = d.([]byte)
	default:
		var err error
		data, err = json.Marshal(d)
		if err != nil {
			return err
		}
	}
	sub, err := svr.GetSubscribe(key)
	if err != nil {
		return err
	}
	msg := &grpc2.StreamMessage{
		Code: code,
		Key:  svr.nodeCfgSvc.GetNodeKey(),
		Data: data,
	}
	return sub.GetStream().Send(msg)
}

func (svr *GrpcServerV2) IsStopped() (res bool) {
	return svr.stopped
}

func NewGrpcServerV2() (svr *GrpcServerV2, err error) {
	// server
	svr = &GrpcServerV2{
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

	svr.nodeSvr, err = NewNodeServerV2()
	if err != nil {
		return nil, err
	}
	svr.modelBaseServiceSvr = NewModelBaseServiceV2Server()
	svr.taskSvr, err = NewTaskServerV2()
	if err != nil {
		return nil, err
	}
	svr.dependenciesSvr = NewDependenciesServerV2()

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

var _serverV2 *GrpcServerV2

func GetGrpcServerV2() (svr *GrpcServerV2, err error) {
	if _serverV2 != nil {
		return _serverV2, nil
	}
	_serverV2, err = NewGrpcServerV2()
	if err != nil {
		return nil, err
	}
	return _serverV2, nil
}
