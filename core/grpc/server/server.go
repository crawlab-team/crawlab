package server

import (
	"encoding/json"
	"fmt"
	"github.com/apex/log"
	config2 "github.com/crawlab-team/crawlab/core/config"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/container"
	"github.com/crawlab-team/crawlab/core/entity"
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/grpc/middlewares"
	"github.com/crawlab-team/crawlab/core/interfaces"
	grpc2 "github.com/crawlab-team/crawlab/grpc"
	"github.com/crawlab-team/crawlab/trace"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/spf13/viper"
	"go/types"
	"google.golang.org/grpc"
	"net"
	"sync"
)

var subs = sync.Map{}

type Server struct {
	// dependencies
	nodeCfgSvc          interfaces.NodeConfigService
	nodeSvr             *NodeServer
	taskSvr             *TaskServer
	messageSvr          *MessageServer
	modelDelegateSvr    *ModelDelegateServer
	modelBaseServiceSvr *ModelBaseServiceServer

	// settings
	cfgPath string
	address interfaces.Address

	// internals
	svr     *grpc.Server
	l       net.Listener
	stopped bool
}

func (svr *Server) Init() (err error) {
	// register
	if err := svr.Register(); err != nil {
		return err
	}

	return nil
}

func (svr *Server) Start() (err error) {
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
			if err == grpc.ErrServerStopped {
				return
			}
			trace.PrintError(err)
			log.Error(errors.ErrorGrpcServerFailedToServe.Error())
		}
	}()

	return nil
}

func (svr *Server) Stop() (err error) {
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

func (svr *Server) Register() (err error) {
	grpc2.RegisterModelDelegateServer(svr.svr, *svr.modelDelegateSvr)       // model delegate
	grpc2.RegisterModelBaseServiceServer(svr.svr, *svr.modelBaseServiceSvr) // model base service
	grpc2.RegisterNodeServiceServer(svr.svr, *svr.nodeSvr)                  // node service
	grpc2.RegisterTaskServiceServer(svr.svr, *svr.taskSvr)                  // task service
	grpc2.RegisterMessageServiceServer(svr.svr, *svr.messageSvr)            // message service

	return nil
}

func (svr *Server) SetAddress(address interfaces.Address) {
	svr.address = address
}

func (svr *Server) GetConfigPath() (path string) {
	return svr.cfgPath
}

func (svr *Server) SetConfigPath(path string) {
	svr.cfgPath = path
}

func (svr *Server) GetSubscribe(key string) (sub interfaces.GrpcSubscribe, err error) {
	res, ok := subs.Load(key)
	if !ok {
		return nil, trace.TraceError(errors.ErrorGrpcStreamNotFound)
	}
	sub, ok = res.(interfaces.GrpcSubscribe)
	if !ok {
		return nil, trace.TraceError(errors.ErrorGrpcInvalidType)
	}
	return sub, nil
}

func (svr *Server) SetSubscribe(key string, sub interfaces.GrpcSubscribe) {
	subs.Store(key, sub)
}

func (svr *Server) DeleteSubscribe(key string) {
	subs.Delete(key)
}

func (svr *Server) SendStreamMessage(key string, code grpc2.StreamMessageCode) (err error) {
	return svr.SendStreamMessageWithData(key, code, nil)
}

func (svr *Server) SendStreamMessageWithData(key string, code grpc2.StreamMessageCode, d interface{}) (err error) {
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
			panic(err)
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

func (svr *Server) IsStopped() (res bool) {
	return svr.stopped
}

func (svr *Server) recoveryHandlerFunc(p interface{}) (err error) {
	err = errors.NewError(errors.ErrorPrefixGrpc, fmt.Sprintf("%v", p))
	trace.PrintError(err)
	return err
}

func NewServer() (svr2 interfaces.GrpcServer, err error) {
	// server
	svr := &Server{
		cfgPath: config2.GetConfigPath(),
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

	// dependency injection
	if err := container.GetContainer().Invoke(func(
		nodeCfgSvc interfaces.NodeConfigService,
		modelDelegateSvr *ModelDelegateServer,
		modelBaseServiceSvr *ModelBaseServiceServer,
		nodeSvr *NodeServer,
		taskSvr *TaskServer,
		messageSvr *MessageServer,
	) {
		// dependencies
		svr.nodeCfgSvc = nodeCfgSvc
		svr.modelDelegateSvr = modelDelegateSvr
		svr.modelBaseServiceSvr = modelBaseServiceSvr
		svr.nodeSvr = nodeSvr
		svr.taskSvr = taskSvr
		svr.messageSvr = messageSvr

		// server
		svr.nodeSvr.server = svr
		svr.taskSvr.server = svr
		svr.messageSvr.server = svr
	}); err != nil {
		return nil, err
	}

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

var _server interfaces.GrpcServer

func GetServer() (svr interfaces.GrpcServer, err error) {
	if _server != nil {
		return _server, nil
	}
	_server, err = NewServer()
	if err != nil {
		return nil, err
	}
	return _server, nil
}
