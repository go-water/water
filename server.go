package water

import (
	"context"
	"github.com/go-water/water/circuitbreaker"
	"github.com/go-water/water/endpoint"
	"github.com/go-water/water/ratelimit"
	"github.com/sony/gobreaker"
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
	breaker      *gobreaker.CircuitBreaker
}

func NewHandler(srv Service, options ...ServerOption) Handler {
	s := new(Server)
	for _, option := range options {
		option(s)
	}

	s.e = srv.Endpoint()
	if s.limit != nil {
		s.e = ratelimit.NewErrorLimiter(s.limit)(s.e)
	}
	if s.breaker != nil {
		s.e = circuitbreaker.GoBreaker(s.breaker)(s.e)
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
