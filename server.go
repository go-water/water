package water

import (
	"context"
	"go.uber.org/zap"
	"sync"
)

type Handler interface {
	ServerWater(ctx context.Context, req any) (any, error)
	GetLogger() *zap.Logger
}

type Server struct {
	c            *Config
	e            Endpoint
	finalizer    []ServerFinalizerFunc
	errorHandler ErrorHandler
	l            *zap.Logger
	p            sync.Pool
}

func NewHandler(srv Service, options ...ServerOption) Handler {
	s := &Server{
		e: srv.Endpoint(),
		c: new(Config),
	}

	for _, option := range options {
		option(s)
	}

	s.p.New = func() any {
		return srv.GetRequest()
	}

	handler := NewLogErrorHandler(s.c.NewLogger(), srv.Name())
	srv.SetLogger(handler.l)
	s.l = handler.l
	s.errorHandler = handler

	return s
}

func (s *Server) ServerWater(ctx context.Context, req any) (resp any, err error) {
	defer func() {
		s.p.Put(req)
	}()

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

func (s *Server) GetRequest() any {
	return s.p.Get()
}
