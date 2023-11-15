package water

import (
	"context"
	"github.com/go-water/water/endpoint"
	"github.com/go-water/water/ratelimit"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

type Handler interface {
	ServerWater(ctx context.Context, req any) (any, error)
	GetLogger() *zap.Logger
}

type Server struct {
	e            endpoint.Endpoint
	finalizer    []ServerFinalizerFunc
	errorHandler ErrorHandler
	l            *zap.Logger
	limit        *rate.Limiter
}

func NewHandler(srv Service, options ...ServerOption) Handler {
	s := new(Server)
	for _, option := range options {
		option(s)
	}

	if s.limit != nil {
		s.e = ratelimit.NewErrorLimiter(s.limit)(srv.Endpoint())
	} else {
		s.e = srv.Endpoint()
	}

	handler := NewLogErrorHandler(log, srv.Name(srv))
	srv.SetLogger(handler.l)
	s.l = handler.l
	s.errorHandler = handler

	return s
}

func (s *Server) ServerWater(ctx context.Context, req any) (resp any, err error) {
	if len(s.finalizer) > 0 {
		defer func() {
			for _, fn := range s.finalizer {
				fn(ctx, err)
			}
		}()
	}

	resp, err = s.e(ctx, req)
	if err != nil {
		s.errorHandler.Handle(ctx, err)
		return nil, err
	}

	return resp, nil
}

func (s *Server) GetLogger() *zap.Logger {
	return s.l
}
