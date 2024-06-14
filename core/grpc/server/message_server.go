package server

import (
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab/core/container"
	"github.com/crawlab-team/crawlab/core/entity"
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/service"
	grpc "github.com/crawlab-team/crawlab/grpc"
	"github.com/crawlab-team/crawlab/trace"
	"io"
)

type MessageServer struct {
	grpc.UnimplementedMessageServiceServer

	// dependencies
	modelSvc service.ModelService
	cfgSvc   interfaces.NodeConfigService

	// internals
	server interfaces.GrpcServer
}

func (svr MessageServer) Connect(stream grpc.MessageService_ConnectServer) (err error) {
	finished := make(chan bool)
	for {
		msg, err := stream.Recv()
		nodeKey := "unknown node key"
		if msg != nil {
			nodeKey = msg.NodeKey
		}
		if err == io.EOF {
			log.Infof("[MessageServer] received signal EOF from node[%s], now quit", nodeKey)
			return nil
		}
		if err != nil {
			log.Errorf("[MessageServer] receiving message error from node[%s]: %v", nodeKey, err)
			return err
		}
		switch msg.Code {
		case grpc.StreamMessageCode_CONNECT:
			log.Infof("[MessageServer] received connect request from node[%s], key: %s", nodeKey, msg.Key)
			svr.server.SetSubscribe(msg.Key, &entity.GrpcSubscribe{
				Stream:   stream,
				Finished: finished,
			})
		case grpc.StreamMessageCode_DISCONNECT:
			log.Infof("[MessageServer] received disconnect request from node[%s], key: %s", nodeKey, msg.Key)
			svr.server.DeleteSubscribe(msg.Key)
			return nil
		case grpc.StreamMessageCode_SEND:
			log.Debugf("[MessageServer] received send request from node[%s] to %s", nodeKey, msg.To)
			sub, err := svr.server.GetSubscribe(msg.To)
			if err != nil {
				return err
			}
			svr.redirectMessage(sub, msg)
		}
	}
}

func (svr MessageServer) redirectMessage(sub interfaces.GrpcSubscribe, msg *grpc.StreamMessage) {
	stream := sub.GetStream()
	if stream == nil {
		trace.PrintError(errors.ErrorGrpcStreamNotFound)
		return
	}
	log.Debugf("[MessageServer] redirect message: %v", msg)
	if err := stream.Send(msg); err != nil {
		trace.PrintError(err)
		return
	}
}

func NewMessageServer() (res *MessageServer, err error) {
	// message server
	svr := &MessageServer{}

	// dependency injection
	if err := container.GetContainer().Invoke(func(
		modelSvc service.ModelService,
		cfgSvc interfaces.NodeConfigService,
	) {
		svr.modelSvc = modelSvc
		svr.cfgSvc = cfgSvc
	}); err != nil {
		return nil, err
	}

	return svr, nil
}
