package water

import (
	"context"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"strings"
	"time"
)

type Service interface {
	Endpoint() Endpoint
	DecodeRequest(ctx context.Context, req interface{}) (interface{}, error)
	EncodeResponse(ctx context.Context, resp interface{}) (interface{}, error)
	ServerBefore() ServerOption
	ServerAfter() ServerOption
	GetLogger() grpclog.DepthLoggerV2
	Name() string
}

func NewHandler(svc Service, options ...ServerOption) *Server {
	return NewServer(
		PrintRpcCall(svc.GetLogger(), svc.Endpoint()),
		svc.DecodeRequest,
		svc.EncodeResponse,
		append(options, svc.ServerBefore(), svc.ServerAfter())...,
	)
}

func PrintRpcCall(logger grpclog.DepthLoggerV2, e Endpoint) Endpoint {
	return func(ctx context.Context, request interface{}) (res interface{}, err error) {
		start := time.Now()
		defer func() {
			var traceID []string
			if md, ok := metadata.FromIncomingContext(ctx); ok {
				if t, yes := md["trace-id"]; yes {
					traceID = t
				}
			}

			if err != nil {
				logger.Errorf("%dms %s %+v <%s>", time.Since(start).Milliseconds(), strings.Join(traceID, "."), request, err.Error())
			} else {
				logger.Infof("%dms %s %+v", time.Since(start).Milliseconds(), strings.Join(traceID, "."), request)
			}
		}()

		res, err = e(ctx, request)
		return res, err
	}
}
