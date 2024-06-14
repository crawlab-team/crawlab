package server

import (
	"context"
	grpc "github.com/crawlab-team/crawlab/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DependenciesServerV2 struct {
	grpc.UnimplementedDependencyServiceV2Server
}

func (svr DependenciesServerV2) Connect(stream grpc.DependencyServiceV2_ConnectServer) (err error) {
	return status.Errorf(codes.Unimplemented, "method Connect not implemented")
}

func (svr DependenciesServerV2) Sync(ctx context.Context, request *grpc.DependenciesServiceV2SyncRequest) (response *grpc.Response, err error) {
	return nil, status.Errorf(codes.Unimplemented, "method Sync not implemented")
}

func (svr DependenciesServerV2) Install(stream grpc.DependencyServiceV2_InstallServer) (err error) {
	return status.Errorf(codes.Unimplemented, "method Install not implemented")
}

func (svr DependenciesServerV2) UninstallDependencies(stream grpc.DependencyServiceV2_UninstallDependenciesServer) (err error) {
	return status.Errorf(codes.Unimplemented, "method UninstallDependencies not implemented")
}

func NewDependenciesServerV2() *DependenciesServerV2 {
	return &DependenciesServerV2{}
}
