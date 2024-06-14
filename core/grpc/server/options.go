package server

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
)

type Option func(svr interfaces.GrpcServer)

func WithConfigPath(path string) Option {
	return func(svr interfaces.GrpcServer) {
		svr.SetConfigPath(path)
	}
}

func WithAddress(address interfaces.Address) Option {
	return func(svr interfaces.GrpcServer) {
		svr.SetAddress(address)
	}
}

type NodeServerOption func(svr *NodeServer)

func WithServerNodeServerService(server interfaces.GrpcServer) NodeServerOption {
	return func(svr *NodeServer) {
		svr.server = server
	}
}

type TaskServerOption func(svr *TaskServer)

func WithServerTaskServerService(server interfaces.GrpcServer) TaskServerOption {
	return func(svr *TaskServer) {
		svr.server = server
	}
}

type MessageServerOption func(svr *MessageServer)

func WithServerMessageServerService(server interfaces.GrpcServer) MessageServerOption {
	return func(svr *MessageServer) {
		svr.server = server
	}
}
