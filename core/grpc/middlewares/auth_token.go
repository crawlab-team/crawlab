package middlewares

import (
	"context"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func GetAuthTokenFunc(nodeCfgSvc interfaces.NodeConfigService) grpc_auth.AuthFunc {
	return func(ctx context.Context) (ctx2 context.Context, err error) {
		// authentication (token verification)
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errors.ErrorGrpcUnauthorized
		}

		// auth key from incoming context
		res, ok := md[constants.GrpcHeaderAuthorization]
		if !ok {
			return ctx, errors.ErrorGrpcUnauthorized
		}
		if len(res) != 1 {
			return ctx, errors.ErrorGrpcUnauthorized
		}
		authKey := res[0]

		// validate
		svrAuthKey := nodeCfgSvc.GetAuthKey()
		if authKey != svrAuthKey {
			return ctx, errors.ErrorGrpcUnauthorized
		}

		return ctx, nil
	}
}

func GetAuthTokenUnaryChainInterceptor(nodeCfgSvc interfaces.NodeConfigService) grpc.UnaryClientInterceptor {
	// set auth key
	md := metadata.Pairs(constants.GrpcHeaderAuthorization, nodeCfgSvc.GetAuthKey())
	//header := metadata.MD{}
	//header[constants.GrpcHeaderAuthorization] = []string{nodeCfgSvc.GetAuthKey()}
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = metadata.NewOutgoingContext(context.Background(), md)
		//opts = append(opts, grpc.Header(&header))
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func GetAuthTokenStreamChainInterceptor(nodeCfgSvc interfaces.NodeConfigService) grpc.StreamClientInterceptor {
	// set auth key
	md := metadata.Pairs(constants.GrpcHeaderAuthorization, nodeCfgSvc.GetAuthKey())
	//header := metadata.MD{}
	//header[constants.GrpcHeaderAuthorization] = []string{nodeCfgSvc.GetAuthKey()}
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		ctx = metadata.NewOutgoingContext(context.Background(), md)
		//opts = append(opts, grpc.Header(&header))
		s, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			return nil, err
		}
		return s, nil
	}
}
